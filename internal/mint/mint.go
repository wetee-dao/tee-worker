package mint

import (
	"fmt"
	"sync"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"wetee.app/worker/internal/mint/chain"
	"wetee.app/worker/util"
)

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

func InitMint(mgr manager.Manager) error {
	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}

	// 创建Metrics Client
	metricsClient, err := versioned.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}

	client, _ := chain.ClientInit()

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

func StartMint() {

	fmt.Println("MintKey => ", Signer.Address)
	client := MinterIns.ChainClient

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
		break
	}

	// 触发轮训事件
	sub, err := client.Api.RPC.Chain.SubscribeNewHeads()
	// 失败后等待10秒重新尝试
	if err != nil {
		util.LogWithRed("SubscribeNewHeads", err)
		time.Sleep(time.Second * 10)
		goto mintStart
	}
	defer sub.Unsubscribe()

	count := 0

	for {
		head := <-sub.Chan()
		fmt.Printf("Chain is at block: #%v\n", head.Number)
		count++
		// 执行任务检查
		// 执行
		if count == 10 {
			sub.Unsubscribe()
			break
		}
	}
}

func getPodInfo() {
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
}
