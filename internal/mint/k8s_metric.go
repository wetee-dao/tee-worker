package mint

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/internal/util"
)

// 获取工作日志和硬件资源使用量
func (m *Minter) GetMetric(ctx *context.Context, nameSpace string, pod model.Pod, now time.Time, stage uint32, isCache bool) ([]string, map[string][]int64, error) {
	// 获取上次记录的时间
	name := GetPodName(pod.PodId)
	if isCache {
		name = name + "-cache"
	}

	// 获取上次记录的缓存时间
	from, err := store.GetCacheId(name)
	if err != nil {
		if isCache {
			stage = 9
		}
		from = now.Add(-6 * time.Second * time.Duration(stage)).Unix()
	}

	// 通过 K8s API 获取指定命名空间中的 Pod 列表
	clientset := m.K8sClient
	pods, err := clientset.CoreV1().Pods(nameSpace).List(*ctx, metav1.ListOptions{
		LabelSelector: "tee=" + GetPodName(pod.PodId),
	})
	if err != nil {
		return nil, nil, err
	}

	// 判断是否找到匹配的 Pod
	if len(pods.Items) == 0 {
		return nil, nil, errors.New("pods is empty")
	}

	// 获取指定 Pod 的日志和硬件资源使用量信息
	logs, crs, err := m.queryMetric(*ctx, pod, nameSpace, pods.Items[0].Name, from)
	if err != nil {
		util.LogError("getMetricInfo", err)
	}

	// 返回获取到的日志和硬件资源使用量
	return logs, crs, nil
}

// 获取容器的资源信息和日志
func (m *Minter) queryMetric(ctx context.Context, pod model.Pod, nameSpace, name string, form int64) ([]string, map[string][]int64, error) {
	podLogOpts := &v1.PodLogOptions{
		SinceTime: &metav1.Time{
			Time: time.Unix(form, 0),
		},
	}

	// 如果是不是TASK类型，则获取c0容器的日志
	if pod.Ptype.SCRIPT == nil {
		podLogOpts.Container = "c0"
	}

	// 获取Pod的logs
	// Get the logs of the Pod
	clientset := m.K8sClient
	req := clientset.CoreV1().Pods(nameSpace).GetLogs(name, podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer podLogs.Close()

	// 读取到的logs是 []string
	// Read the logs line by line
	logs := []string{}
	scanner := bufio.NewScanner(podLogs)
	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("failed to read log line: %v", err)
	}

	use := map[string][]int64{}

	// 获取Pod的内存使用情况
	// Gets the memory usage of the Pod
	metricsClient := m.MetricsClient
	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(nameSpace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if pod.Ptype.CPU != nil {
			return nil, nil, err
		} else {
			use["d"] = []int64{0, 0, 0}
		}
	}

	// 遍历Pod的容器，获取内存使用情况
	// Walk through the Pod container to get memory usage
	for i, container := range podMetrics.Containers {
		fmt.Println("["+fmt.Sprint(i)+"]CPU", "===>", container.Usage.Cpu().MilliValue(), " M")
		fmt.Println("["+fmt.Sprint(i)+"]MEM", "===>", container.Usage.Memory().Value()/(1024*1024), " MB")

		use[container.Name] = []int64{container.Usage.Cpu().MilliValue(), container.Usage.Memory().Value() / (1024 * 1024), 0}
	}

	return logs, use, nil
}
