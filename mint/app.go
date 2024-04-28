package mint

import (
	"context"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"wetee.app/worker/mint/proof"
	"wetee.app/worker/store"
	"wetee.app/worker/util"
)

func (m *Minter) DoWithAppState(ctx *context.Context, c ContractStateWrap, stage uint32, head types.Header) (*gtypes.RuntimeCall, error) {
	if c.App == nil || c.WorkState == nil {
		return nil, errors.New("app is nil")
	}

	app := c.App
	state := c.WorkState

	_, err := m.CheckAppStatus(ctx, c)
	if err != nil {
		util.LogWithRed("checkPodStatus", err)
		return nil, err
	}

	// 判断是否上传工作证明
	// Check if work proof needs to be uploaded
	if uint64(app.Status) != 1 || uint64(head.Number)-state.BlockNumber < uint64(stage) {
		return nil, nil
	}

	util.LogWithRed("=========================================== WorkProofUpload APP")

	workId := c.ContractState.WorkId
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)
	nameSpace := AccountToSpace(c.ContractState.User[:])

	// 获取pod信息
	// Get pod information
	clientset := m.K8sClient
	pods, err := clientset.CoreV1().Pods(nameSpace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=" + name,
	})
	if err != nil {
		util.LogWithRed("getPod", err)
		return nil, err
	}

	if len(pods.Items) == 0 {
		util.LogWithRed("pods is empty")
		return nil, errors.New("pods is empty")
	}
	fmt.Println("pods: ", pods.Items[0].Name)

	// 获取log和硬件资源使用量
	// Get log and hardware resource usage
	logs, crs, err := m.getMetricInfo(*ctx, workId, nameSpace, pods.Items[0].Name, uint64(head.Number)-state.BlockNumber)
	// 如果获取log和硬件资源使用量失败就不提交相关数据
	if err != nil {
		util.LogWithRed("getMetricInfo", err)
	}

	return proof.MakeWorkProof(workId, logs, crs, state.BlockNumber)
}

// checkAppStatus check app status
// 校对应用状态
func (m *Minter) CheckAppStatus(ctx *context.Context, state ContractStateWrap) (*appsv1.Deployment, error) {
	address := AccountToSpace(state.ContractState.User[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(address)
	workId := state.ContractState.WorkId
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	app := state.App
	deployment, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	version := state.Version
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}

		// 重新创建
		envs, err := m.GetEnvsFromSettings(workId, state.Settings)
		if err != nil {
			return nil, err
		}
		err = m.CreateApp(ctx, state.ContractState.User[:], workId, app, envs, version)
		if err != nil {
			return nil, err
		}
		deployment, err = nameSpace.Get(*ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
	}

	return deployment, err
}

// CreateOrUpdateApp create or update app
// 校对应用链上状态后创建或更新应用
func (m *Minter) CreateApp(ctx *context.Context, user []byte, workId gtypes.WorkId, app *gtypes.TeeApp, envs []v1.EnvVar, version uint64) error {
	saddress := AccountToSpace(user)
	errc := m.checkNameSpace(*ctx, saddress)
	if errc != nil {
		return errc
	}

	nameSpace := m.K8sClient.AppsV1().Deployments(saddress)
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	err := store.SetSecrets(workId, &store.Secrets{
		Env: map[string]string{
			"": "",
		},
	})
	if err != nil {
		return err
	}

	ports := GetContainerPortFormService(name, app.Port)
	ports = append(ports, v1.ContainerPort{
		Name:          "port0",
		ContainerPort: int32(8888),
		Protocol:      "TCP",
	})
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: map[string]string{"version": fmt.Sprint(version)},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": name},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "c1",
							Image: string(app.Image),
							Ports: ports,
							Env:   envs,
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:                 resource.MustParse(fmt.Sprint(app.Cr.Cpu) + "m"),
									v1.ResourceMemory:              resource.MustParse(fmt.Sprint(app.Cr.Mem) + "M"),
									"alibabacloud.com/sgx_epc_MiB": *resource.NewQuantity(int64(10), resource.DecimalExponent),
								},
								Requests: v1.ResourceList{
									v1.ResourceCPU:                 resource.MustParse(fmt.Sprint(app.Cr.Cpu) + "m"),
									v1.ResourceMemory:              resource.MustParse(fmt.Sprint(app.Cr.Mem) + "M"),
									"alibabacloud.com/sgx_epc_MiB": *resource.NewQuantity(int64(10), resource.DecimalExponent),
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = nameSpace.Create(*ctx, &deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// 创建机密认证服务
	serviceSpace := m.K8sClient.CoreV1().Services(saddress)
	sports := GetServicePortFormService(name, app.Port)
	sports = append(sports, v1.ServicePort{
		Name:       name + "-8888",
		Protocol:   "TCP",
		Port:       8888,
		TargetPort: intstr.FromInt(8888),
	})

	aservice := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name + "-expose",
			Labels: map[string]string{"service": name},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"app": name},
			Type:     "NodePort",
			Ports:    sports,
		},
	}
	_, err = serviceSpace.Create(*ctx, &aservice, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Println("================================================= Create", err)

	return err
}

// update app service
// 更新APP
func (m *Minter) UpdateApp(ctx *context.Context, user []byte, workId gtypes.WorkId, app *gtypes.TeeApp, envs []v1.EnvVar, version uint64) error {
	saddress := AccountToSpace(user)
	nameSpace := m.K8sClient.AppsV1().Deployments(saddress)
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	existing, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err == nil {
		fmt.Println("================================================= Updating", name)
		existing.ObjectMeta.Annotations = map[string]string{
			"version": fmt.Sprint(version),
		}
		existing.Spec.Template.Spec.Containers[0].Env = envs
		existing.Spec.Template.Spec.Containers[0].Image = string(app.Image)
		_, err = nameSpace.Update(*ctx, existing, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
	}

	return err
}
