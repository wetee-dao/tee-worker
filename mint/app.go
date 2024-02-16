package mint

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	chain "github.com/wetee-dao/go-sdk"
	gtype "github.com/wetee-dao/go-sdk/gen/types"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/dao"
	"wetee.app/worker/util"
)

func (m *Minter) DoWithAppState(ctx *context.Context, c ContractStateWrap, stage uint32, head types.Header) error {
	if c.App == nil || c.WorkState == nil {
		return errors.New("app is nil")
	}

	app := c.App
	state := c.WorkState

	// 状态为停止状态，停止Pod
	if uint64(app.Status) == 2 {
		m.StopApp(c.ContractState.WorkId)
		return nil
	}

	_, err := m.CheckAppStatus(ctx, c)
	if err != nil {
		util.LogWithRed("checkPodStatus", err)
		return err
	}

	// 判断是否上传工作证明
	// Check if work proof needs to be uploaded
	if uint64(app.Status) == 1 && uint64(head.Number)-state.BlockNumber >= uint64(stage) {
		util.LogWithRed("=========================================== WorkProofUpload APP")

		workId := c.ContractState.WorkId
		name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)
		nameSpace := AccountToAddress(c.ContractState.User[:])

		// 获取pod信息
		// Get pod information
		clientset := m.K8sClient
		pods, err := clientset.CoreV1().Pods(nameSpace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: "app=" + name,
		})
		fmt.Println("pods: ", pods.Items[0].Name)
		if err != nil {
			util.LogWithRed("getMetricInfo", err)
			return err
		}

		// 获取log和硬件资源使用量
		// Get log and hardware resource usage
		logs, crs, err := m.getMetricInfo(*ctx, workId, nameSpace, pods.Items[0].Name, uint64(head.Number)-state.BlockNumber)
		if err != nil {
			util.LogWithRed("getMetricInfo", err)
			return err
		}

		// 获取log hash
		// Get log hash
		logHash, err := getWorkLogHash(name, logs, state.BlockNumber)
		if err != nil {
			util.LogWithRed("getWorkLogHash", err)
			return err
		}

		// 获取计算资源hash
		// Get Computing resource hash
		crHash, cr, err := getWorkCrHash(name, crs, state.BlockNumber)
		if err != nil {
			util.LogWithRed("getWorkCrHash", err)
			return err
		}

		// 初始化worker对象
		// init worker object
		worker := chain.Worker{
			Client: m.ChainClient,
			Signer: Signer,
		}

		// 上传工作证明
		// Upload work proof
		err = worker.WorkProofUpload(c.ContractState.WorkId, logHash, crHash, gtype.Cr{
			Cpu:  cr[0],
			Mem:  cr[1],
			Disk: 0,
		}, []byte(""), false)
		if err != nil {
			util.LogWithRed("WorkProofUpload", err)
			return err
		}
	}

	return nil
}

// checkAppStatus check app status
// 校对应用状态
func (m *Minter) CheckAppStatus(ctx *context.Context, state ContractStateWrap) (*appsv1.Deployment, error) {
	address := AccountToAddress(state.ContractState.User[:])
	nameSpace := m.K8sClient.AppsV1().Deployments(address)
	workId := state.ContractState.WorkId
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	app := state.App
	if uint8(app.Status) == 2 {
		m.StopApp(workId)
		return nil, errors.New("app stop")
	}

	deployment, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	version := state.Version
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}

		// 重新创建
		envs, err := m.GetEnvsFromSettings(workId, state.Settings)
		if err != nil {
			return nil, err
		}
		err = m.CreateApp(ctx, state.ContractState.User[:], workId, app, envs, version)
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
func (m *Minter) CreateApp(ctx *context.Context, user []byte, workId gtype.WorkId, app *gtype.TeeApp, envs []v1.EnvVar, version uint64) error {
	saddress := AccountToAddress(user)
	errc := m.checkNameSpace(*ctx, saddress)
	if errc != nil {
		return errc
	}

	nameSpace := m.K8sClient.AppsV1().Deployments(saddress)
	name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)

	err := dao.SetSecrets(workId, &dao.Secrets{
		Env: map[string]string{
			"": "",
		},
	})
	if err != nil {
		return err
	}

	if uint8(app.Status) == 2 {
		m.StopApp(workId)
		return nil
	}

	resource.NewMilliQuantity(int64(app.Cr.Mem)*1024*1024, resource.BinarySI)

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
							Env: envs,
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									// v1.ResourceCPU:                 resource.MustParse(fmt.Sprint(app.Cr.Cpu) + "m"),
									// v1.ResourceMemory:              resource.MustParse(fmt.Sprint(app.Cr.Mem) + "M"),
									"alibabacloud.com/sgx_epc_MiB": *resource.NewQuantity(int64(20), resource.DecimalExponent),
								},
								Requests: v1.ResourceList{
									// v1.ResourceCPU:                 resource.MustParse(fmt.Sprint(app.Cr.Cpu) + "m"),
									// v1.ResourceMemory:              resource.MustParse(fmt.Sprint(app.Cr.Mem) + "M"),
									"alibabacloud.com/sgx_epc_MiB": *resource.NewQuantity(int64(20), resource.DecimalExponent),
								},
							},
						},
					},
				},
			},
		},
	}

	bt, _ := json.Marshal(20)
	fmt.Println(string(bt))

	_, err = nameSpace.Create(*ctx, &deployment, metav1.CreateOptions{})
	fmt.Println("================================================= Create", err)

	return err
}

func (m *Minter) UpdateApp(ctx *context.Context, user []byte, workId gtype.WorkId, app *gtype.TeeApp, envs []v1.EnvVar, version uint64) error {
	if uint8(app.Status) == 2 {
		m.StopApp(workId)
		return nil
	}

	saddress := AccountToAddress(user)
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
		existing.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = int32(app.Port[0])
		_, err = nameSpace.Update(*ctx, existing, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
	}

	return err
}
