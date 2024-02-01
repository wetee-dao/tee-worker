package mint

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
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

	// 创建Chain Client
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

	// 获取签名账户
	Signer, err = GetMintKey()
	lock.Unlock()

	return err
}

// start mint
// 开始挖矿
func (m *Minter) StartMint() {
	fmt.Println("MintKey => ", Signer.Address)
	client := m.ChainClient
	chainAPI := client.Api

	// 初始化worker对象
	// Initialize the worker object
	worker := chain.Worker{
		Client: client,
		Signer: Signer,
	}

	// 挖矿开始
mintStart:

	// 等待集群开启
	// Waiting for cluster start
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

	// 订阅区块事件
	// Subscribe to block events
	sub, err := chainAPI.RPC.Chain.SubscribeNewHeads()
	if err != nil {
		util.LogWithRed("SubscribeNewHeads", err)
		// 失败后等待10秒重新尝试
		// Wait 10 seconds to try again
		time.Sleep(time.Second * 10)
		goto mintStart
	}
	defer sub.Unsubscribe()

	for {
		head := <-sub.Chan()
		fmt.Printf("Chain is at block: #%v\n", head.Number)
		blockHash, _ := chainAPI.RPC.Chain.GetBlockHash(uint64(head.Number))

		// 读取/处理新的区块信息
		// Read/process new block information
		events, err := system.GetEvents(chainAPI.RPC.State, blockHash)
		if err != nil {
			util.LogWithRed("GetEventsLatest", err)
			continue
		}

		// 处理事件
		// Processing event
		for _, event := range events {
			err := m.DoWithEvent(event, clusterId)
			if err != nil {
				continue
			}
		}

		// 获取合约列表
		// Get contract list
		cs, err := m.GetClusterContracts(clusterId, &blockHash)
		if err != nil {
			util.LogWithRed("GetClusterContracts", err)
			continue
		}

		// 获取收费周期
		// Get the charge cycle
		stage, err := worker.GetStage()
		if err != nil {
			util.LogWithRed("GetStage", err)
			continue
		}

		fmt.Println("===========================================GetClusterContracts: ", cs)

		// 校对合约状态
		// Check contract status
		for _, c := range cs {
			ctx := context.Background()

			// 如果是APP类型，检查Pod状态，检查是否需要上传工作证明
			// If it is APP type, check Pod status, check if it needs to upload work proof
			if c.ContractState.WorkId.Wtype.IsAPP {
				err := m.DoWithAppState(&ctx, c, stage, head)
				if err != nil {
					util.LogWithRed("DoWithAppState", err)
					continue
				}
			}

			// 如果是TASK类型，检查Pod状态，Pod如果执行完成，则上传日志和结果
			// If it is TASK type, check Pod status, Pod if it is executed, upload logs and results
			if c.ContractState.WorkId.Wtype.IsTASK {
				err := m.DoWithTaskState(&ctx, c, stage, head)
				if err != nil {
					util.LogWithRed("DoWithTaskState", err)
					continue
				}
			}
		}
	}
}

func (m *Minter) DoWithEvent(event types.EventRecord, clusterId uint64) error {
	e := event.Event
	ctx := context.Background()
	var err error

	// 处理任务消息
	// Handling Worker Messages
	if e.IsWeteeWorker {
		startEvent := e.AsWeteeWorkerField0
		if startEvent.IsWorkRuning {
			workId := startEvent.AsWorkRuningWorkId1
			user := startEvent.AsWorkRuningUser0
			cid := startEvent.AsWorkRuningClusterId2
			if cid == clusterId {
				version, _ := chain.GetVersion(m.ChainClient, workId)
				if workId.Wtype.IsAPP {
					appIns := chain.App{
						Client: m.ChainClient,
					}
					app, _ := appIns.GetApp(user[:], workId.Id)
					err = m.CreateApp(&ctx, user[:], workId, app, version)
					util.LogWithRed("===========================================CreateOrUpdateApp error: ", err)
				} else {
					taskIns := chain.Task{
						Client: m.ChainClient,
					}
					task, _ := taskIns.GetTask(user[:], workId.Id)
					err = m.CreateTask(&ctx, user[:], workId, task, version)
					util.LogWithRed("===========================================CreateOrUpdateTask error: ", err)
				}
			}
		}
	}

	// 处理机密应用消息
	// Handling App Messages
	if e.IsWeteeApp {
		appEvent := e.AsWeteeAppField0
		if appEvent.IsWorkStopped {
			workId := appEvent.AsWorkStoppedWorkId1

			util.LogWithRed("===========================================StopPod", workId)
			err = m.StopApp(workId)
			util.LogWithRed("===========================================StopPod error: ", err)
		}
		if appEvent.IsWorkUpdated {
			workId := appEvent.AsWorkUpdatedWorkId1
			user := appEvent.AsWorkUpdatedUser0

			util.LogWithRed("===========================================WorkUpdated: ", workId)
			version, _ := chain.GetVersion(m.ChainClient, workId)
			appIns := chain.App{
				Client: m.ChainClient,
			}
			app, _ := appIns.GetApp(user[:], workId.Id)
			err = m.UpdateApp(&ctx, user[:], workId, app, version)
			util.LogWithRed("===========================================CreateOrUpdatePod error: ", err)
		}
	}
	return err
}
