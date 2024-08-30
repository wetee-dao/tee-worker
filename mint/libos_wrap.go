package mint

import (
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func (m *Minter) WrapLibos(deployment *appsv1.Deployment, version *gtypes.TEEVersion) {
	if version.IsCVM {
		// 获取APPID
		envs := []corev1.EnvVar{
			deployment.Spec.Template.Spec.Containers[0].Env[0],
		}

		deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, corev1.Container{
			Name:  "libos",
			Image: "registry.cn-hangzhou.aliyuncs.com/wetee_dao/cvm:2024-08-25-14_38",
			Env:   envs,
		})
	}
}
