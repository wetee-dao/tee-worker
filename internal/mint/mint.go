package mint

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/manager"

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

	side *sidechain.SideChain
	// preRecerve is the channel to receive SendEncryptedSecretRequest
	preRecerve map[string]chan any
}

// InitCluster
// 初始化矿工
func InitCluster(mgr manager.Manager, privateKey *model.PrivKey, side *sidechain.SideChain) error {
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
		side:          side,
		preRecerve:    make(map[string]chan any),
	}
	MinterIns.PrivateKey = privateKey

	return err
}

// start mint
// 开始挖矿
func (m *Minter) StartMint() {
	signer := m.PrivateKey.ToSigner()

	// 等待集群开启
	// Waiting for cluster start
	for {
		// 获取clusterId
		cluster, err := chains.MainChain.GetMintWorker(signer.AccountID())
		if err != nil {
			fmt.Println("ClusterId => clusterId not found, mint not started ", err)
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
		err = store.SetClusterId(cluster.Id)
		if err != nil {
			util.LogError("worker.SetClusterId", err)
			time.Sleep(time.Second * 10)
			continue
		}
		break
	}

	clusterId, _ := store.GetClusterId()

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
			util.LogError("VersionContracts", err)
			sleepFrom(start, time.Second*6)
			continue
		}

		// 计算新增、更新、未改变的合约
		// Calculate pod version
		added, updated, unchanged, deleted, err := CalcPodVersionFromCache(podVersions)
		if err != nil {
			util.LogError("SetPods", err)
			sleepFrom(start, time.Second*6)
			continue
		}

		// 删除过期的合约
		// Delete contracts
		for _, p := range deleted {
			err := m.StopPod(p)
			if err != nil && !strings.Contains(err.Error(), "not found") {
				util.LogError("DelPod "+fmt.Sprint(p.PodId), err)
				continue
			}

			err = store.DelPod(p)
			if err != nil {
				util.LogWithRed("DelPod", err)
			}
		}

		// 等待部署的程序
		// Pod list from chain to deploy
		addAndUpdated := append(added, updated...)
		ids := make([]uint64, 0, len(addAndUpdated))
		for _, p := range addAndUpdated {
			ids = append(ids, p.PodId)
		}
		deployList, err := chains.MainChain.GetPodsByIds(ids)
		if err != nil {
			util.LogError("GetPodsByIds", err)
			sleepFrom(start, time.Second*6)
			continue
		}
		if len(podVersions) > 0 {
			fmt.Println()
			util.LogWithGreen("BLOCK  ", "-", head)
			util.LogWithGray("POD ALL", "-", len(podVersions))
			util.LogWithGray("TODO   ", "-", len(deployList))
		}

		// 触发TEE调用
		// Trigger TEE calls
		// m.trigger(cs, clusterId, uint64(head.Number))

		// 部署程序
		// Deploy POD
		ctx := context.Background()
		for _, p := range deployList {
			// 跳过20个块内已经部署的程序
			// Skip programs deployed within 20 blocks
			if p.DeploySkipUtil >= uint32(head) {
				continue
			}

			if p.Ptype.CPU != nil {
				_, err = m.DeployOrUpdateAPP(&ctx, p)
				if err != nil {
					util.LogError("DeployOrUpdateAPP", err)
					continue
				}
			} else if p.Ptype.SCRIPT != nil {
				_, err = m.DeployOrUpdateTASK(&ctx, p)
				if err != nil {
					util.LogError("DeployOrUpdateTASK", err)
					continue
				}
			} else if p.Ptype.GPU != nil {
				_, err = m.DeployOrUpdateGPU(&ctx, p)
				if err != nil {
					util.LogError("DeployOrUpdateGPU", err)
					continue
				}
			}

			// 20个区块内不重复处理
			// 20 blocks do not repeat processing
			p.DeploySkipUtil += 20
			err = store.SetPod(p)
			if err != nil {
				util.LogWithRed("SetPod", err)
			}
		}

		// 获取收费周期
		// Get the charge cycle
		stage, err := chains.MainChain.GetMintInterval()
		if err != nil {
			util.LogError("GetMintInterval", err)
			sleepFrom(start, time.Second*6)
			continue
		}

		// 获取程序运行费用
		// Mint POD
		mintCalls := make([]*model.TeeCall, 0, 20)
		for _, pv := range unchanged {
			p, err := store.GetPod(pv.PodId)
			if err != nil {
				util.LogError("GetPod", err)
				continue
			}

			// 更新最后挖矿区块
			p.LastMintBlockNumber = pv.LastMint

			// 跳过200个块以内已经部署的POD
			if p.SkipUtil >= uint32(head) {
				continue
			}

			if pv.Status != 1 {
				continue
			}
			m.checkPodStatus(&ctx, *p)

			if p.Ptype.CPU != nil {
				call, err := m.MintAPP(&ctx, *p, stage, uint32(head))
				if err != nil {
					util.LogError("MintAPP", err)
					continue
				}
				if call != nil {
					mintCalls = append(mintCalls, call)
				}
			} else if p.Ptype.SCRIPT != nil {
				call, err := m.MintTASK(&ctx, *p, stage, uint32(head))
				if err != nil {
					util.LogError("MintTASK", err)
					continue
				}
				if call != nil {
					mintCalls = append(mintCalls, call)
				}
			} else if p.Ptype.GPU != nil {
				call, err := m.MintGPU(&ctx, *p, stage, uint32(head))
				if err != nil {
					util.LogError("MintGPU", err)
					continue
				}
				if call != nil {
					mintCalls = append(mintCalls, call)
				}
			}

			store.SetPodSkipUtil(p.PodId, uint32(head)+10, pv.LastMint)
		}

		// 读取待同步到区块链的调用
		// query call for sync to chain
		calls, keys, _ := store.GetPendingCall()
		mintCalls = append(mintCalls, calls...)

		// 批量提交调用到区块链
		// sync tx to chain
		if len(mintCalls) > 0 {
			tx := &model.SysCall{
				Payload: &model.SysCall_HubCall{
					HubCall: &model.HubCall{Call: mintCalls},
				},
			}

			err = m.side.SubmitCallFromNode(tx)
			if err != nil && !strings.Contains(err.Error(), "tx already exists") && err.Error() != "Tx already received from peer" {
				util.LogError("SubmitTx", err)
			} else {
				err := store.DeletePendingCalls(keys)
				if err != nil {
					util.LogError("DeletePendingCalls", err)
				}
			}
		}

		sleepFrom(start, time.Second*6)
	}
}

