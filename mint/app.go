package mint

import (
	"context"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/mint/proof"
	"wetee.app/worker/util"
)

// DoWithAppState
// 获取app状态
func (m *Minter) DoWithAppState(ctx *context.Context, c ContractStateWrap, stage uint32, head types.Header) (*gtypes.RuntimeCall, error) {
	if c.App == nil || c.WorkState == nil {
		return nil, errors.New("app is nil")
	}

	app := c.App
	state := c.WorkState

	_, err := m.CheckAppStatus(ctx, c)
	if err != nil {
		util.LogError("checkPodStatus", err)
		return nil, err
	}

	if app.Status != 3 {
		return nil, nil
	}

	// 判断是否上传工作证明
	// Check if work proof needs to be uploaded
	// App状态 0: created, 1: deploying, 2: stop, 3: deoloyed
	if uint64(head.Number)-state.BlockNumber < uint64(stage) {
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
		LabelSelector: "app=" + name,
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
func (m *Minter) CheckAppStatus(ctx *context.Context, state ContractStateWrap) (*appsv1.Deployment, error) {
	address := AccountToSpace(state.ContractState.User[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(address)
	workId := state.ContractState.WorkId
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	app := state.App
	version := state.Version
	deployment, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}

		// 重新创建
		err = m.CreateApp(ctx, state.ContractState.User[:], workId, app, state.Envs, version)
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
func (m *Minter) CreateApp(ctx *context.Context, user []byte, workId gtypes.WorkId, app *gtypes.TeeApp, envs []*gtypes.Env, version uint64) error {
	saddress := AccountToSpace(user)
	err := m.checkNameSpace(*ctx, saddress)
	if err != nil {
		return err
	}

	nameSpace := m.K8sClient.AppsV1().Deployments(saddress)
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	// 构建容器
	main := gtypes.Container{
		Image:   app.Image,
		Command: app.Command,
		Port:    app.Port,
		Cr:      app.Cr,
	}
	cs := append([]gtypes.Container{main}, app.SideContainer...)

	pContainers, err := m.buildPodContainer(ctx, workId, saddress, name, cs, envs)
	if err != nil {
		return err
	}

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: map[string]string{"version": fmt.Sprint(version)},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": name},
				},
				Spec: v1.PodSpec{Containers: pContainers},
			},
		},
	}

	// 初始化磁盘
	err = m.DeploymentPVCWrap(ctx, saddress, name, cs, &deployment)
	if err != nil {
		return err
	}

	// 初始化TEE
	m.DeploymentTEEWrap(&deployment, &app.TeeVersion)

	// ADD Libos
	m.WrapLibos(&deployment, &app.TeeVersion)

	_, err = nameSpace.Create(*ctx, &deployment, metav1.CreateOptions{})
	fmt.Println("================================================= Create pod", err)
	if err != nil {
		return err
	}

	return err
}

// update app service
// 更新APP
func (m *Minter) UpdateApp(ctx *context.Context, user []byte, workId gtypes.WorkId, app *gtypes.TeeApp, envs []v1.EnvVar, version uint64) error {
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
		_, err = nameSpace.Update(*ctx, existing, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
	}

	return err
}
