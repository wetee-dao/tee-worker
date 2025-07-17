package mint

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/internal/mint/proof"
	"wetee.app/worker/internal/util"
)

// DoWithApp
// 获取app状态
func (m *Minter) DoAPP(ctx *context.Context, pod model.Pod, stage uint32, currBlock uint32) (*types.Call, int64, error) {
	util.LogWithCyan("===========================================", "DEPLOY APP", pod.PodId)
	_, err := m.CheckAPP(ctx, pod)
	if err != nil {
		util.LogError("CheckAPP", err)
		return nil, 0, err
	}

	nameSpace := AccountToSpace(pod.Owner[:])

	// Check if work proof needs to be uploaded
	// status 0=>created  1=>deoloying 2=>error  3=>stop
	if currBlock-pod.LastMintBlockNumber < stage {
		// 如果当前区块高度小于当前工作高度+阶段高度则不上传工作证明 但是保存工作证明到本地
		// logs, crs, err := m.GetLogAndCr(ctx, nameSpace, app, now, stage, true)
		// if err != nil {
		// 	util.LogError("getMetricInfo", err)
		// 	return nil, 0, err
		// }
		// return nil, 0, proof.CacheWorkProof(app.PodId, logs, crs, now, uint64(blockNumber))
	}

	// Get logs and use compute resource
	now := time.Now()
	logs, crs, err := m.GetMetric(ctx, nameSpace, pod, now, stage, false)
	if err != nil {
		return nil, 0, errors.Wrap(err, "GetLogAndCr")
	}

	// make pod start proof
	return proof.MakeWorkProof(pod, logs, crs, now)
}

// CheckAPP check app status
// 校对应用状态
func (m *Minter) CheckAPP(ctx *context.Context, pod model.Pod) (*appsv1.Deployment, error) {
	// get namespace name
	name := GetPodName(pod.PodId)
	address := AccountToSpace(pod.Owner[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(address)

	// get deployment
	deployment, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, errors.Wrap(err, "CheckAPP Get Deployment")
		}
	}

	// delete pod
	if deployment.Name != "" {
		err := m.StopPod(pod)
		if err != nil {
			return nil, errors.Wrap(err, "StopApp")
		}
	}

	// create pod
	deployment, err = m.CreateApp(ctx, pod, []*gtypes.Env1{})
	if err != nil {
		return nil, errors.Wrap(err, "CreateApp")
	}

	return deployment, err
}

// CreateOrUpdateApp create or update app
// 校对应用链上状态后创建或更新应用
func (m *Minter) CreateApp(ctx *context.Context, pod model.Pod, envs []*gtypes.Env1) (*appsv1.Deployment, error) {
	// get namespace name
	name := GetPodName(pod.PodId)
	nameSpaceStr := AccountToSpace(pod.Owner[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(nameSpaceStr)
	err := m.checkNameSpace(*ctx, nameSpaceStr)
	if err != nil {
		return nil, err
	}

	// build pod
	// build pod ports
	pContainers, err := m.buildPodContainer(ctx, pod, nameSpaceStr, name, pod.Containers, envs)
	if err != nil {
		return nil, err
	}

	// build deployment
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
				Spec: v1.PodSpec{Containers: pContainers},
			},
		},
	}

	// init PVC
	err = m.DeploymentPVCWrap(ctx, nameSpaceStr, name, pod.Containers, &deployment)
	if err != nil {
		return nil, err
	}

	// wrap TEE
	m.DeploymentTEEWrap(&deployment, pod.TeeType)

	// add Libos
	m.WrapLibos(&deployment, pod.TeeType)

	// create pod
	deployed, err := nameSpace.Create(*ctx, &deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "Create Deployment")
	}

	return deployed, nil
}
