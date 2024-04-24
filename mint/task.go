package mint

import (
	"context"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/mint/proof"
	"wetee.app/worker/store"
	"wetee.app/worker/util"
)

func (m *Minter) DoWithTaskState(ctx *context.Context, c ContractStateWrap, stage uint32, head types.Header) (*gtypes.RuntimeCall, error) {
	if c.Task == nil || c.WorkState == nil {
		return nil, errors.New("task is nil")
	}

	app := c.Task
	state := c.WorkState

	// 处于调度状态，不处理
	if uint64(app.Status) == 4 {
		return nil, nil
	}

	pod, err := m.CheckTaskStatus(ctx, c)
	if err != nil {
		util.LogWithRed("checkTaskStatus", err)
		return nil, err
	}

	// 判断是否上传工作证明
	// Determine whether to upload proof of employment
	if pod.Status.Phase != v1.PodSucceeded && pod.Status.Phase != v1.PodFailed {
		return nil, nil
	}
	util.LogWithRed("===========================================WorkProofUpload TASK")
	nameSpace := AccountToSpace(c.ContractState.User[:])
	workId := c.ContractState.WorkId
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	// 获取log和硬件资源使用量
	// Obtain the log and hardware resource usage
	logs, crs, err := m.getMetricInfo(*ctx, workId, nameSpace, name, uint64(head.Number)-state.BlockNumber)
	if err != nil {
		util.LogWithRed("getMetricInfo", err)
		return nil, err
	}

	m.StopApp(c.ContractState.WorkId, "")
	return proof.MakeWorkProof(workId, logs, crs, state.BlockNumber)
}

// check task status，if task is running, return pod, if task not run, create pod
func (m *Minter) CheckTaskStatus(ctx *context.Context, state ContractStateWrap) (*v1.Pod, error) {
	address := AccountToSpace(state.ContractState.User[:])
	nameSpace := m.K8sClient.CoreV1().Pods(address)
	workId := state.ContractState.WorkId
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	app := state.Task
	version := state.Version

	pod, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil {
		if err.Error() == "pods \""+name+"\" not found" {
			envs, err := m.GetEnvsFromSettings(workId, state.Settings)
			if err != nil {
				return nil, err
			}
			err = m.CreateTask(ctx, state.ContractState.User[:], workId, app, envs, version)
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
func (m *Minter) CreateTask(ctx *context.Context, user []byte, workId gtypes.WorkId, app *gtypes.TeeTask, envs []v1.EnvVar, version uint64) error {
	saddress := AccountToSpace(user[:])
	errc := m.checkNameSpace(*ctx, saddress)
	if errc != nil {
		return errc
	}

	nameSpace := m.K8sClient.CoreV1().Pods(saddress)
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	err := store.SetSecrets(workId, &store.Secrets{
		Env: map[string]string{
			"": "",
		},
	})
	if err != nil {
		return err
	}

	_, err = nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err == nil {
		existingPod, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		existingPod.ObjectMeta.Annotations = map[string]string{"version": fmt.Sprint(version)}
		existingPod.Spec.Containers[0].Image = string(app.Image)
		existingPod.Spec.Containers[0].Ports[0].ContainerPort = int32(app.Port[0])
		_, err = nameSpace.Update(*ctx, existingPod, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
		return err
	}

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
					Image: string(app.Image),
					Ports: []v1.ContainerPort{
						{
							Name:          name + "-" + "0",
							ContainerPort: int32(app.Port[0]),
							Protocol:      "TCP",
						},
					},
					Env: envs,
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							"alibabacloud.com/sgx_epc_MiB": *resource.NewQuantity(int64(20), resource.DecimalExponent),
						},
						Requests: v1.ResourceList{
							"alibabacloud.com/sgx_epc_MiB": *resource.NewQuantity(int64(20), resource.DecimalExponent),
						},
					},
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	_, err = nameSpace.Create(*ctx, pod, metav1.CreateOptions{})
	fmt.Println("================================================= Create", err)

	return err
}
