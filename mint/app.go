package mint

import (
	"context"
	"fmt"

	stypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/gen/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/dao"
	"wetee.app/worker/util"
)

// checkAppStatus check app status
// 校对应用状态
func (m *Minter) CheckAppStatus(ctx context.Context, state chain.ContractStateWrap) (*v1.Pod, error) {
	k8s := m.K8sClient.CoreV1()
	address := AccountToAddress(state.ContractState.User[:])
	nameSpace := k8s.Pods(address)
	workID := state.ContractState.WorkId
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	appIns := chain.App{
		Client: m.ChainClient,
		Signer: Signer,
	}
	app, err := appIns.GetApp(state.ContractState.User[:], workID.Id)
	if err != nil {
		return nil, err
	}
	if uint8(app.Status) == 2 {
		m.StopApp(workID)
		return nil, errors.New("app stop")
	}

	pod, err := nameSpace.Get(ctx, name, metav1.GetOptions{})
	if err != nil && err.Error() != "pods \""+name+"\" not found" {
		return nil, err
	}

	version, err := chain.GetVersion(m.ChainClient, workID)
	if err != nil {
		return nil, err
	}

	if pod.ObjectMeta.Annotations["version"] != fmt.Sprint(version) {
		err = m.CreateApp(ctx, state.ContractState.User[:], workID, version)
		if err != nil {
			return nil, err
		}
		pod, err = nameSpace.Get(ctx, name, metav1.GetOptions{})
	}

	return pod, err
}

// CreateOrUpdateApp create or update app
// 校对应用链上状态后创建或更新应用
func (m *Minter) CreateApp(ctx context.Context, user []byte, workID types.WorkId, version uint64) error {
	saddress := AccountToAddress(user)
	errc := m.checkNameSpace(ctx, saddress)
	if errc != nil {
		return errc
	}

	appIns := chain.App{
		Client: m.ChainClient,
		Signer: Signer,
	}
	app, err := appIns.GetApp(user[:], workID.Id)
	if err != nil {
		return err
	}

	k8s := m.K8sClient.CoreV1()
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

	existingPod, err := nameSpace.Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		if uint8(app.Status) == 2 {
			m.StopApp(workID)
			return nil
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
				Kind:       "App",
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
			},
		}
		_, err = nameSpace.Create(ctx, pod, metav1.CreateOptions{})
		fmt.Println("================================================= Create", err)
	}

	return err
}

func (m *Minter) UpdateApp(ctx context.Context, user []byte, workID types.WorkId, blockHash stypes.Hash, version uint64) error {
	k8s := m.K8sClient.CoreV1()

	saddress := AccountToAddress(user)
	nameSpace := k8s.Pods(saddress)
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	appIns := chain.App{
		Client: m.ChainClient,
		Signer: Signer,
	}
	app, err := appIns.GetApp(user[:], workID.Id)
	if err != nil {
		return err
	}

	existingPod, err := nameSpace.Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		if uint8(app.Status) == 2 {
			m.StopApp(workID)
			return nil
		}
		existingPod.ObjectMeta.Annotations = map[string]string{
			"version": fmt.Sprint(version),
		}
		existingPod.Spec.Containers[0].Image = string(app.Image)
		existingPod.Spec.Containers[0].Ports[0].ContainerPort = int32(app.Port[0])
		_, err = nameSpace.Update(ctx, existingPod, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
	}

	return err
}

// StopApp
// 停止应用
func (m *Minter) StopApp(workID types.WorkId) error {
	ctx := context.Background()
	user, err := chain.GetAccount(m.ChainClient, workID)
	if err != nil {
		return err
	}

	saddress := AccountToAddress(user[:])

	k8s := m.K8sClient.CoreV1()
	nameSpace := k8s.Pods(saddress)
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)
	return nameSpace.Delete(ctx, name, metav1.DeleteOptions{})
}
