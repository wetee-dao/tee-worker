package mint

import (
	"context"
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/internal/mint/proof"
	"wetee.app/worker/internal/util"
)

func (m *Minter) DoWithTaskState(ctx *context.Context, app model.Pod, stage uint32, blockNumber uint32) (*types.Call, int64, error) {
	// 处于调度状态，不处理
	if uint64(app.Status) == 4 {
		return nil, 0, nil
	}

	pod, err := m.CheckTaskStatus(ctx, app)
	if err != nil {
		util.LogError("checkTaskStatus", err)
		return nil, 0, err
	}

	// 判断是否上传工作证明
	// Determine whether to upload proof of employment
	if pod.Status.Phase != v1.PodSucceeded && pod.Status.Phase != v1.PodFailed {
		return nil, 0, nil
	}
	util.LogError("===========================================WorkProofUpload TASK")
	nameSpace := AccountToSpace(app.Owner[:])
	name := GetPodName(app.PodId)

	// 获取log和硬件资源使用量
	// Obtain the log and hardware resource usage
	t := uint64(blockNumber) - uint64(app.LastMintBlockNumber)
	from := time.Now().Add(-6 * time.Second * time.Duration(t)).Unix()
	logs, crs, err := m.getMetricInfo(*ctx, app, nameSpace, name, from)
	if err != nil {
		util.LogError("getMetricInfo", err)
		return nil, 0, err
	}

	m.StopApp(app)

	now := time.Now()
	return proof.MakeWorkProof(app, logs, crs, now, uint64(app.LastMintBlockNumber))
}

// check task status，if task is running, return pod, if task not run, create pod
func (m *Minter) CheckTaskStatus(ctx *context.Context, app model.Pod) (*v1.Pod, error) {
	address := AccountToSpace(app.Owner[:])
	nameSpace := m.K8sClient.CoreV1().Pods(address)
	name := GetPodName(app.PodId)

	pod, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil {
		if err.Error() == "pods \""+name+"\" not found" {
			err = m.CreateTask(ctx, app.Owner[:][:], app, []*gtypes.Env1{}, app.Version)
			if err != nil {
				return nil, err
			}
			return nameSpace.Get(*ctx, name, metav1.GetOptions{})
		}

		return nil, err
	}

	return pod, nil
}

// create task
func (m *Minter) CreateTask(ctx *context.Context, user []byte, app model.Pod, envs []*gtypes.Env1, version uint32) error {
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

	cenvs, err := m.BuildEnvsFromSettings(app.PodId, envs)
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
