package mint

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"wetee.app/worker/db"
	"wetee.app/worker/mint/chain"
	"wetee.app/worker/mint/chain/gen/system"
	"wetee.app/worker/mint/chain/gen/types"
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
	Signer, err = chain.GetMintKey()
	return err
}

// start mint
// 开始挖矿
func StartMint() {
	fmt.Println("MintKey => ", Signer.Address)
	client := MinterIns.ChainClient
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
		fmt.Println("ClusterId => ", clusterId)
		db.SetClusterId(clusterId)

		break
	}

	clusterId, _ := db.GetClusterId()

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
					fmt.Println("===========================================WorkRuning ID: ", workId)
					err = CreateOrUpdatePod(user[:], workId, blockHash.Hex())
					fmt.Println("===========================================CreateOrUpdatePod error: ", err)
				}
			}
			if e.IsWeteeApp {
				appEvent := e.AsWeteeAppField0
				if appEvent.IsWorkStopped {
					workId := appEvent.AsWorkStoppedWorkId1
					fmt.Println("===========================================WorkStopped", workId)
					err := StopPod(workId)
					fmt.Println("===========================================StopPod error: ", err)
				}
				if appEvent.IsWorkUpdated {
					workId := appEvent.AsWorkUpdatedWorkId1
					user := appEvent.AsWorkUpdatedUser0
					fmt.Println("===========================================WorkUpdated ID: ", workId)
					err = CreateOrUpdatePod(user[:], workId, blockHash.Hex())
					fmt.Println("===========================================CreateOrUpdatePod error: ", err)
				}
				// e.AsWeteeAppField0.IsClusterCreated
			}
			// if e.IsWeteeTask {
			// fmt.Println("e.AsWeteeTaskField0")
			// }
		}

		// 获取合约列表
		cs, err := worker.GetClusterContracts(clusterId, &blockHash)
		fmt.Println("GetClusterContracts", err)
		fmt.Println("GetClusterContracts", cs)

		// 校对合约状态
		for _, c := range cs {
			err := checkPodStatus(c, blockHash.Hex())
			fmt.Println("checkPodStatus", err)
		}
	}
}

func checkPodStatus(state chain.ContractStateWrap, blockHash string) error {
	ctx := context.Background()
	k8s := MinterIns.K8sClient.CoreV1()
	address := hex.EncodeToString(state.ContractState.User[:])
	nameSpace := k8s.Pods(address[1:])
	workID := state.ContractState.WorkId
	name := GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	pod, err := nameSpace.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if pod.ObjectMeta.ResourceVersion == state.BlockHash {
		return nil
	}

	return CreateOrUpdatePod(state.ContractState.User[:], workID, blockHash)
}

func CreateOrUpdatePod(user []byte, workID types.WorkId, blockHash string) error {
	address := hex.EncodeToString(user[:])
	saddress := address[1:] //去掉前面的 0x
	ctx := context.Background()
	errc := checkNameSpace(ctx, saddress)
	if errc != nil {
		return errc
	}

	appIns := chain.App{
		Client: MinterIns.ChainClient,
		Signer: Signer,
	}
	app, err := appIns.GetApp(user[:], workID.Id)
	if err != nil {
		return err
	}

	k8s := MinterIns.K8sClient.CoreV1()
	nameSpace := k8s.Pods(saddress)
	name := GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	_, err = nameSpace.Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		existingPod, err := nameSpace.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		existingPod.ObjectMeta.ResourceVersion = blockHash
		existingPod.Spec.Containers[0].Image = string(app.Image)
		existingPod.Spec.Containers[0].Ports[0].ContainerPort = int32(app.Port[0])
		_, err = nameSpace.Update(ctx, existingPod, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
	} else {
		pod := &v1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "App",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            name,
				ResourceVersion: blockHash,
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "c1",
						Image: string(app.Image),
						Ports: []v1.ContainerPort{
							{
								Name:          string(app.Name) + "0",
								ContainerPort: 80,
								Protocol:      "TCP",
							},
						},
					},
				},
			},
		}
		_, err = nameSpace.Create(ctx, pod, metav1.CreateOptions{})
		fmt.Println("================================================= Create", err)
	}

	return err
}

func StopPod(workID types.WorkId) error {
	ctx := context.Background()
	user, err := chain.GetAccount(MinterIns.ChainClient, workID)
	if err != nil {
		return err
	}
	address := hex.EncodeToString(user[:])
	saddress := address[1:] //去掉前面的 0x

	k8s := MinterIns.K8sClient.CoreV1()
	nameSpace := k8s.Pods(saddress)
	name := GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)
	return nameSpace.Delete(ctx, name, metav1.DeleteOptions{})
}

func checkNameSpace(ctx context.Context, address string) error {
	k8s := MinterIns.K8sClient.CoreV1()
	nameSpaces, err := k8s.Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	var found bool = false
	for _, namespace := range nameSpaces.Items {
		if namespace.Name == address {
			found = true
			break
		}
	}
	if !found {
		_, err = k8s.Namespaces().Create(ctx, &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: address,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetWorkTypeStr(work types.WorkId) string {
	if work.Wtype.IsAPP {
		return "app"
	}

	if work.Wtype.IsTASK {
		return "task"
	}

	return "unknown"
}

// func getPodInfo() {
// podLogOpts := &corev1.PodLogOptions{
// 	Container: "worker",
// 	SinceTime: &metav1.Time{
// 		Time: time.Now().Add(-1 * time.Minute),
// 	},
// }
// req := clientset.CoreV1().Pods("default").GetLogs("testnginx", podLogOpts)
// podLogs, err := req.Stream(ctx)
// if err != nil {
// 	fmt.Println("error in opening stream " + err.Error())
// 	return
// }
// defer podLogs.Close()

// // Read the logs line by line
// logs := ""
// scanner := bufio.NewScanner(podLogs)
// for scanner.Scan() {
// 	logs += scanner.Text() + "\n"
// }
// if err := scanner.Err(); err != nil {
// 	fmt.Printf("failed to read log line: %v", err)
// }

// fmt.Println("logs: ================================================")
// fmt.Println(logs)
// fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

// 获取Pod的内存使用情况
// podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses("default").List(ctx, metav1.ListOptions{})
// if err != nil {
// 	fmt.Println(err)
// 	return
// }

// for i := 0; i < len(podMetrics.Items); i++ {
// 	pod := podMetrics.Items[i]
// 	// 遍历Pod的容器，获取内存使用情况
// 	for _, container := range pod.Containers {
// 		fmt.Printf("Pod %s CPU使用情况: %s \n", pod.Name, container.Usage.Cpu())
// 		fmt.Printf("Pod %s 内存使用情况: %s \n", pod.Name, container.Usage.Memory())
// 	}
// }
// }
