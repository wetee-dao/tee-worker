package mint

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type InitData struct {
	Envs     map[string]string
	Files    map[int]map[string]string
	Encrypts map[int]map[string]uint64
}

// Build Envs
// 获取配置文件
func (m *Minter) BuildEnvsFromSettings(podId uint64, nameSpace string, index int, settingEnvs []model.Env, init *InitData) ([]corev1.EnvVar, error) {
	envs := []corev1.EnvVar{}

	if init.Envs == nil {
		init.Envs = make(map[string]string)
	}
	if init.Files == nil {
		init.Files = make(map[int]map[string]string)
	}
	if init.Files[index] == nil {
		init.Files[index] = make(map[string]string)
	}
	if init.Encrypts == nil {
		init.Encrypts = make(map[int]map[string]uint64)
	}
	if init.Encrypts[index] == nil {
		init.Encrypts[index] = make(map[string]uint64)
	}

	for _, setting := range settingEnvs {
		if setting.Env != nil {
			envs = append(envs, corev1.EnvVar{
				Name:  string(setting.Env.F0),
				Value: string(setting.Env.F1),
			})
		} else if setting.File != nil {
			init.Files[index][string(setting.File.F0)] = hex.EncodeToString(setting.File.F1)
		} else if setting.Encrypt != nil {
			init.Encrypts[index][string(setting.Encrypt.F0)] = setting.Encrypt.F1
		}
	}

	// fbt, _ := json.Marshal(files)
	// envs = append(envs, corev1.EnvVar{
	// 	Name:  "__FILES__",
	// 	Value: string(fbt),
	// })

	// ebt, _ := json.Marshal(encrypts)
	// envs = append(envs, corev1.EnvVar{
	// 	Name:  "__ENCRYPTS__",
	// 	Value: string(ebt),
	// })

	return envs, nil
}

func (m *Minter) WrapDeploymentInitData(deployment *appsv1.Deployment, version model.TEEType, initData *InitData) {
	if version.SGX != nil {
		for k, v := range initData.Envs {
			deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
				Name:  k,
				Value: v,
			})
		}

		encrypts, _ := json.Marshal(initData.Encrypts)
		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "__ENCRYPTS__",
			Value: string(encrypts),
		})

		files, _ := json.Marshal(initData.Files)
		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "__FILES__",
			Value: string(files),
		})
	} else if version.CVM != nil {

	}
}

func (m *Minter) WrapPodInitData(pod *corev1.Pod, version model.TEEType, initData *InitData) {
	for k, v := range initData.Envs {
		pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	encrypts, _ := json.Marshal(initData.Encrypts)
	pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
		Name:  "__ENCRYPTS__",
		Value: string(encrypts),
	})

	files, _ := json.Marshal(initData.Files)
	pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
		Name:  "__FILES__",
		Value: string(files),
	})
}

// WrapNodeService
// 包装环境变量
func (m *Minter) WrapEnvs(envs []corev1.EnvVar, nameSpace, name string, nodeSers *corev1.Service) error {
	mdata := make(map[string]string)
	mdata["cluster_domain"] = m.HostDomain
	mdata["project_domain"] = nameSpace + ".svc.cluster.local"
	// mdata["server_ssl"] = strings.Join(util.GetSslRoot(), "|")
	for i, port := range nodeSers.Spec.Ports {
		if port.NodePort != 0 {
			mdata["ser_"+fmt.Sprint(i)+"_nodeport"] = fmt.Sprint(port.NodePort)
		}
	}

	for i, env := range envs {
		if strings.Contains(env.Value, "{{.") {
			v, err := renderTemplate(env.Value, mdata)
			if err != nil {
				return err
			}
			envs[i].Value = v
		}
	}

	return nil
}
