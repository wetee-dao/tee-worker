package mint

import (
	"bytes"
	"context"
	"html/template"
	"math/rand"
	"slices"

	"fmt"

	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

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
	nodePorts, teePorts := []v1.ServicePort{}, []v1.ServicePort{}

	// 计算所有的服务端口
	for i, container := range cs {
		nodeports, teeports := m.BuildServicePortFormService(name, i, container.Port)

		nodePorts = append(nodePorts, nodeports...)
		teePorts = append(teePorts, teeports...)
	}

	// 添加机密认证服务
	nodePorts = append(nodePorts, v1.ServicePort{
		Name:       name + "-65535",
		Protocol:   "TCP",
		Port:       65535,
		TargetPort: intstr.FromInt(65535),
	})

	// 创建对外服务
	service := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name + "-expose",
			Labels: map[string]string{"service": name},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"tee": name},
			Type:     "NodePort",
			Ports:    nodePorts,
		},
	}

	nodeSers, err := serviceSpace.Create(*ctx, &service, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("====== CREATE service", err)
		return nil, err
	}

	// 创建项目内端口
	teeService := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"service": name},
		},
		Spec: v1.ServiceSpec{
			Selector:  map[string]string{"tee": name},
			ClusterIP: "None",
			Ports:     teePorts,
		},
	}
	_, err = serviceSpace.Create(*ctx, &teeService, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("====== CREATE tee service", err)
		return nil, err
	}

	// 构建容器
	for i, container := range cs {
		// 获取服务端口
		ports := BuildContainerPortFormService(name, container.Port)

		// 构建来自用户的环境变量
		containerEnvs, err := m.BuildEnvsFromSettings(pod.PodId, filterEnvs(envs, uint16(i)))
		if err != nil {
			fmt.Println("====== CREATE user envs", err)
			return nil, err
		}

		// 构建来自集群的环境变量
		err = m.WrapEnvs(containerEnvs, nameSpace, name, nodeSers)
		if err != nil {
			fmt.Println("====== CREATE cluster envs", err)
			return nil, err
		}

		podContainers = append(podContainers, v1.Container{
			Name:    "c" + fmt.Sprint(i),
			Image:   string(container.Image),
			Ports:   ports,
			Env:     containerEnvs,
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

// Get Service Port From Service
// 获取对外服务端口
func (m *Minter) BuildServicePortFormService(name string, index int, services []model.Service) ([]v1.ServicePort, []v1.ServicePort) {
	nodePorts := []v1.ServicePort{}
	headlessPorts := []v1.ServicePort{}
	name = name + "-" + fmt.Sprint(index)

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

// StopApp
// 停止应用
func (m *Minter) StopPod(p model.Pod) error {
	ctx := context.Background()
	name := GetPodName(p.PodId)
	space := AccountToSpace(p.Owner[:])

	serviceSpace := m.K8sClient.CoreV1().Services(space)
	list, err := serviceSpace.List(ctx, metav1.ListOptions{
		LabelSelector: "service=" + name,
	})
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		err := serviceSpace.Delete(ctx, item.Name, metav1.DeleteOptions{})
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

// renderTemplate
func renderTemplate(templateString string, data map[string]string) (string, error) {
	tmpl, err := template.New("template").Parse(templateString)
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

func contains(s []int32, e int32) bool {
	return slices.Contains(s, e)
}

func GetPodName(pid uint64) string {
	return "tee-" + fmt.Sprint(pid)
}
