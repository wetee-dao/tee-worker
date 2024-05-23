package mint

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base32"
	"encoding/hex"
	"html/template"
	"math/rand"
	"strings"

	"fmt"
	"time"

	"github.com/pkg/errors"
	chain "github.com/wetee-dao/go-sdk"
	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"wetee.app/worker/store"
	"wetee.app/worker/util"
)

// 获取容器的资源信息和日志
func (m *Minter) getMetricInfo(ctx context.Context, wid gtypes.WorkId, nameSpace, name string, stage uint64) ([]string, map[string][]int64, error) {
	podLogOpts := &corev1.PodLogOptions{
		SinceTime: &metav1.Time{
			Time: time.Now().Add(-6 * time.Second * time.Duration(stage)),
		},
	}

	// 如果是不是TASK类型，则获取c0容器的日志
	if !wid.Wtype.IsTASK {
		podLogOpts.Container = "c0"
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
	address := base32.HexEncoding.EncodeToString(user[:])
	address = strings.ReplaceAll(strings.ToLower(address), "=", "")
	return strings.TrimRight(address, "000000000000000000")
}

// Hex Address To Account
// 将hex地址转换为用户公钥
func HexStringToSpace(address string) string {
	address = strings.ReplaceAll(address, "0x", "")
	user, _ := hex.DecodeString(address)
	return AccountToSpace(user)
}

// Get Envs from Work
// 获取环境变量
func (m *Minter) BuildEnvs(workId gtypes.WorkId) ([]corev1.EnvVar, error) {
	settings, err := m.GetSettingsFromWork(workId, nil)
	if err != nil {
		return []corev1.EnvVar{}, errors.Wrap(err, "GetSettingsFromWork error")
	}

	return m.BuildEnvsFromSettings(workId, settings)
}

// Build Envs
// 获取配置文件
func (m *Minter) BuildEnvsFromSettings(workId gtypes.WorkId, settings []*gtypes.Env) ([]corev1.EnvVar, error) {
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

// WrapNodeService
// 包装环境变量
func (m *Minter) WrapEnvs(envs []corev1.EnvVar, nameSpace, name string, nodeSers *v1.Service) error {
	mdata := make(map[string]string)
	mdata["cluster_domain"] = m.HostDomain
	mdata["project_domain"] = nameSpace + ".svc.cluster.local"
	mdata["gen_ssl"] = strings.Join(util.GetSslRoot(), "|")
	for i, port := range nodeSers.Spec.Ports {
		if port.NodePort != 0 {
			mdata["ser_"+fmt.Sprint(i)+"_nodeport"] = fmt.Sprint(port.NodePort)
		}
	}

	for i, env := range envs {
		if strings.Contains(env.Value, "{{.") {
			v, err := renderTemplate(env.Value, mdata)
			if err != nil {
				return err
			}

			envs[i].Value = v
		}
	}

	return nil
}

// BuildCommand
// 构建启动命令
func (m *Minter) BuildCommand(cmd *gtypes.Command) []string {
	if cmd.IsNONE {
		return []string{}
	}
	if cmd.IsBASH {
		return []string{"bash", "-c", string(cmd.AsBASHField0)}
	}
	if cmd.IsSH {
		return []string{"/bin/sh", "-c", string(cmd.AsSHField0)}
	}
	if cmd.IsZSH {
		return []string{"/bin/zsh", "-c", string(cmd.AsZSHField0)}
	}
	return []string{}
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

// Get Container Port From Service
// 获取容器服务端口
func BuildContainerPortFormService(name string, services []gtypes.Service) []corev1.ContainerPort {
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
func (m *Minter) BuildServicePortFormService(name string, services []gtypes.Service) ([]corev1.ServicePort, []corev1.ServicePort) {
	nodePorts := []corev1.ServicePort{}
	headlessPorts := []corev1.ServicePort{}
	for i, ser := range services {
		var protocol corev1.Protocol
		var port uint16

		// 获取服务端口
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

		if ser.IsTcp || ser.IsUdp {
			if port != 0 {
				nodePorts = append(nodePorts, corev1.ServicePort{
					Name:       name + "-" + fmt.Sprint(i) + "-" + fmt.Sprint(port) + "-nodeport",
					Port:       int32(port),
					TargetPort: intstr.FromInt(int(port)),
					Protocol:   protocol,
				})
			} else {
				nodePort := m.randNodeport()
				if ser.IsTcp {
					services[i].AsTcpField0 = nodePort
				}
				if ser.IsUdp {
					services[i].AsUdpField0 = nodePort
				}
				nodePorts = append(nodePorts, corev1.ServicePort{
					Name:       name + "-" + fmt.Sprint(i) + "-" + fmt.Sprint(nodePort) + "-nodeport",
					Port:       int32(nodePort),
					TargetPort: intstr.FromInt(int(nodePort)),
					NodePort:   int32(nodePort),
					Protocol:   protocol,
				})
			}
		} else {
			headlessPorts = append(headlessPorts, corev1.ServicePort{
				Name:       name + "-" + fmt.Sprint(i) + "-" + fmt.Sprint(port) + "-headless",
				Port:       int32(port),
				TargetPort: intstr.FromInt(int(port)),
				Protocol:   protocol,
			})
		}
	}
	return nodePorts, headlessPorts
}

// 随机生成NodePort端口
func (m *Minter) randNodeport() uint16 {
	// 查询当前已分配的NodePort端口列表
	usedPorts := []int32{}
	services, err := m.K8sClient.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, svc := range services.Items {
		for _, port := range svc.Spec.Ports {
			if port.NodePort != 0 {
				usedPorts = append(usedPorts, port.NodePort)
			}
		}
	}

	var port uint16 = 0
	// 生成一个随机的NodePort端口
	for {
		randomPort := int32(rand.Intn(2767) + 30000)
		if !contains(usedPorts, randomPort) {
			port = uint16(randomPort)
			break
		}
	}

	return port
}

func renderTemplate(templateString string, data map[string]string) (string, error) {
	tmpl, err := template.New("myTemplate").Parse(templateString)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (m *Minter) buildPodContainer(
	ctx *context.Context,
	workId gtypes.WorkId,
	nameSpace, name string,
	cs []gtypes.Container,
	envs []*gtypes.Env,
) ([]v1.Container, error) {
	ty := ""
	if workId.Wtype.IsAPP {
		ty = "app"
	} else if workId.Wtype.IsGPU {
		ty = "gpu"
	}

	podContainers := make([]v1.Container, 0, len(cs))

	serviceSpace := m.K8sClient.CoreV1().Services(nameSpace)
	nodeports, projectPorts := []corev1.ServicePort{}, []corev1.ServicePort{}

	// 计算所有的服务端口
	for _, container := range cs {
		cnodeports, cprojectPorts := m.BuildServicePortFormService(name, container.Port)

		// 创建对外端口
		nodeports = append(nodeports, cnodeports...)
		projectPorts = append(projectPorts, cprojectPorts...)
	}

	// 添加机密认证服务
	nodeports = append(nodeports, v1.ServicePort{
		Name:       name + "-8888",
		Protocol:   "TCP",
		Port:       8888,
		TargetPort: intstr.FromInt(8888),
	})

	// 创建对外服务
	aservice := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name + "-expose",
			Labels: map[string]string{"service": name},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{ty: name},
			Type:     "NodePort",
			Ports:    nodeports,
		},
	}

	nodeSers, err := serviceSpace.Create(*ctx, &aservice, metav1.CreateOptions{})
	fmt.Println("================================================= Create service", err)
	if err != nil {
		return nil, err
	}

	// 创建项目内端口
	pservice := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"service": name},
		},
		Spec: v1.ServiceSpec{
			Selector:  map[string]string{ty: name},
			ClusterIP: "None",
			Ports:     projectPorts,
		},
	}
	_, err = serviceSpace.Create(*ctx, &pservice, metav1.CreateOptions{})
	fmt.Println("================================================= Create project service", err)
	if err != nil {
		return nil, err
	}

	// 构建容器
	for i, container := range cs {
		ports := BuildContainerPortFormService(name, container.Port)
		if i == 0 {
			ports = append(ports, v1.ContainerPort{
				Name:          "port0",
				ContainerPort: int32(8888),
				Protocol:      "TCP",
			})
		}

		cnevs, err := m.BuildEnvsFromSettings(workId, filterEnvs(envs, uint16(i)))
		if err != nil {
			return nil, err
		}

		err = m.WrapEnvs(cnevs, nameSpace, name, nodeSers)
		fmt.Println("================================================= Create WrapEnvs", err)
		if err != nil {
			return nil, err
		}

		podContainers = append(podContainers, v1.Container{
			Name:    "c" + fmt.Sprint(i),
			Image:   string(container.Image),
			Ports:   ports,
			Env:     cnevs,
			Command: m.BuildCommand(&container.Command),
			Resources: v1.ResourceRequirements{
				Limits: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse(fmt.Sprint(container.Cr.Cpu) + "m"),
					v1.ResourceMemory: resource.MustParse(fmt.Sprint(container.Cr.Mem) + "M"),
				},
				Requests: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse(fmt.Sprint(container.Cr.Cpu) + "m"),
					v1.ResourceMemory: resource.MustParse(fmt.Sprint(container.Cr.Mem) + "M"),
				},
			},
		})
	}

	return podContainers, nil
}

func filterEnvs(envs []*gtypes.Env, index uint16) []*gtypes.Env {
	var fenvs []*gtypes.Env
	for i, env := range envs {
		if env.Index == index {
			fenvs = append(fenvs, envs[i])
		}
	}
	return fenvs
}

func contains(s []int32, e int32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
