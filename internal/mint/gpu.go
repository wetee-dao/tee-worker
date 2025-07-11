package mint

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/internal/mint/proof"
	"wetee.app/worker/internal/util"
)

func (m *Minter) DoWithGpuAppState(ctx *context.Context, app model.Pod, stage uint32, blockNumber uint32) (*gtypes.RuntimeCall, error) {
	_, err := m.CheckGpuAppStatus(ctx, app)
	if err != nil {
		util.LogError("checkPodStatus", err)
		return nil, err
	}

	if app.Status != 3 {
		return nil, nil
	}

	nameSpace := AccountToSpace(app.Owner[:])
	now := time.Now()

	// 判断是否上传工作证明
	// Check if work proof needs to be uploaded
	// App状态 0: created, 1: deploying, 2: stop, 3: deoloyed
	if uint32(blockNumber)-app.LastMintBlockNumber < stage {
		if (uint64(blockNumber)+app.PodId)%10 != 0 {
			return nil, nil
		}
		// 如果当前区块高度小于当前工作高度+阶段高度则不上传工作证明 但是保存工作证明到本地
		logs, crs, err := m.GetLogAndCr(ctx, nameSpace, app, now, stage, true)
		if err != nil {
			util.LogError("getMetricInfo", err)
			return nil, err
		}
		return nil, proof.CacheWorkProof(app.PodId, logs, crs, now, uint64(blockNumber))
	}

	util.LogError("=========================================== WorkProofUpload GPU")

	logs, crs, err := m.GetLogAndCr(ctx, nameSpace, app, now, stage, false)
	if err != nil {
		util.LogError("getMetricInfo", err)
		return nil, err
	}

	return proof.MakeWorkProof(app, logs, crs, now, uint64(app.LastMintBlockNumber))
}

// checkAppStatus check app status
// 校对应用状态
func (m *Minter) CheckGpuAppStatus(ctx *context.Context, app model.Pod) (*appsv1.Deployment, error) {
	address := AccountToSpace(app.Owner[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(address)
	name := GetPodName(app.PodId)

	deployment, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}

		// 重新创建
		err = m.CreateGpuApp(ctx, app.Owner[:], app, []*gtypes.Env1{}, uint64(app.Version))
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
func (m *Minter) CreateGpuApp(ctx *context.Context, user []byte, app model.Pod, envs []*gtypes.Env1, version uint64) error {
	saddress := AccountToSpace(user)
	err := m.checkNameSpace(*ctx, saddress)
	if err != nil {
		return err
	}

	nameSpace := m.K8sClient.AppsV1().Deployments(saddress)
	name := GetPodName(app.PodId)

	// 构建容器
	// 构建容器端口
	pContainers, err := m.buildPodContainer(ctx, app, saddress, name, app.Containers, envs)
	if err != nil {
		return err
	}

	// 添加gpu资源
	for i := 0; i < len(pContainers); i++ {
		pContainers[i].Resources.Limits["nvidia.com/gpu"] = *resource.NewQuantity(int64(app.Containers[i-1].Cr.Gpu), resource.DecimalExponent)
		pContainers[i].Resources.Requests["nvidia.com/gpu"] = *resource.NewQuantity(int64(app.Containers[i-1].Cr.Gpu), resource.DecimalExponent)
	}

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
	m.WrapAiModel(&app, &deployment)

	// ADD Libos
	m.WrapLibos(&deployment, app.TeeType)

	_, err = nameSpace.Create(*ctx, &deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return err
}

func (m *Minter) UpdateGpuApp(ctx *context.Context, user []byte, workId gtypes.WorkId, app *gtypes.GpuApp, envs []v1.EnvVar, version uint64) error {
	saddress := AccountToSpace(user)
	nameSpace := m.K8sClient.AppsV1().Deployments(saddress)
	name := GetPodName(workId.Id)

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

func (m *Minter) WrapAiModel(app *model.Pod, deployment *appsv1.Deployment) {
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
