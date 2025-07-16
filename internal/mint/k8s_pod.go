package mint

import (
	"bufio"
	"bytes"
	"context"
	"html/template"
	"math/rand"

	"fmt"
	"time"

	"github.com/pkg/errors"
	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/internal/util"
)

// 获取容器的资源信息和日志
func (m *Minter) getMetricInfo(ctx context.Context, wid model.Pod, nameSpace, name string, form int64) ([]string, map[string][]int64, error) {
	podLogOpts := &v1.PodLogOptions{
		SinceTime: &metav1.Time{
			Time: time.Unix(form, 0),
		},
	}

	// 如果是不是TASK类型，则获取c0容器的日志
	if wid.Ptype.SCRIPT == nil {
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

	var use map[string][]int64 = map[string][]int64{}

	// 获取Pod的内存使用情况
	// Gets the memory usage of the Pod
	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(nameSpace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if wid.Ptype.CPU != nil {
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

// BuildCommand
// 构建启动命令
func (m *Minter) BuildCommand(cmd *model.Command) []string {
	if cmd.NONE != nil {
		return []string{}
	}
	if cmd.BASH != nil {
		return []string{"bash", "-c", string(*cmd.BASH)}
	}
	if cmd.SH != nil {
		return []string{"/bin/sh", "-c", string(*cmd.SH)}
	}
	if cmd.ZSH != nil {
		return []string{"/bin/zsh", "-c", string(*cmd.ZSH)}
	}
	return []string{}
}

// StopApp
// 停止应用
func (m *Minter) StopApp(p model.Pod) error {
	ctx := context.Background()
	space := AccountToSpace(p.Owner[:])
	if space == "" {
		// user, err := module.GetAccount(m.ChainClient, workId)
		// if err != nil {
		// 	return err
		// }
		// space = AccountToSpace(user[:])
	}

	name := GetPodName(p.PodId)
	util.LogError("StopApp: ", name)

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

	if p.Ptype.CPU != nil || p.Ptype.GPU != nil {
		nameSpace := m.K8sClient.AppsV1().Deployments(space)

		return nameSpace.Delete(ctx, name, metav1.DeleteOptions{})
	}

	nameSpace := m.K8sClient.CoreV1().Pods(space)
	return nameSpace.Delete(ctx, name, metav1.DeleteOptions{})
}

// Get Container Port From Service
// 获取容器服务端口
func BuildContainerPortFormService(name string, services []model.Service) []v1.ContainerPort {
	ports := []v1.ContainerPort{}
	for _, ser := range services {
		protocol := v1.ProtocolTCP
		var port uint16

		// 获取服务端口
		if ser.ProjectUdp != nil {
			protocol = v1.ProtocolUDP
			port = *ser.ProjectUdp
		} else if ser.ProjectTcp != nil {
			protocol = v1.ProtocolTCP
			port = *ser.ProjectTcp
		} else if ser.Tcp != nil {
			protocol = v1.ProtocolTCP
			port = *ser.Tcp
		} else if ser.Udp != nil {
			protocol = v1.ProtocolUDP
			port = *ser.Udp
		} else if ser.Http != nil {
			protocol = v1.ProtocolTCP
			port = *ser.Http
		} else if ser.Https != nil {
			protocol = v1.ProtocolTCP
			port = *ser.Https
		}

		ports = append(ports, v1.ContainerPort{
			Name:          name + "-" + fmt.Sprint(port),
			ContainerPort: int32(port),
			Protocol:      protocol,
		})
	}
	return ports
}

// Get Service Port From Service
// 获取对外服务端口
func (m *Minter) BuildServicePortFormService(name string, services []model.Service) ([]v1.ServicePort, []v1.ServicePort) {
	nodePorts := []v1.ServicePort{}
	headlessPorts := []v1.ServicePort{}
	for i, ser := range services {
		var protocol v1.Protocol
		var port uint16

		// 获取服务端口
		if ser.ProjectUdp != nil {
			protocol = v1.ProtocolUDP
			port = *ser.ProjectUdp
		} else if ser.ProjectTcp != nil {
			protocol = v1.ProtocolTCP
			port = *ser.ProjectTcp
		} else if ser.Tcp != nil {
			protocol = v1.ProtocolTCP
			port = *ser.Tcp
		} else if ser.Udp != nil {
			protocol = v1.ProtocolUDP
			port = *ser.Udp
		} else if ser.Http != nil {
			protocol = v1.ProtocolTCP
			port = *ser.Http
		} else if ser.Https != nil {
			protocol = v1.ProtocolTCP
			port = *ser.Https
		}

		if ser.Tcp != nil || ser.Udp != nil {
			if port != 0 {
				nodePorts = append(nodePorts, v1.ServicePort{
					Name:       name + "-" + fmt.Sprint(i) + "-" + fmt.Sprint(port) + "-nodeport",
					Port:       int32(port),
					TargetPort: intstr.FromInt(int(port)),
					Protocol:   protocol,
				})
			} else {
				nodePort := m.randNodeport()
				if ser.Tcp != nil {
					services[i].Tcp = &nodePort
				}
				if ser.Udp != nil {
					services[i].Udp = &nodePort
				}
				nodePorts = append(nodePorts, v1.ServicePort{
					Name:       name + "-" + fmt.Sprint(i) + "-" + fmt.Sprint(nodePort) + "-nodeport",
					Port:       int32(nodePort),
					TargetPort: intstr.FromInt(int(nodePort)),
					NodePort:   int32(nodePort),
					Protocol:   protocol,
				})
			}
		} else {
			headlessPorts = append(headlessPorts, v1.ServicePort{
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

// Build Pod Container
func (m *Minter) buildPodContainer(
	ctx *context.Context,
	pod model.Pod,
	nameSpace, name string,
	cs []model.Container,
	envs []*gtypes.Env1,
) ([]v1.Container, error) {
	podContainers := make([]v1.Container, 0, len(cs))

	serviceSpace := m.K8sClient.CoreV1().Services(nameSpace)
	nodeports, projectPorts := []v1.ServicePort{}, []v1.ServicePort{}

	// 计算所有的服务端口
	for _, container := range cs {
		cnodeports, cprojectPorts := m.BuildServicePortFormService(name, container.Port)

		// 创建对外端口
		nodeports = append(nodeports, cnodeports...)
		projectPorts = append(projectPorts, cprojectPorts...)
	}

	// 添加机密认证服务
	nodeports = append(nodeports, v1.ServicePort{
		Name:       name + "-65535",
		Protocol:   "TCP",
		Port:       65535,
		TargetPort: intstr.FromInt(65535),
	})

	// 创建对外服务
	aservice := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name + "-expose",
			Labels: map[string]string{"service": name},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"tee_pod": name},
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
			Selector:  map[string]string{"tee_pod": name},
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
		// if i == 0 {
		// 	ports = append(ports, v1.ContainerPort{
		// 		Name:          "port0",
		// 		ContainerPort: int32(65535),
		// 		Protocol:      "TCP",
		// 	})
		// }

		cnevs, err := m.BuildEnvsFromSettings(pod.PodId, filterEnvs(envs, uint16(i)))
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

// 获取工作日志和硬件资源使用量
func (m *Minter) GetLogAndCr(ctx *context.Context, nameSpace string, work model.Pod, now time.Time, stage uint32, isCache bool) ([]string, map[string][]int64, error) {
	// 获取上次记录的时间
	name := "TEE-" + fmt.Sprint(work.PodId)
	if isCache {
		name = name + "-cache"
	}
	from, err := store.GetCacheId(name)
	if err != nil {
		if isCache {
			stage = 9
		}
		from = now.Add(-6 * time.Second * time.Duration(stage)).Unix()
	}

	// 通过 K8s API 获取指定命名空间中的 Pod 列表
	clientset := m.K8sClient
	pods, err := clientset.CoreV1().Pods(nameSpace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: getLabelSelector(work.PodId),
	})
	if err != nil {
		util.LogError("getPod", err)
		return nil, nil, err
	}

	// 判断是否找到匹配的 Pod
	if len(pods.Items) == 0 {
		// 记录错误信息：没有找到匹配的 Pod
		util.LogError("pods is empty")
		return nil, nil, errors.New("pods is empty")
	}

	// 打印获取到的 Pod 名称
	fmt.Println("pods: ", pods.Items[0].Name)

	// 获取指定 Pod 的日志和硬件资源使用量信息
	logs, crs, err := m.getMetricInfo(*ctx, work, nameSpace, pods.Items[0].Name, from)

	// 如果获取 log 和硬件资源使用量的过程中出现错误，则记录错误日志
	if err != nil {
		util.LogError("getMetricInfo", err)
	}

	// 返回获取到的日志和硬件资源使用量
	return logs, crs, nil
}

func contains(s []int32, e int32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getLabelSelector(popId uint64) string {
	// 获取工作类型字符串表示
	name := "TEE-" + fmt.Sprint(popId)
	return name
}

func GetPodName(pid uint64) string {
	return "tee-" + fmt.Sprint(pid)
}
