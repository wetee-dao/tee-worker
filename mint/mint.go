package mint

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/gen/system"
	"github.com/wetee-dao/go-sdk/gen/types"
	"wetee.app/worker/dao"
	"wetee.app/worker/util"
)

// Minter
// 矿工
type Minter struct {
	K8sClient     *kubernetes.Clientset
	MetricsClient *versioned.Clientset
	ChainClient   *chain.ChainClient
}

var (
	MinterIns *Minter
	lock      sync.Mutex
	Signer    *signature.KeyringPair
)

// InitMint
// 初始化矿工
func InitMint(mgr manager.Manager) error {
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

	client, err := chain.ClientInit()
	if err != nil {
		return err
	}

	// 初始化minter
	lock.Lock()
	MinterIns = &Minter{
		K8sClient:     clientset,
		MetricsClient: metricsClient,
		ChainClient:   client,
	}
	lock.Unlock()

	// 获取签名账户
	Signer, err = GetMintKey()
	return err
}

// start mint
// 开始挖矿
func (m *Minter) StartMint() {
	fmt.Println("MintKey => ", Signer.Address)
	client := m.ChainClient
	chainAPI := client.Api

	// 获取挖矿状态
	worker := chain.Worker{
		Client: client,
		Signer: Signer,
	}

	// 挖矿开始
mintStart:

	// 等待集群开启
	for {
		clusterId, err := worker.Getk8sClusterAccounts(Signer.PublicKey)
		if err != nil {
			fmt.Println("ClusterId => ", err)
			time.Sleep(time.Second * 10)
			continue
		}
		dao.SetClusterId(clusterId)

		break
	}

	clusterId, _ := dao.GetClusterId()
	fmt.Println("ClusterId => ", clusterId)

	// 触发区块监听
	sub, err := chainAPI.RPC.Chain.SubscribeFinalizedHeads()

	// 失败后等待10秒重新尝试
	if err != nil {
		util.LogWithRed("SubscribeNewHeads", err)
		time.Sleep(time.Second * 10)
		goto mintStart
	}
	defer sub.Unsubscribe()

	for {
		head := <-sub.Chan()
		fmt.Printf("Chain is at block: #%v\n", head.Number)
		blockHash, _ := chainAPI.RPC.Chain.GetBlockHash(uint64(head.Number))

		// 读取/处理新的区块信息
		events, err := system.GetEvents(chainAPI.RPC.State, blockHash)
		if err != nil {
			util.LogWithRed("GetEventsLatest", err)
			continue
		}

		for _, event := range events {
			e := event.Event
			if e.IsWeteeWorker {
				startEvent := e.AsWeteeWorkerField0
				if startEvent.IsWorkRuning {
					workId := startEvent.AsWorkRuningWorkId1
					user := startEvent.AsWorkRuningUser0
					cid := startEvent.AsWorkRuningClusterId2
					if cid == clusterId {
						fmt.Println("===========================================WorkRuning: ", workId)

						version, _ := chain.GetVersion(m.ChainClient, workId)
						if workId.Wtype.IsAPP {
							ctx := context.Background()
							err = m.CreateApp(ctx, user[:], workId, version)
							fmt.Println("===========================================CreateOrUpdateApp error: ", err)
						} else {
							err = m.CreateTask(user[:], workId, version)
							fmt.Println("===========================================CreateOrUpdateTask error: ", err)
						}
					}
				}
			}
			if e.IsWeteeApp {
				appEvent := e.AsWeteeAppField0
				if appEvent.IsWorkStopped {
					workId := appEvent.AsWorkStoppedWorkId1

					fmt.Println("===========================================StopPod", workId)
					err := m.StopApp(workId)
					fmt.Println("===========================================StopPod error: ", err)
				}
				if appEvent.IsWorkUpdated {
					workId := appEvent.AsWorkUpdatedWorkId1
					user := appEvent.AsWorkUpdatedUser0

					fmt.Println("===========================================WorkUpdated: ", workId)
					version, _ := chain.GetVersion(m.ChainClient, workId)
					ctx := context.Background()
					err = m.UpdateApp(ctx, user[:], workId, blockHash, version)
					fmt.Println("===========================================CreateOrUpdatePod error: ", err)
				}
			}
		}

		// 获取合约列表
		cs, err := worker.GetClusterContracts(clusterId, &blockHash)
		fmt.Println("GetClusterContracts", cs)
		if err != nil {
			fmt.Println("GetClusterContracts", err)
			continue
		}
		stage, err := worker.GetStage()
		if err != nil {
			fmt.Println("GetStage", err)
			continue
		}

		// 校对合约状态
		for _, c := range cs {
			state, err := worker.GetWorkContract(c.ContractState.WorkId, clusterId)
			if err != nil {
				fmt.Println("GetWorkContract", err)
				continue
			}

			// 如果是APP类型，检查Pod状态，检查是否需要上传工作证明
			if c.ContractState.WorkId.Wtype.IsAPP {
				ctx := context.Background()
				_, err := m.CheckAppStatus(ctx, c)
				if err != nil {
					fmt.Println("checkPodStatus", err)
					continue
				}

				appIns := chain.App{
					Client: m.ChainClient,
					Signer: Signer,
				}
				app, err := appIns.GetApp(c.ContractState.User[:], c.ContractState.WorkId.Id)
				if err != nil {
					fmt.Println("appIns.GetApp", err)
					m.StopApp(c.ContractState.WorkId)
					continue
				}

				// 判断是否上传工作证明
				if uint64(app.Status) == 1 && uint64(head.Number)-state.BlockNumber >= uint64(stage) {
					fmt.Println("=========================================== WorkProofUpload APP")
					nameSpace := AccountToAddress(c.ContractState.User[:])
					workID := c.ContractState.WorkId
					name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

					// 获取log和硬件资源使用量
					ctx := context.Background()
					logs, crs, err := m.getMetricInfo(ctx, nameSpace, name, uint64(head.Number)-state.BlockNumber)
					if err != nil {
						fmt.Println("getMetricInfo", err)
						continue
					}

					// 获取log hash
					logHash, err := getWorkLogHash(name, logs, state.BlockNumber)
					if err != nil {
						fmt.Println("getWorkLogHash", err)
						continue
					}

					// 获取 cr hash
					crHash, cr, err := getWorkCrHash(name, crs, state.BlockNumber)
					if err != nil {
						fmt.Println("getWorkCrHash", err)
						continue
					}

					// 上传工作证明
					err = worker.WorkProofUpload(c.ContractState.WorkId, logHash, crHash, types.Cr{
						Cpu:  cr[0],
						Mem:  cr[1],
						Disk: 0,
					}, []byte(""))
					if err != nil {
						fmt.Println("WorkProofUpload", err)
						continue
					}
				}
			}

			// 如果是TASK类型，检查Pod状态，Pod如果执行完成，则上传日志和结果
			if c.ContractState.WorkId.Wtype.IsTASK {
				ctx := context.Background()
				pod, err := m.CheckTaskStatus(ctx, c)
				if err != nil {
					fmt.Println("checkTaskStatus", err)
					continue
				}

				// 判断是否上传工作证明
				if pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed {
					fmt.Println("===========================================WorkProofUpload TASK")
					nameSpace := AccountToAddress(c.ContractState.User[:])
					workID := c.ContractState.WorkId
					name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

					// 获取log和硬件资源使用量
					ctx := context.Background()
					logs, crs, err := m.getMetricInfo(ctx, nameSpace, name, uint64(head.Number)-state.BlockNumber)
					if err != nil {
						fmt.Println("getMetricInfo", err)
						continue
					}

					// 获取log hash
					logHash, err := getWorkLogHash(name, logs, state.BlockNumber)
					if err != nil {
						fmt.Println("getWorkLogHash", err)
						continue
					}

					// 获取 cr hash
					crHash, cr, err := getWorkCrHash(name, crs, state.BlockNumber)
					if err != nil {
						fmt.Println("getWorkCrHash", err)
						continue
					}

					// 上传工作证明结束任务
					err = worker.WorkProofUpload(c.ContractState.WorkId, logHash, crHash, types.Cr{
						Cpu:  cr[0],
						Mem:  cr[1],
						Disk: cr[2],
					}, []byte(""))
					if err != nil {
						fmt.Println("WorkProofUpload", err)
						continue
					}
				}
			}
		}
	}
}
