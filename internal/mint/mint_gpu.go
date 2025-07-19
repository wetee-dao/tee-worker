package mint

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/internal/mint/proof"
	"wetee.app/worker/internal/util"
)

func (m *Minter) MintGPU(ctx *context.Context, pod model.Pod, stage uint32, currBlock uint32) (*types.Call, int64, error) {
	nameSpace := AccountToSpace(pod.Owner[:])

	// 判断是否上传工作证明
	// Check if work proof needs to be uploaded
	// status 0=>created  1=>deoloying 2=>error  3=>stop
	if uint32(currBlock)-pod.LastMintBlockNumber < stage {
		return nil, 0, nil
	}

	util.LogWithCyan("===========================================", "MINT APP", pod.PodId)
	now := time.Now()
	logs, crs, err := m.GetMetric(ctx, nameSpace, pod, now, stage, false)
	if err != nil {
		util.LogError("GetMetric", err)
		return nil, 0, err
	}

	return proof.MakeWorkProof(pod, logs, crs, now)
}

// checkAppStatus check app status
// 校对应用状态
func (m *Minter) DeployOrUpdateGPU(ctx *context.Context, pod model.Pod) (*appsv1.Deployment, error) {
	address := AccountToSpace(pod.Owner[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(address)
	name := GetPodName(pod.PodId)

	util.LogWithCyan("===========================================", "DEPLOY APP", pod.PodId)
	deployment, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}

		// 重新创建
		err = m.CreateGpuApp(ctx, pod, []*gtypes.Env1{})
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
func (m *Minter) CreateGpuApp(ctx *context.Context, pod model.Pod, envs []*gtypes.Env1) error {
	// get namespace name
	name := GetPodName(pod.PodId)
	saddress := AccountToSpace(pod.Owner[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(saddress)
	err := m.checkNameSpace(*ctx, saddress)
	if err != nil {
		return err
	}

	// build pod
	// build pod ports
	pContainers, err := m.buildPodContainer(ctx, pod, saddress, name, pod.Containers, envs)
	if err != nil {
		return err
	}

	// add gpu resources
	for i := 0; i < len(pContainers); i++ {
		pContainers[i].Resources.Limits["nvidia.com/gpu"] = *resource.NewQuantity(int64(pod.Containers[i].Cr.Gpu), resource.DecimalExponent)
		pContainers[i].Resources.Requests["nvidia.com/gpu"] = *resource.NewQuantity(int64(pod.Containers[i].Cr.Gpu), resource.DecimalExponent)
	}

	// build deployment
	nvidiaClass := "nvidia"
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: map[string]string{"version": fmt.Sprint(pod.Version)},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"tee": name},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"tee": name},
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

	// add model
	m.WrapAiModel(&pod, &deployment)

	// ADD Libos
	m.WrapLibos(&deployment, pod.TeeType)

	_, err = nameSpace.Create(*ctx, &deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return err
}

func (m *Minter) WrapAiModel(app *model.Pod, deployment *appsv1.Deployment) {
	meta := map[string]string{}
	json.Unmarshal([]byte(app.Meta), &meta)
	switch meta["ai-model"] {
	case "sd":
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
	case "ollama":
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
