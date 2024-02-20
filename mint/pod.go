package mint

import (
	"bufio"
	"context"
	"encoding/hex"

	"fmt"
	"time"

	"github.com/pkg/errors"
	chain "github.com/wetee-dao/go-sdk"
	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/store"
	"wetee.app/worker/util"
)

func (m *Minter) getMetricInfo(ctx context.Context, wid gtypes.WorkId, nameSpace, name string, stage uint64) ([]string, map[string][]int64, error) {
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

	// 获取Pod的内存使用情况
	// Gets the memory usage of the Pod
	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(nameSpace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if wid.Wtype.IsAPP {
			return nil, nil, err
		} else {
			mem["d"] = []int64{0, 0, 0}
		}
	} else {
		// 遍历Pod的容器，获取内存使用情况
		// Walk through the Pod container to get memory usage
		for _, container := range podMetrics.Containers {
			fmt.Println("Pod ", podMetrics.Name, " CPU使用情况: ", container.Usage.Cpu().MilliValue(), " M")
			fmt.Println("Pod ", podMetrics.Name, " 内存使用情况: ", container.Usage.Memory().Value()/(1024*1024), " MB")

			mem[container.Name] = []int64{container.Usage.Cpu().MilliValue(), container.Usage.Memory().Value() / (1024 * 1024), 0}
		}
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

func (m *Minter) GetEnvs(workId gtypes.WorkId) ([]corev1.EnvVar, error) {
	// 用于应用联系控制面板的凭证
	wid, err := store.SealAppID(workId)
	if err != nil {
		return []corev1.EnvVar{}, err
	}

	envs := []corev1.EnvVar{
		{Name: "APPID", Value: wid},
		{Name: "IN_TEE", Value: string("1")},
	}

	settings, err := m.GetSettingsFromWork(workId, nil)
	if err != nil {
		return envs, errors.Wrap(err, "GetSettingsFromWork error")
	}

	for _, setting := range settings {
		envs = append(envs, corev1.EnvVar{
			Name:  string(setting.K),
			Value: string(setting.V),
		})
	}

	return envs, nil
}

func (m *Minter) GetEnvsFromSettings(workId gtypes.WorkId, settings []*gtypes.AppSetting) ([]corev1.EnvVar, error) {
	// 用于应用联系控制面板的凭证
	wid, err := store.SealAppID(workId)
	fmt.Println(wid)
	if err != nil {
		return []corev1.EnvVar{}, err
	}

	envs := []corev1.EnvVar{
		{Name: "APPID", Value: wid},
		{Name: "IN_TEE", Value: string("1")},
	}

	for _, setting := range settings {
		envs = append(envs, corev1.EnvVar{
			Name:  string(setting.K),
			Value: string(setting.V),
		})
	}

	return envs, nil
}

// StopApp
// 停止应用
func (m *Minter) StopApp(workId gtypes.WorkId) error {
	ctx := context.Background()
	user, err := chain.GetAccount(m.ChainClient, workId)
	if err != nil {
		return err
	}

	saddress := AccountToAddress(user[:])
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)
	// util.LogWithRed("===========================================StopPod", workId, " ", name)

	if workId.Wtype.IsAPP {
		nameSpace := m.K8sClient.AppsV1().Deployments(saddress)

		return nameSpace.Delete(ctx, name, metav1.DeleteOptions{})
	}

	nameSpace := m.K8sClient.CoreV1().Pods(saddress)
	return nameSpace.Delete(ctx, name, metav1.DeleteOptions{})
}
