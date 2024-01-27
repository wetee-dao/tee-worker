package mint

import (
	"context"
	"encoding/hex"
	"fmt"

	stypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/gen/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/dao"
	"wetee.app/worker/util"
)

func checkTaskStatus(state chain.ContractStateWrap, blockHash stypes.Hash) (*v1.Pod, error) {
	ctx := context.Background()
	k8s := MinterIns.K8sClient.CoreV1()
	address := hex.EncodeToString(state.ContractState.User[:])
	nameSpace := k8s.Pods(address[1:])
	workID := state.ContractState.WorkId
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	pod, err := nameSpace.Get(ctx, name, metav1.GetOptions{})
	if err != nil && err.Error() != "pods \""+name+"\" not found" {
		return nil, err
	}

	version, err := chain.GetVersion(MinterIns.ChainClient, workID)
	if err != nil {
		return nil, err
	}

	if pod.ObjectMeta.Annotations["version"] == fmt.Sprint(version) {
		err := CreateOrUpdateTask(state.ContractState.User[:], workID, blockHash, version)
		if err != nil {
			return nil, err
		}
		return nameSpace.Get(ctx, name, metav1.GetOptions{})
	}

	return pod, nil
}

func CreateOrUpdateTask(user []byte, workID types.WorkId, blockHash stypes.Hash, version uint64) error {
	saddress := AccountToAddress(user[:])
	ctx := context.Background()
	errc := checkNameSpace(ctx, saddress)
	if errc != nil {
		return errc
	}

	appIns := chain.App{
		Client: MinterIns.ChainClient,
		Signer: Signer,
	}
	app, err := appIns.GetApp(user[:], workID.Id)
	if err != nil {
		return err
	}

	k8s := MinterIns.K8sClient.CoreV1()
	nameSpace := k8s.Pods(saddress)
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	err = dao.SetSecrets(workID, &dao.Secrets{
		Env: map[string]string{
			"": "",
		},
	})
	if err != nil {
		return err
	}

	_, err = nameSpace.Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		existingPod, err := nameSpace.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		existingPod.ObjectMeta.Annotations = map[string]string{
			"version": fmt.Sprint(version),
		}
		existingPod.Spec.Containers[0].Image = string(app.Image)
		existingPod.Spec.Containers[0].Ports[0].ContainerPort = int32(app.Port[0])
		_, err = nameSpace.Update(ctx, existingPod, metav1.UpdateOptions{})
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
								"alibabacloud.com/sgx_epc_MiB": *resource.NewMilliQuantity(int64(20), resource.DecimalSI),
							},
							Requests: v1.ResourceList{
								"alibabacloud.com/sgx_epc_MiB": *resource.NewMilliQuantity(int64(20), resource.DecimalSI),
							},
						},
					},
				},
			},
		}
		_, err = nameSpace.Create(ctx, pod, metav1.CreateOptions{})
		fmt.Println("================================================= Create", err)
	}

	return err
}
