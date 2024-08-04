package mint

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/mint/proof"
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
		util.LogError("checkPodStatus", err)
		return nil, err
	}

	// 判断是否上传工作证明
	// Check if work proof needs to be uploaded
	if app.Status == 0 || app.Status == 2 || (app.Status == 3 && uint64(head.Number)-uint64(state.BlockNumber) < uint64(stage)) {
		return nil, nil
	}

	util.LogError("=========================================== WorkProofUpload APP")

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
		util.LogError("getPod", err)
		return nil, err
	}

	if len(pods.Items) == 0 {
		util.LogError("pods is empty")
		return nil, errors.New("pods is empty")
	}
	fmt.Println("pods: ", pods.Items[0].Name)

	// 获取log和硬件资源使用量
	// Get log and hardware resource usage
	logs, crs, err := m.getMetricInfo(*ctx, workId, nameSpace, pods.Items[0].Name, uint64(head.Number)-uint64(state.BlockNumber))
	// 如果获取log和硬件资源使用量失败就不提交相关数据
	if err != nil {
		util.LogError("getMetricInfo", err)
	}

	return proof.MakeWorkProof(workId, logs, crs, uint64(state.BlockNumber))
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
		err = m.CreateGpuApp(ctx, state.ContractState.User[:], workId, app, state.Envs, version)
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
func (m *Minter) CreateGpuApp(ctx *context.Context, user []byte, workId gtypes.WorkId, app *gtypes.GpuApp, envs []*gtypes.Env, version uint64) error {
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

	// 构建容器
	main := gtypes.Container{
		Image:   app.Image,
		Command: app.Command,
		Port:    app.Port,
		Cr:      app.Cr,
	}
	cs := append([]gtypes.Container{main}, app.SideContainer...)

	// 构建容器端口
	pContainers, err := m.buildPodContainer(ctx, workId, saddress, name, cs, envs)
	if err != nil {
		return err
	}

	// 添加gpu资源
	pContainers[0].Resources.Limits["nvidia.com/gpu"] = *resource.NewQuantity(int64(app.Cr.Gpu), resource.DecimalExponent)
	pContainers[0].Resources.Requests["nvidia.com/gpu"] = *resource.NewQuantity(int64(app.Cr.Gpu), resource.DecimalExponent)

	nvidiaClass := "nvidia"
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
					Containers:       pContainers,
					NodeSelector: map[string]string{
						"TEE": "CVM-SEV",
					},
				},
			},
		},
	}

	// 添加模型
	m.WrapAiModel(app, &deployment)

	_, err = nameSpace.Create(*ctx, &deployment, metav1.CreateOptions{})
	if err != nil {
		return err
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

func (m *Minter) WrapAiModel(app *gtypes.GpuApp, deployment *appsv1.Deployment) {
	meta := map[string]string{}
	json.Unmarshal([]byte(app.Meta), &meta)
	if meta["ai-model"] == "sd" {
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[0].VolumeMounts, v1.VolumeMount{
			Name:      "model-volume",
			MountPath: "/app/stable-diffusion-webui/models/Stable-diffusion",
			ReadOnly:  true,
		}, v1.VolumeMount{
			Name:      "openai-volume",
			MountPath: "/app/stable-diffusion-webui/openai",
			ReadOnly:  true,
		})

		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, v1.Volume{
			Name: "model-volume",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/home/wetee/AI/SD/model",
				},
			},
		}, v1.Volume{
			Name: "openai-volume",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/home/wetee/AI/SD/openai",
				},
			},
		})
	} else if meta["ai-model"] == "ollama" {
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[0].VolumeMounts, v1.VolumeMount{
			Name:      "ollama-volume",
			MountPath: "/root/.ollama",
			ReadOnly:  false,
		})

		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, v1.Volume{
			Name: "ollama-volume",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/home/wetee/AI/ollama",
				},
			},
		})
	}
}
