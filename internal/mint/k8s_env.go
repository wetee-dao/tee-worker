package mint

import (
	"fmt"
	"strings"

	gtypes "github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	corev1 "k8s.io/api/core/v1"
	"wetee.app/worker/internal/store"
)

// Get Envs from Work
// 获取环境变量
func (m *Minter) BuildEnvs(podId uint64) ([]corev1.EnvVar, error) {
	// settings, err := m.GetSettingsFromWork(workId, nil)
	// if err != nil {
	// 	return []corev1.EnvVar{}, errors.Wrap(err, "GetSettingsFromWork error")
	// }
	settings := []*gtypes.Env1{}

	return m.BuildEnvsFromSettings(podId, settings)
}

// Build Envs
// 获取配置文件
func (m *Minter) BuildEnvsFromSettings(podId uint64, settings []*gtypes.Env1) ([]corev1.EnvVar, error) {
	// 用于应用联系控制面板的凭证
	wid, err := store.SealAppID(podId)
	if err != nil {
		return []corev1.EnvVar{}, err
	}

	envs := []corev1.EnvVar{
		{Name: "APPID", Value: wid},
		{Name: "PODID", Value: fmt.Sprint(podId)},
	}

	for _, setting := range settings {
		// TODO add file
		if setting.K.IsFile {
			continue
		}
		envs = append(envs, corev1.EnvVar{
			Name:  string(setting.K.AsEnvField0),
			Value: string(setting.V),
		})
	}

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

// Filter Envs for container
func filterEnvs(envs []*gtypes.Env1, index uint16) []*gtypes.Env1 {
	var fenvs []*gtypes.Env1
	for i, env := range envs {
		if env.Index == index {
			fenvs = append(fenvs, envs[i])
		}
	}
	return fenvs
}
