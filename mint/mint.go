package mint

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/module"
	"github.com/wetee-dao/go-sdk/pallet/system"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/internal/peer"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/mint/proof"
	types "wetee.app/worker/type"
	"wetee.app/worker/util"
)

// Minter
// 矿工
type Minter struct {
	K8sClient     *kubernetes.Clientset
	MetricsClient *versioned.Clientset
	ChainClient   *chain.ChainClient
	P2Peer        *peer.Peer
	Nodes         []*types.Node
	PrivateKey    *types.PrivKey
	HostDomain    string

	// App lanch event
	// 应用启动事件
	AppLanch []gtypes.WorkId
	mu       sync.RWMutex

	// preRecerve is the channel to receive SendEncryptedSecretRequest
	preRecerve map[string]chan interface{}
}

var (
	lock            sync.Mutex
	MinterIns       *Minter
	DefaultChainUrl string = "ws://wetee-node.worker-addon.svc.cluster.local:9944"
)

// InitCluster
// 初始化矿工
func InitCluster(mgr manager.Manager) error {
	// 创建K8s Client
	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}

	// 创建Metrics Client
	metricsClient, err := versioned.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}

	// 初始化minter
	lock.Lock()
	MinterIns = &Minter{
		K8sClient:     clientset,
		MetricsClient: metricsClient,
		ChainClient:   nil,
		HostDomain:    "",
		preRecerve:    make(map[string]chan interface{}),
	}

	// 获取签名账户
	_, privateKey, err := GetMintKey()
	if err != nil {
		return err
	}

	MinterIns.PrivateKey = privateKey
	lock.Unlock()

	// 此处不捕获错误，因为如果初始化失败，程序可以继续运行
	InitChainClient(DefaultChainUrl)

	return err
}

func InitChainClient(url string) error {
	if MinterIns.ChainClient != nil {
		return nil
	}
	client, err := chain.ClientInit(url, true)
	if err != nil {
		return err
	}
	MinterIns.ChainClient = client
	store.SetChainUrl(url)

	return nil
}

