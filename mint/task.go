package mint

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/wetee-dao/go-sdk/gen/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/dao"
	"wetee.app/worker/util"
)

// check task status，if task is running, return pod, if task not run, create pod
func (m *Minter) CheckTaskStatus(ctx *context.Context, state ContractStateWrap) (*v1.Pod, error) {
	address := hex.EncodeToString(state.ContractState.User[:])
	nameSpace := m.K8sClient.CoreV1().Pods(address[1:])
	workID := state.ContractState.WorkId
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	app := state.Task
	version := state.Version

	pod, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil {
		if err.Error() == "pods \""+name+"\" not found" {
			err := m.CreateTask(ctx, state.ContractState.User[:], workID, app, version)
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
func (m *Minter) CreateTask(ctx *context.Context, user []byte, workID types.WorkId, app *types.TeeTask, version uint64) error {
	saddress := AccountToAddress(user[:])
	errc := m.checkNameSpace(*ctx, saddress)
	if errc != nil {
		return errc
	}

	nameSpace := m.K8sClient.CoreV1().Pods(saddress)
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	err := dao.SetSecrets(workID, &dao.Secrets{
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
		existingPod.ObjectMeta.Annotations = map[string]string{
			"version": fmt.Sprint(version),
		}
		existingPod.Spec.Containers[0].Image = string(app.Image)
		existingPod.Spec.Containers[0].Ports[0].ContainerPort = int32(app.Port[0])
		_, err = nameSpace.Update(*ctx, existingPod, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
	} else {
		// 用于应用联系控制面板的凭证
		wid, err := dao.SealAppID(workID)
		if err != nil {
			return err
		}
		pod := &v1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Task",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Annotations: map[string]string{
					"version": fmt.Sprint(version),
				},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "c1",
						Image: string(app.Image),
						Ports: []v1.ContainerPort{
							{
								Name:          string(app.Name) + "0",
								ContainerPort: int32(app.Port[0]),
								Protocol:      "TCP",
							},
						},
						Env: []v1.EnvVar{
							{
								Name:  "APPID",
								Value: wid,
							},
							{
								Name:  "IN_TEE",
								Value: string("1"),
							},
						},
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
	}

	return err
}
