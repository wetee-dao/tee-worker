package mint

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	chains "github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	sidechain "github.com/wetee-dao/tee-dsecret/side-chain"

	"wetee.app/worker/internal/store"
	"wetee.app/worker/internal/util"
)

var (
	MinterIns *Minter
)

// Minter
// 矿工
type Minter struct {
	K8sClient     *kubernetes.Clientset
	MetricsClient *versioned.Clientset
	PrivateKey    *model.PrivKey
	HostDomain    string
	mu            sync.RWMutex

	// preRecerve is the channel to receive SendEncryptedSecretRequest
	preRecerve map[string]chan any
}

// InitCluster
// 初始化矿工
func InitCluster(mgr manager.Manager, privateKey *model.PrivKey) error {
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

	MinterIns = &Minter{
		K8sClient:     clientset,
		MetricsClient: metricsClient,
		HostDomain:    "",
		preRecerve:    make(map[string]chan interface{}),
	}
	MinterIns.PrivateKey = privateKey

	return err
}

// start mint
// 开始挖矿
func (m *Minter) StartMint() {
	signer, _ := m.PrivateKey.ToSigner()
	fmt.Println("MintKey => ", signer.Address)

	// 等待集群开启
	// Waiting for cluster start
	for {
		// 获取clusterId
		cluster, err := chains.MainChain.GetMintWorker(signer.AccountID())
		if err != nil {
			fmt.Println("ClusterId => clusterId not found, mint not started")
			time.Sleep(time.Second * 10)
			continue
		}

		// _, err = m.UploadClusterProof()
		// if err != nil {
		// 	fmt.Println("worker.ClusterProofUpload => ", err)
		// 	time.Sleep(time.Second * 10)
		// 	continue
		// }

		if cluster.Ip.Domain.IsSome() {
			MinterIns.HostDomain = string(cluster.Ip.Domain.V)
		} else {
			MinterIns.HostDomain = ""
		}

		// 保存clusterId
		store.SetClusterId(cluster.Id)

		break
	}

	clusterId, _ := store.GetClusterId()
	fmt.Println("ClusterId => ", clusterId)

	for {
		start := time.Now()
		client := chains.MainChain.GetClient()
		head, err := client.GetBlockNumber()
		if err != nil {
			fmt.Println("GetBlockNumber => ", err)
			sleepFrom(start, time.Second*6)
			continue
		}

		// Get contract list
		podVersions, err := chains.MainChain.GetPodsVersionByWorker(clusterId)
		if err != nil {
			util.LogError("GetClusterContracts", err)
			sleepFrom(start, time.Second*6)
			continue
		}

		// Calculate pod version
		added, updated, deleted, err := CalcPodVersionFromCache(podVersions)
		if err != nil {
			util.LogError("SetPods", err)
			sleepFrom(start, time.Second*6)
			continue
		}

		// 删除过期的合约
		// Delete expired contracts
		for _, p := range deleted {
			err := m.StopApp(p)
			if err != nil && !strings.Contains(err.Error(), "not found") {
				util.LogError("DeleteRuning "+fmt.Sprint(p.PodId), err)
				continue
			}
			store.DelPod(p)
		}

		addAndUpdated := append(added, updated...)
		ids := make([]uint64, 0, len(addAndUpdated))
		for _, p := range addAndUpdated {
			ids = append(ids, p.PodId)
		}

		// 获取收费周期
		// Get the charge cycle
		// stage, err := worker.GetStage()
		// if err != nil {
		// 	util.LogError("GetStage", err)
		// 	continue
		// }
		var stage uint32 = 30

		util.LogWithGray("POD ALL", ">", len(podVersions))
		todoList, err := chains.MainChain.GetPodsByIds(ids)
		if err != nil {
			util.LogError("GetPodsByIds", err)
			sleepFrom(start, time.Second*6)
			continue
		}
		util.LogWithCyan("TODO   ", ">", len(todoList))

		// 触发TEE调用
		// Trigger TEE calls
		// m.trigger(cs, clusterId, uint64(head.Number))

		// 校对合约状态
		// Check contract status
		calls := make([]*model.IndexCall, 0, 20)
		for _, p := range todoList {
			ctx := context.Background()

			if p.Ptype.CPU != nil {
				// 如果是APP类型，检查Pod状态，检查是否需要上传工作证明
				// If it is APP type, check Pod status, check if it needs to upload work proof
				call, t, err := m.DoWithAppState(&ctx, p, stage, uint32(head))
				if err != nil {
					util.LogError("DoWithAppState", err)
					continue
				}
				if call != nil {
					bt, _ := codec.Encode(call)
					calls = append(calls, &model.IndexCall{
						Index: t,
						Call:  bt,
					})
				}
			} else if p.Ptype.SCRIPT != nil {
				// 如果是TASK类型，检查Pod状态，Pod如果执行完成，则上传日志和结果
				// If it is TASK type, check Pod status, Pod if it is executed, upload logs and results
				call, t, err := m.DoWithTaskState(&ctx, p, stage, uint32(head))
				if err != nil {
					util.LogError("DoWithTaskState", err)
					continue
				}
				if call != nil {
					bt, _ := codec.Encode(call)
					calls = append(calls, &model.IndexCall{
						Index: t,
						Call:  bt,
					})
				}
			} else if p.Ptype.GPU != nil {
				// 如果是GPU类型，检查Pod状态，检查是否需要上传工作证明
				call, t, err := m.DoWithGpuAppState(&ctx, p, stage, uint32(head))
				if err != nil {
					util.LogError("DoWithGpuAppState", err)
					continue
				}
				if call != nil {
					bt, _ := codec.Encode(call)
					calls = append(calls, &model.IndexCall{
						Index: t,
						Call:  bt,
					})
				}
			}

			store.SetPod(p)
		}

		if len(calls) > 0 {
			_, err := sidechain.SubmitTx(&model.Tx{
				Payload: &model.Tx_HubCall{
					HubCall: &model.HubCall{
						Call: calls,
					},
				},
			})
			fmt.Println("submit tx", err)
		}

		sleepFrom(start, time.Second*6)
	}
}

func CalcPodVersionFromCache(newArr []model.PodVersion) (added, updated []model.PodVersion, deleted []model.Pod, e error) {
	// get data from cache
	oldArr, err := store.GetPods()
	if err != nil {
		return nil, nil, nil, err
	}

	oldMap := make(map[uint64]model.Pod)
	newMap := make(map[uint64]model.PodVersion)

	for _, v := range oldArr {
		oldMap[v.PodId] = *v
	}
	for _, v := range newArr {
		newMap[v.PodId] = v
	}

	// find add and update
	for id, newVal := range newMap {
		oldVal, exists := oldMap[id]
		if !exists {
			// add
			added = append(added, newVal)
		} else if oldVal.Version != newVal.Version {
			// update
			updated = append(updated, newVal)
		}
	}

	// find deleted
	for id, oldVal := range oldMap {
		if _, exists := newMap[id]; !exists {
			deleted = append(deleted, oldVal)
		}
	}
	return
}

func sleepFrom(start time.Time, duration time.Duration) {
	now := time.Now()
	if now.Sub(start) >= duration {
		return
	}
	time.Sleep(duration - now.Sub(start))
}
