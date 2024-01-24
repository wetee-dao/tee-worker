package mint

import (
	"context"
	"encoding/hex"
	"fmt"

	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/gen/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/dao"
	"wetee.app/worker/util"
)

func checkPodStatus(state chain.ContractStateWrap, blockHash string) error {
	ctx := context.Background()
	k8s := MinterIns.K8sClient.CoreV1()
	address := hex.EncodeToString(state.ContractState.User[:])
	nameSpace := k8s.Pods(address[1:])
	workID := state.ContractState.WorkId
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	pod, err := nameSpace.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if pod.ObjectMeta.ResourceVersion == state.BlockHash {
		return nil
	}

	return CreateOrUpdatePod(state.ContractState.User[:], workID, blockHash)
}

func CreateOrUpdatePod(user []byte, workID types.WorkId, blockHash string) error {
	address := hex.EncodeToString(user[:])
	saddress := address[1:] //去掉前面的 0x
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
		existingPod.ObjectMeta.ResourceVersion = blockHash
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
				Kind:       "App",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            name,
				ResourceVersion: blockHash,
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "c1",
						Image: string(app.Image),
						Ports: []v1.ContainerPort{
							{
								Name:          string(app.Name) + "0",
								ContainerPort: 80,
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
					},
				},
			},
		}
		_, err = nameSpace.Create(ctx, pod, metav1.CreateOptions{})
		fmt.Println("================================================= Create", err)
	}

	return err
}

func StopPod(workID types.WorkId) error {
	ctx := context.Background()
	user, err := chain.GetAccount(MinterIns.ChainClient, workID)
	if err != nil {
		return err
	}
	address := hex.EncodeToString(user[:])
	saddress := address[1:] //去掉前面的 0x

	k8s := MinterIns.K8sClient.CoreV1()
	nameSpace := k8s.Pods(saddress)
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)
	return nameSpace.Delete(ctx, name, metav1.DeleteOptions{})
}
