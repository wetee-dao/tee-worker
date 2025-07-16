package mint

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/internal/mint/proof"
	"wetee.app/worker/internal/util"
)

// DoWithAppState
// 获取app状态
func (m *Minter) DoWithAppState(ctx *context.Context, app model.Pod, stage uint32, blockNumber uint32) (*types.Call, int64, error) {
	_, err := m.CheckAppStatus(ctx, app)
	if err != nil {
		util.LogError("checkPodStatus", err)
		return nil, 0, err
	}

	// if app.Status != 3 {
	// 	return nil, 0, nil
	// }

	nameSpace := AccountToSpace(app.Owner[:])
	now := time.Now()

	// 判断是否上传工作证明
	// Check if work proof needs to be uploaded
	// App状态 0: created, 1: deploying, 2: stop, 3: deoloyed
	if blockNumber-app.LastMintBlockNumber < stage {
		// if (uint64(blockNumber)+app.PodId)%10 != 0 {
		// 	return nil, 0, nil
		// }

		// 如果当前区块高度小于当前工作高度+阶段高度则不上传工作证明 但是保存工作证明到本地
		// logs, crs, err := m.GetLogAndCr(ctx, nameSpace, app, now, stage, true)
		// if err != nil {
		// 	util.LogError("getMetricInfo", err)
		// 	return nil, 0, err
		// }
		// return nil, 0, proof.CacheWorkProof(app.PodId, logs, crs, now, uint64(blockNumber))
	}

	util.LogError("=========================================== WorkProofUpload APP")

	logs, crs, err := m.GetLogAndCr(ctx, nameSpace, app, now, stage, false)
	// if err != nil {
	// 	util.LogError("GetLogAndCr", err)
	// 	return nil, 0, err
	// }

	return proof.MakeWorkProof(app, logs, crs, now, uint64(app.LastMintBlockNumber))
}

// checkAppStatus check app status
// 校对应用状态
func (m *Minter) CheckAppStatus(ctx *context.Context, app model.Pod) (*appsv1.Deployment, error) {
	address := AccountToSpace(app.Owner[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(address)
	name := GetPodName(app.PodId)

	version := app.Version
	deployment, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}

		// 重新创建
		err = m.CreateApp(ctx, app.Owner[:], app, []*gtypes.Env1{}, version)
		if err != nil {
			return nil, err
		}

		deployment, err = nameSpace.Get(*ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
	} else {

	}

	return deployment, err
}

// CreateOrUpdateApp create or update app
// 校对应用链上状态后创建或更新应用
func (m *Minter) CreateApp(ctx *context.Context, user []byte, app model.Pod, envs []*gtypes.Env1, version uint32) error {
	saddress := AccountToSpace(user)
	err := m.checkNameSpace(*ctx, saddress)
	if err != nil {
		return err
	}

	nameSpace := m.K8sClient.AppsV1().Deployments(saddress)
	name := GetPodName(app.PodId)

	// 构建容器
	pContainers, err := m.buildPodContainer(ctx, app, saddress, name, app.Containers, envs)
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
	err = m.DeploymentPVCWrap(ctx, saddress, name, app.Containers, &deployment)
	if err != nil {
		return err
	}

	// 初始化TEE
	m.DeploymentTEEWrap(&deployment, app.TeeType)

	// ADD Libos
	m.WrapLibos(&deployment, app.TeeType)

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
	name := GetPodName(workId.Id)

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
