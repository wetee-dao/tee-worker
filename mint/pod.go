package mint

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getMetricInfo(ctx context.Context, nameSpace, name string, stage uint64) ([]string, map[string][]int64, error) {
	podLogOpts := &corev1.PodLogOptions{
		SinceTime: &metav1.Time{
			Time: time.Now().Add(-6 * time.Second * time.Duration(stage)),
		},
	}
	clientset := MinterIns.K8sClient
	metricsClient := MinterIns.MetricsClient

	req := clientset.CoreV1().Pods(nameSpace).GetLogs(name, podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer podLogs.Close()

	// Read the logs line by line
	logs := []string{}
	scanner := bufio.NewScanner(podLogs)
	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("failed to read log line: %v", err)
	}

	fmt.Println("logs: ================================================")
	fmt.Println(logs)
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

	// 获取Pod的内存使用情况
	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(nameSpace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	// 遍历Pod的容器，获取内存使用情况
	var mem map[string][]int64 = map[string][]int64{}
	for _, container := range podMetrics.Containers {
		fmt.Println("Pod ", podMetrics.Name, " CPU使用情况: ", container.Usage.Cpu().Value())
		fmt.Println("Pod ", podMetrics.Name, " 内存使用情况: ", container.Usage.Memory().Value()/1024/1024)

		mem[container.Name] = []int64{container.Usage.Cpu().Value(), container.Usage.Memory().Value() / 1024 / 1024}
	}

	return logs, mem, nil
}

func AccountToAddress(user []byte) string {
	address := hex.EncodeToString(user[:])
	saddress := address[1:] //去掉前面的 0x
	return saddress
}
