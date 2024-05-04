package mint

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/mint/proof"
	"wetee.app/worker/store"
	"wetee.app/worker/util"
)

func (m *Minter) DoWithGpuAppState(ctx *context.Context, c ContractStateWrap, stage uint32, head types.Header) (*gtypes.RuntimeCall, error) {
	if c.GpuApp == nil || c.WorkState == nil {
		return nil, errors.New("app is nil")
	}

	app := c.GpuApp
	state := c.WorkState

	_, err := m.CheckGpuAppStatus(ctx, c)
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
		LabelSelector: "gpu=" + name,
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
func (m *Minter) CheckGpuAppStatus(ctx *context.Context, state ContractStateWrap) (*appsv1.Deployment, error) {
	address := AccountToSpace(state.ContractState.User[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(address)
	workId := state.ContractState.WorkId
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	app := state.GpuApp
	deployment, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	version := state.Version
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}

		// 重新创建
		envs, err := m.GetEnvsFromSettings(workId, state.Envs)
		if err != nil {
			return nil, err
		}
		err = m.CreateGpuApp(ctx, state.ContractState.User[:], workId, app, envs, version)
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
func (m *Minter) CreateGpuApp(ctx *context.Context, user []byte, workId gtypes.WorkId, app *gtypes.GpuApp, envs []v1.EnvVar, version uint64) error {
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

	nvidiaClass := "nvidia"
	metaJson := map[string]string{}
	json.Unmarshal(app.Meta, &metaJson)
	command := []string{}
	if c, ok := metaJson["c"]; ok {
		command = strings.Split(c, " ")
	}

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: map[string]string{"version": fmt.Sprint(version)},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"gpu": name},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"gpu": name},
				},
				Spec: v1.PodSpec{
					RuntimeClassName: &nvidiaClass,
					NodeSelector: map[string]string{
						"TEE": "CVM-SEV",
					},
					Containers: []v1.Container{
						{
							Name:  "c1",
							Image: string(app.Image),
							Ports: GetContainerPortFormService(name, app.Port),
							Env: []v1.EnvVar{
								{Name: "IN_TEE", Value: string("1")},
								{Name: "COMMANDLINE_ARGS", Value: string(" --no-half-vae --lowvram --share --xformers ")},
							},
							Command: command,
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse(fmt.Sprint(app.Cr.Cpu) + "m"),
									v1.ResourceMemory: resource.MustParse(fmt.Sprint(app.Cr.Mem) + "M"),
									"nvidia.com/gpu":  *resource.NewQuantity(int64(app.Cr.Gpu), resource.DecimalExponent),
								},
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse(fmt.Sprint(app.Cr.Cpu) + "m"),
									v1.ResourceMemory: resource.MustParse(fmt.Sprint(app.Cr.Mem) + "M"),
									"nvidia.com/gpu":  *resource.NewQuantity(int64(app.Cr.Gpu), resource.DecimalExponent),
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "model-volume",
									MountPath: "/app/stable-diffusion-webui/models/Stable-diffusion",
								},
								{
									Name:      "openai-volume",
									MountPath: "/app/stable-diffusion-webui/openai",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "model-volume",
							VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{
									Path: "/home/wetee/work/wetee/worker/AI/model",
								},
							},
						},
						{
							Name: "openai-volume",
							VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{
									Path: "/home/wetee/work/wetee/worker/AI/openai",
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

	// 创建机密服务
	sports := GetServicePortFormService(name, app.Port)
	if len(sports) > 0 {
		ServiceSpace := m.K8sClient.CoreV1().Services(saddress)
		service := v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:   name + "-expose",
				Labels: map[string]string{"service": name},
			},
			Spec: v1.ServiceSpec{
				Selector: map[string]string{"gpu": name},
				Type:     "NodePort",
				Ports:    GetServicePortFormService(name, app.Port),
			},
		}
		_, err = ServiceSpace.Create(*ctx, &service, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	return err
}

func (m *Minter) UpdateGpuApp(ctx *context.Context, user []byte, workId gtypes.WorkId, app *gtypes.GpuApp, envs []v1.EnvVar, version uint64) error {
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
		// existing.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = int32(app.Port[0])
		_, err = nameSpace.Update(*ctx, existing, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
	}

	return err
}
