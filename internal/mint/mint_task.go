package mint

import (
	"context"
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	inkutil "github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/internal/mint/proof"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/internal/util"
)

// mint task
func (m *Minter) MintTASK(ctx *context.Context, pod model.Pod, stage uint32, blockNumber uint32) (*model.TeeCall, error) {
	task, err := m.DeployOrUpdateTASK(ctx, pod)
	if err != nil {
		util.LogError("CheckTASK", err)
		return nil, err
	}

	// 判断是否上传工作证明
	// Determine whether to upload proof of employment
	if task.Status.Phase != v1.PodSucceeded && task.Status.Phase != v1.PodFailed {
		return nil, nil
	}

	util.LogWithBlue("MINT TASK", "===========================================", pod.PodId)
	nameSpace := AccountToSpace(pod.Owner[:])
	name := GetPodName(pod.PodId)

	// 获取log和硬件资源使用量
	// Obtain the log and hardware resource usage
	t := uint64(blockNumber) - uint64(pod.LastMintBlockNumber)
	from := time.Now().Add(-6 * time.Second * time.Duration(t)).Unix()
	logs, crs, err := m.queryMetric(*ctx, pod, nameSpace, name, from)
	if err != nil {
		util.LogError("getMetricInfo", err)
		return nil, err
	}

	m.StopPod(pod)
	return proof.MakeWorkProof(pod, inkutil.NewNone[types.AccountID](), logs, crs, time.Now())
}

// check task status，if task is running, return pod, if task not run, create pod
func (m *Minter) DeployOrUpdateTASK(ctx *context.Context, pod model.Pod) (*v1.Pod, error) {
	address := AccountToSpace(pod.Owner[:])
	nameSpace := m.K8sClient.CoreV1().Pods(address)
	name := GetPodName(pod.PodId)

	util.LogWithBlue("DEPLOY TASK", "===========================================", pod.PodId)
	task, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil {
		if err.Error() == "pods \""+name+"\" not found" {
			err = m.CreateTask(ctx, pod.Owner[:], pod, pod.Version)
			if err != nil {
				return nil, err
			}
			return nameSpace.Get(*ctx, name, metav1.GetOptions{})
		}

		return nil, err
	}

	return task, nil
}

// create task
func (m *Minter) CreateTask(ctx *context.Context, user []byte, app model.Pod, version uint32) error {
	saddress := AccountToSpace(user[:])
	errc := m.checkNameSpace(*ctx, saddress)
	if errc != nil {
		return errc
	}

	nameSpace := m.K8sClient.CoreV1().Pods(saddress)
	name := GetPodName(app.PodId)

	_, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err == nil {
		existingPod, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		existingPod.ObjectMeta.Annotations = map[string]string{"version": fmt.Sprint(version)}
		existingPod.Spec.Containers[0].Image = string(app.Containers[0].Image)
		_, err = nameSpace.Update(*ctx, existingPod, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
		return err
	}

	// 用于应用联系控制面板的凭证
	wid, err := store.SealAppID(app.PodId)
	if err != nil {
		return err
	}
	initEnvs := InitData{
		Envs: map[string]string{
			"APPID":      wid,
			"PODID":      fmt.Sprint(app.PodId),
			"NAME_SPACE": saddress,
		},
	}

	cenvs, err := m.BuildEnvsFromSettings(app.PodId, saddress, 0, app.Containers[0].Env, &initEnvs)
	if err != nil {
		return err
	}

	container := app.Containers[0]
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{Kind: "Task", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: map[string]string{"version": fmt.Sprint(version)},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "c1",
					Image: string(container.Image),
					Ports: BuildContainerPortFormService(name, container.Port),
					Env:   cenvs,
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
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	m.PodTEEWrap(pod, app.TeeType)

	_, err = nameSpace.Create(*ctx, pod, metav1.CreateOptions{})
	fmt.Println("================================================= Create", err)

	return err
}
