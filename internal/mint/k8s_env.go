package mint

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	corev1 "k8s.io/api/core/v1"
	"wetee.app/worker/internal/store"
)

// Build Envs
// 获取配置文件
func (m *Minter) BuildEnvsFromSettings(podId uint64, nameSpace string, settings []model.Env) ([]corev1.EnvVar, error) {
	// 用于应用联系控制面板的凭证
	wid, err := store.SealAppID(podId)
	if err != nil {
		return []corev1.EnvVar{}, err
	}

	envs := []corev1.EnvVar{
		{Name: "APPID", Value: wid},
		{Name: "PODID", Value: fmt.Sprint(podId)},
		{Name: "NAME_SPACE", Value: nameSpace},
	}

	files := map[string]string{}
	encrypts := map[string]uint64{}
	for _, setting := range settings {
		if setting.Env != nil {
			envs = append(envs, corev1.EnvVar{
				Name:  string(setting.Env.F0),
				Value: string(setting.Env.F1),
			})
		} else if setting.File != nil {
			files[string(setting.File.F0)] = hex.EncodeToString(setting.File.F1)
		} else if setting.Encrypt != nil {
			encrypts[string(setting.Encrypt.F0)] = setting.Encrypt.F1
		}
	}

	fbt, _ := json.Marshal(files)
	envs = append(envs, corev1.EnvVar{
		Name:  "__FILES__",
		Value: string(fbt),
	})

	ebt, _ := json.Marshal(encrypts)
	envs = append(envs, corev1.EnvVar{
		Name:  "__ENCRYPTS__",
		Value: string(ebt),
	})

	return envs, nil
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