// 测试已经部署的程序的状态
// check pod status
func (m *Minter) checkPodStatus(ctx *context.Context, pod model.Pod) {
	// get namespace name
	name := GetPodName(pod.PodId)
	address := AccountToSpace(pod.Owner[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(address)

	// get deployment
	deployment, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			util.LogWithRed("====== check pod status err", err)
			return
		}
	}

	// check if deployment exists
	if deployment.Name != "" {
		return
	}

	// deploy app
	if pod.Ptype.CPU != nil {
		m.DeployOrUpdateAPP(ctx, pod)
	} else if pod.Ptype.GPU != nil {
		m.DeployOrUpdateGPU(ctx, pod)
	}
}

func CalcPodVersionFromCache(newPod []model.PodVersion) (added, updated, unchanged []model.PodVersion, deleted []model.Pod, err error) {
	// get data from cache
	oldPod, err := store.GetPods()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	oldMap := make(map[uint64]model.Pod)
	newMap := make(map[uint64]model.PodVersion)
	for _, v := range oldPod {
		oldMap[v.PodId] = *v
	}
	for _, v := range newPod {
		newMap[v.PodId] = v
	}

	// find add and update
	for id, newVal := range newMap {
		oldVal, exists := oldMap[id]
		if !exists {
			// add
			added = append(added, newVal)
		} else {
			if oldVal.Version != newVal.Version {
				// update
				updated = append(updated, newVal)
			} else {
				// unchanged update last mint
				unchanged = append(unchanged, newVal)
				store.SetPodLastMint(id, newVal.LastMint)
			}
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
