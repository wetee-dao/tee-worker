package mint

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/util"
)

// Get Envs from Work
// 获取环境变量
func (m *Minter) BuildEnvs(workId gtypes.WorkId) ([]corev1.EnvVar, error) {
	settings, err := m.GetSettingsFromWork(workId, nil)
	if err != nil {
		return []corev1.EnvVar{}, errors.Wrap(err, "GetSettingsFromWork error")
	}

	return m.BuildEnvsFromSettings(workId, settings)
}

// Build Envs
// 获取配置文件
func (m *Minter) BuildEnvsFromSettings(workId gtypes.WorkId, settings []*gtypes.Env1) ([]corev1.EnvVar, error) {
	// 用于应用联系控制面板的凭证
	wid, err := store.SealAppID(workId)
	if err != nil {
		return []corev1.EnvVar{}, err
	}

	chainUrl := DefaultChainUrl
	url, err := store.GetChainUrl()
	if err == nil {
		chainUrl = url
	}

	envs := []corev1.EnvVar{
		{Name: "APPID", Value: wid},
		{Name: "CHAIN_ADDR", Value: chainUrl},
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
func (m *Minter) WrapEnvs(envs []corev1.EnvVar, nameSpace, name string, nodeSers *v1.Service) error {
	mdata := make(map[string]string)
	mdata["cluster_domain"] = m.HostDomain
	mdata["project_domain"] = nameSpace + ".svc.cluster.local"
	mdata["gen_ssl"] = strings.Join(util.GetSslRoot(), "|")
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

func filterEnvs(envs []*gtypes.Env1, index uint16) []*gtypes.Env1 {
	var fenvs []*gtypes.Env1
	for i, env := range envs {
		if env.Index == index {
			fenvs = append(fenvs, envs[i])
		}
	}
	return fenvs
}