// start mint
// 开始挖矿
func (m *Minter) StartMint() {
	signer, _ := m.PrivateKey.ToSigner()
	fmt.Println("MintKey => ", signer.Address)

	// 挖矿开始
mintStart:

	var worker module.Worker

	// 等待集群开启
	// Waiting for cluster start
	for {
		// 此处不捕获错误，因为如果初始化失败，程序可以继续运行
		InitChainClient(DefaultChainUrl)
		if MinterIns.ChainClient == nil {
			fmt.Println("Chain connect is not init")
			time.Sleep(time.Second * 10)
			continue
		}

		// 启动p2p
		// Start p2p
		err := m.StartP2P()
		if m.P2Peer != nil {
			m.P2Peer.Discover(context.Background())
		}
		if err != nil {
			fmt.Println("worker.ClusterProofUpload => ", err)
			time.Sleep(time.Second * 10)
			continue
		}

		// 初始化worker对象
		// Initialize the worker object
		worker = module.Worker{
			Client: m.ChainClient,
			Signer: signer,
		}

		// 获取clusterId
		clusterId, err := worker.Getk8sClusterAccounts(signer.PublicKey)
		if err != nil {
			fmt.Println("ClusterId => clusterId not found, mint not started")
			time.Sleep(time.Second * 10)
			continue
		}

		// 获取 TEE 根证书
		report, t, err := proof.GetRemoteReport(signer, nil)
		if err != nil {
			fmt.Println("GetRootDcapReport => ", err)
			time.Sleep(time.Second * 10)
			continue
		}

		// 上传 TEE 证书
		// hash := blake2b.Sum256(report)
		param := types.TeeParam{
			Report:  report,
			Time:    t,
			TeeType: 0,
			Address: signer.SS58Address(42),
			Data:    nil,
		}
		_, err = m.UploadClusterProof(&param)
		if err != nil {
			fmt.Println("worker.ClusterProofUpload => ", err)
			time.Sleep(time.Second * 10)
			continue
		}

		// 获取集群域名
		cluster, err := worker.GetCluster(clusterId)
		if err != nil {
			fmt.Println("worker.GetCluster => ", err)
			time.Sleep(time.Second * 10)
			continue
		}
		MinterIns.HostDomain = string(cluster.Ip[0].Domain.AsSomeField0)

		// 保存clusterId
		store.SetClusterId(clusterId)

		break
	}

	clusterId, _ := store.GetClusterId()
	fmt.Println("ClusterId => ", clusterId)

	client := m.ChainClient
	chainAPI := client.Api

	// 订阅区块事件
	// Subscribe to block events
	sub, err := chainAPI.RPC.Chain.SubscribeNewHeads()
	if err != nil {
		util.LogError("SubscribeNewHeads", err)
		// 失败后等待10秒重新尝试
		// Wait 10 seconds to try again
		time.Sleep(time.Second * 10)
		goto mintStart
	}
	defer sub.Unsubscribe()

	for {
		head := <-sub.Chan()
		util.LogError("Chain is at block: #", fmt.Sprint(head.Number))
		blockHash, _ := chainAPI.RPC.Chain.GetBlockHash(uint64(head.Number))

		// P2P 节点发现 10个区块刷新一次
		// P2P node discovery
		if uint64(head.Number)%10 == 0 {
			m.P2Peer.Discover(context.Background())
			fmt.Println("Peer len:", len(m.P2Peer.Network().Peers()))
		}

		err := client.CheckMetadata()
		if err != nil {
			util.LogError("CheckMetadata", err)
			continue
		}

		// 读取/处理新的区块信息
		// Read/process new block information
		events, err := system.GetEvents(chainAPI.RPC.State, blockHash)
		if err != nil {
			util.LogError("GetEventsLatest", err)
			continue
		}

		// 处理事件
		// Processing event
		for _, event := range events {
			err := m.DoWithEvent(event, clusterId)
			if err != nil {
				util.LogError("DoWithEvent", err)
				continue
			}
		}

		// 获取合约列表
		// Get contract list
		cs, err := m.GetClusterContracts(clusterId, &blockHash)
		if err != nil {
			util.LogError("GetClusterContracts", err)
			continue
		}

		// 删除过期的合约
		// Delete expired contracts
		_, err = DeleteFormCache(cs, func(wid gtypes.WorkId, d store.RuningCache) error {
			name := util.GetWorkTypeStr(wid) + "-" + fmt.Sprint(wid.Id)
			err := m.StopApp(wid, d.NameSpace)
			if err != nil && !strings.Contains(err.Error(), "not found") {
				util.LogError("DeleteRuning "+name+" ", err)
				return err
			}
			return nil
		})
		if err != nil {
			util.LogError("DeleteFormCache", err)
			continue
		}

		// 获取收费周期
		// Get the charge cycle
		stage, err := worker.GetStage()
		if err != nil {
			util.LogError("GetStage", err)
			continue
		}

		fmt.Println("===========================================GetClusterContracts: ", len(cs))
		proofs := make([]gtypes.RuntimeCall, 0, 20)

		// 触发TEE调用
		// Trigger TEE calls
		m.trigger(cs, clusterId, uint64(head.Number))

		// 校对合约状态
		// Check contract status
		for _, c := range cs {
			ctx := context.Background()

			if c.ContractState.WorkId.Wtype.IsAPP {
				// 如果是APP类型，检查Pod状态，检查是否需要上传工作证明
				// If it is APP type, check Pod status, check if it needs to upload work proof
				call, err := m.DoWithAppState(&ctx, c, stage, head)
				if err != nil {
					util.LogError("DoWithAppState", err)
				}
				if call != nil {
					proofs = append(proofs, *call)
				}
			} else if c.ContractState.WorkId.Wtype.IsTASK {
				// 如果是TASK类型，检查Pod状态，Pod如果执行完成，则上传日志和结果
				// If it is TASK type, check Pod status, Pod if it is executed, upload logs and results
				call, err := m.DoWithTaskState(&ctx, c, stage, head)
				if err != nil {
					util.LogError("DoWithTaskState", err)
				}
				if call != nil {
					proofs = append(proofs, *call)
				}
			} else if c.ContractState.WorkId.Wtype.IsGPU {
				// 如果是GPU类型，检查Pod状态，检查是否需要上传工作证明
				call, err := m.DoWithGpuAppState(&ctx, c, stage, head)
				if err != nil {
					util.LogError("DoWithGpuAppState", err)
				}
				if call != nil {
					proofs = append(proofs, *call)
				}
			}
		}

		// 获取启动应用列表
		var lanchs = []gtypes.WorkId{}
		m.mu.Lock()
		lanchs = m.AppLanch
		m.AppLanch = []gtypes.WorkId{}
		m.mu.Unlock()
		for _, wid := range lanchs {
			call, err := proof.MakeStartProof(wid, uint64(head.Number))
			if err != nil {
				util.LogError("MakeStartProof", err)
				continue
			}
			proofs = append(proofs, *call)
		}

		if len(proofs) > 0 {
			// 上传工作证明
			// Upload work proof
			go func(b uint64) {
				err = proof.SubmitWorkProof(client, worker.Signer, proofs)
				if err != nil {
					util.LogError("WorkProofUpload", err)
				} else {
					fmt.Println("Proof.SubmitWorkProof blocknumber =>", b, "success")
				}
			}(uint64(head.Number))
		}
	}
}

func (m *Minter) Addlanch(wid gtypes.WorkId) {
	m.mu.Lock()
	m.AppLanch = append(m.AppLanch, wid)
	m.mu.Unlock()
}

func DeleteFormCache(cs map[gtypes.WorkId]ContractStateWrap, deleteFunc func(gtypes.WorkId, store.RuningCache) error) ([]gtypes.WorkId, error) {
	// 获取缓存
	caches, err := store.GetRuning()
	if err != nil {
		return nil, err
	}

	// 删除已经停止的应用
	var deletes = []gtypes.WorkId{}
	for name, cache := range caches {
		ids := strings.Split(name, "-")
		id, err := strconv.ParseUint(ids[1], 10, 64)
		if err != nil {
			return nil, err
		}

		wid := gtypes.WorkId{
			Wtype: util.GetWorkType(ids[0]),
			Id:    id,
		}

		if _, ok := cs[wid]; !ok {
			util.LogError("DeleteFormCache", fmt.Sprintf("Delete cache %s", name))
			if deleteFunc(wid, cache) == nil {
				delete(caches, name)
				deletes = append(deletes, wid)
			}
		}
	}

	// 重构新的缓存
	for workId, c := range cs {
		name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)
		nameSpace := AccountToSpace(c.ContractState.User[:])
		caches[name] = store.RuningCache{
			NameSpace: nameSpace,
			Status:    "running",
			DeleteAt:  0,
		}
	}

	// 设置新的缓存
	err = store.SetRuning(caches)
	if err != nil {
		return nil, err
	}

	return deletes, nil
}
