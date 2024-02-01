package mint

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/wetee-dao/go-sdk/gen/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Minter) getMetricInfo(ctx context.Context, wid types.WorkId, nameSpace, name string, stage uint64) ([]string, map[string][]int64, error) {
	podLogOpts := &corev1.PodLogOptions{
		SinceTime: &metav1.Time{
			Time: time.Now().Add(-6 * time.Second * time.Duration(stage)),
		},
	}
	clientset := m.K8sClient
	metricsClient := m.MetricsClient

	// 获取Pod的logs
	// Get the logs of the Pod
	req := clientset.CoreV1().Pods(nameSpace).GetLogs(name, podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer podLogs.Close()

	// Read the logs line by line
	// 读取到的logs是 []string
	logs := []string{}
	scanner := bufio.NewScanner(podLogs)
	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("failed to read log line: %v", err)
	}

	fmt.Println("================================================logs: ", logs)
	var mem map[string][]int64 = map[string][]int64{}
	if wid.Wtype.IsAPP {
		// 获取Pod的内存使用情况
		// Gets the memory usage of the Pod
		podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(nameSpace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, nil, err
		}

		// 遍历Pod的容器，获取内存使用情况
		// Walk through the Pod container to get memory usage
		var mem map[string][]int64 = map[string][]int64{}
		for _, container := range podMetrics.Containers {
			fmt.Println("Pod ", podMetrics.Name, " CPU使用情况: ", container.Usage.Cpu().Value())
			fmt.Println("Pod ", podMetrics.Name, " 内存使用情况: ", container.Usage.Memory().Value()/1024/1024)

			mem[container.Name] = []int64{container.Usage.Cpu().Value(), container.Usage.Memory().Value() / 1024 / 1024, 0}
		}
	} else {
		mem["d"] = []int64{0, 0, 0}
	}

	return logs, mem, nil
}

// Account To Hex Address
// 将用户公钥转换为hex地址
func AccountToAddress(user []byte) string {
	address := hex.EncodeToString(user[:])
	saddress := address[1:] //去掉前面的 0x
	return saddress
}
