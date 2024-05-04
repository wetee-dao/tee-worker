package mint

import (
	"bufio"
	"context"
	"encoding/base32"
	"encoding/hex"
	"strings"

	"fmt"
	"time"

	"github.com/pkg/errors"
	chain "github.com/wetee-dao/go-sdk"
	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

	// fmt.Println("================================================logs: ", logs)
	var use map[string][]int64 = map[string][]int64{}

	// 获取Pod的内存使用情况
	// Gets the memory usage of the Pod
	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(nameSpace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if wid.Wtype.IsAPP {
			return nil, nil, err
		} else {
			use["d"] = []int64{0, 0, 0}
		}
	} else {
		// 遍历Pod的容器，获取内存使用情况
		// Walk through the Pod container to get memory usage
		for _, container := range podMetrics.Containers {
			fmt.Println("Pod ", podMetrics.Name, " CPU使用情况: ", container.Usage.Cpu().MilliValue(), " M")
			fmt.Println("Pod ", podMetrics.Name, " 内存使用情况: ", container.Usage.Memory().Value()/(1024*1024), " MB")

			use[container.Name] = []int64{container.Usage.Cpu().MilliValue(), container.Usage.Memory().Value() / (1024 * 1024), 0}
		}
	}

	return logs, use, nil
}

// Account To Hex Address
// 将用户公钥转换为hex地址
func AccountToSpace(user []byte) string {
	// address := hex.EncodeToString(user[:])
	address := base32.HexEncoding.EncodeToString(user[:])
	// address := base64.StdEncoding(user[:])
	// saddress := address[1:] //去掉前面的 0x54c
	return strings.ReplaceAll(strings.ToLower(address), "=", "")
}

func HexStringToSpace(address string) string {
	address = strings.ReplaceAll(address, "0x", "")
	user, _ := hex.DecodeString(address)
	return AccountToSpace(user)
}

// Get Envs from Work
// 获取环境变量
func (m *Minter) GetEnvs(workId gtypes.WorkId) ([]corev1.EnvVar, error) {

	settings, err := m.GetSettingsFromWork(workId, nil)
	if err != nil {
		return []corev1.EnvVar{}, errors.Wrap(err, "GetSettingsFromWork error")
	}

	return m.GetEnvsFromSettings(workId, settings)
}

// 获取配置文件
func (m *Minter) GetEnvsFromSettings(workId gtypes.WorkId, settings []*gtypes.Env) ([]corev1.EnvVar, error) {
	// 用于应用联系控制面板的凭证
	wid, err := store.SealAppID(workId)
	if err != nil {
		return []corev1.EnvVar{}, err
	}

	envs := []corev1.EnvVar{
		{Name: "APPID", Value: wid},
		{Name: "IN_TEE", Value: string("1")},
	}

	for _, setting := range settings {
		// TODO add file
		if setting.K.IsFile {
			continue
		}
		envs = append(envs, corev1.EnvVar{
			Name:  string(setting.K.AsEnvField0),
			Value: string(setting.V),
		})
	}

	return envs, nil
}

// Get Container Port From Service
// 获取容器服务端口
func GetContainerPortFormService(name string, services []gtypes.Service) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{}
	for _, ser := range services {
		protocol := corev1.ProtocolTCP
		port := ser.AsTcpField0
		if ser.IsProjectUdp {
			protocol = corev1.ProtocolUDP
			port = ser.AsProjectUdpField0
		} else if ser.IsProjectTcp {
			protocol = corev1.ProtocolTCP
			port = ser.AsProjectTcpField0
		} else if ser.IsTcp {
			protocol = corev1.ProtocolSCTP
			port = ser.AsTcpField0
		} else if ser.IsUdp {
			protocol = corev1.ProtocolUDP
			port = ser.AsUdpField0
		}
		ports = append(ports, corev1.ContainerPort{
			Name:          name + "-" + fmt.Sprint(port),
			ContainerPort: int32(port),
			Protocol:      protocol,
		})
	}
	return ports
}

// Get Service Port From Service
// 获取对外服务端口
func GetServicePortFormService(name string, services []gtypes.Service) []corev1.ServicePort {
	ports := []corev1.ServicePort{}
	for _, ser := range services {
		protocol := corev1.ProtocolTCP
		port := ser.AsTcpField0
		if ser.IsProjectUdp {
			protocol = corev1.ProtocolUDP
			port = ser.AsProjectUdpField0
		} else if ser.IsProjectTcp {
			protocol = corev1.ProtocolTCP
			port = ser.AsProjectTcpField0
		} else if ser.IsTcp {
			protocol = corev1.ProtocolTCP
			port = ser.AsTcpField0
		} else if ser.IsUdp {
			protocol = corev1.ProtocolUDP
			port = ser.AsUdpField0
		}

		// if ser.IsProjectTcp || ser.IsProjectUdp {
		// 	ports = append(ports, corev1.ServicePort{
		// 		Name:       fmt.Sprint(port),
		// 		Port:       int32(port),
		// 		TargetPort: intstr.FromInt(int(port)),
		// 		Protocol:   protocol,
		// 	})
		// }

		if ser.IsTcp || ser.IsUdp {
			ports = append(ports, corev1.ServicePort{
				Name:       name + "-" + fmt.Sprint(port) + "-nodeport",
				Port:       int32(port),
				TargetPort: intstr.FromInt(int(port)),
				Protocol:   protocol,
			})
		}
	}
	return ports
}

// StopApp
// 停止应用
func (m *Minter) StopApp(workId gtypes.WorkId, space string) error {
	ctx := context.Background()

	if space == "" {
		user, err := chain.GetAccount(m.ChainClient, workId)
		if err != nil {
			return err
		}
		space = AccountToSpace(user[:])
	}

	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)
	util.LogWithRed("StopApp: ", name)

	ServiceSpace := m.K8sClient.CoreV1().Services(space)
	list, err := ServiceSpace.List(ctx, metav1.ListOptions{
		LabelSelector: "service=" + name,
	})
	if err != nil {
		return err
	}
	for _, item := range list.Items {
		err := ServiceSpace.Delete(ctx, item.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	if workId.Wtype.IsAPP || workId.Wtype.IsGPU {
		nameSpace := m.K8sClient.AppsV1().Deployments(space)

		return nameSpace.Delete(ctx, name, metav1.DeleteOptions{})
	}

	nameSpace := m.K8sClient.CoreV1().Pods(space)
	return nameSpace.Delete(ctx, name, metav1.DeleteOptions{})
}
