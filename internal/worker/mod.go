package worker

import (
	"bufio"
	"context"
	"fmt"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/author"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func WorkerInit(mgr manager.Manager) {
	ctx := context.Background()

	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		fmt.Println("error in opening stream" + err.Error())
		return
	}
	podLogOpts := &corev1.PodLogOptions{
		Container: "kube-rbac-proxy",
		SinceTime: &metav1.Time{
			Time: time.Now().Add(-1 * time.Minute),
		},
	}
	req := clientset.CoreV1().Pods("worker-system").GetLogs("worker-controller-manager-66dfb44b9b-9scl8", podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		fmt.Println("error in opening stream " + err.Error())
		return
	}
	defer podLogs.Close()

	// Read the logs line by line
	logs := ""
	scanner := bufio.NewScanner(podLogs)
	for scanner.Scan() {
		logs += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("failed to read log line: %v", err)
	}

	fmt.Println("logs: ================================================")
	fmt.Println(logs)
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<,")

	// 创建Metrics Client
	metricsClient, err := versioned.NewForConfig(mgr.GetConfig())
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取Pod的内存使用情况
	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses("worker-system").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(podMetrics.Items); i++ {
		pod := podMetrics.Items[i]
		// 遍历Pod的容器，获取内存使用情况
		for _, container := range pod.Containers {
			fmt.Printf("Pod %s CPU使用情况: %s \n", pod.Name, container.Usage.Cpu())
			fmt.Printf("Pod %s 内存使用情况: %s \n", pod.Name, container.Usage.Memory())
		}
	}
}

func ChainConnect() {
	// Query the system events and extract information from them. This example runs until exited via Ctrl-C

	// Create our API with a default connection to the local node
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	fmt.Println(genesisHash.Hex())

	from := signature.TestKeyringPairAlice

	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")

	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompactFromUInt(1000000000000000000))

	ext := types.NewExtrinsic(c)

	era := types.ExtrinsicEra{IsMortalEra: false}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()

	key, err := types.CreateStorageKey(meta, "System", "Account", from.PublicKey)

	var sub *author.ExtrinsicStatusSubscription

	var accountInfo types.AccountInfo
	_, err = api.RPC.State.GetStorageLatest(key, &accountInfo)

	nonce := uint32(accountInfo.Nonce)

	fmt.Println("accountInfo.Data.Free===>", accountInfo.Data.Free)
	o := types.SignatureOptions{
		// BlockHash:   blockHash,
		BlockHash:          genesisHash, // BlockHash needs to == GenesisHash if era is immortal. // TODO: add an error?
		Era:                era,
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	err = ext.Sign(from, o)

	sub, err = api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		panic(err)
	}

	defer sub.Unsubscribe()
	timeout := time.After(20 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			fmt.Printf("%#v\n", status)

			if status.IsInBlock {
				fmt.Println("IsInBlock")
				// return
			}
			if status.IsFinalized {
				fmt.Println("IsFinalized")
				var accountInfo types.AccountInfo
				api.RPC.State.GetStorageLatest(key, &accountInfo)
				fmt.Println("accountInfo.Data.Free ===> ", accountInfo.Data.Free)
				return
			}
		case <-timeout:
			fmt.Println("timeout")
			return
		}
	}
}
