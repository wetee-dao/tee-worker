package mint

import (
	"context"
	"fmt"
	"strings"

	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Minter) DeploymentPVCWrap(ctx *context.Context, nameSpace string, name string, cs []gtypes.Container, deployment *appsv1.Deployment) error {
	// 查询所有id所对应的pvc
	pvcList, err := m.K8sClient.CoreV1().PersistentVolumeClaims(nameSpace).List(*ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for cindex := range cs {
		disk := cs[cindex].Cr.Disk

		// 判断是否存在PVC
		l := len(disk)
		for i := 0; i < l; i++ {
			cdisk := disk[i]
			pvcName, pvc := findPvc(name, pvcList.Items, cindex, cdisk)
			// pvcName := name + "-pvc" + strings.ReplaceAll(path, "/", "-")
			// m.K8sClient
			// 不存在就创建
			if pvc == nil {
				pvc = &corev1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pvcName,
						Namespace: nameSpace,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("1Gi"),
							},
						},
					},
				}
				_, err := m.K8sClient.CoreV1().PersistentVolumeClaims(nameSpace).Create(*ctx, pvc, metav1.CreateOptions{})
				if err != nil {
					return err
				}
			}

			// 挂载卷到 deployment
			deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: name + "-store-" + fmt.Sprint(cindex) + "-" + fmt.Sprint(i),
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName,
					},
				},
			})

			// 挂载到容器
			deployment.Spec.Template.Spec.Containers[cindex].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[cindex].VolumeMounts, corev1.VolumeMount{
				Name:      name + "-store-" + fmt.Sprint(cindex) + "-" + fmt.Sprint(i),
				MountPath: string(cdisk.Path.AsSSDField0),
			})
		}
	}

	return nil
}

// 查询数组中是否存在目标元素
func findPvc(name string, arr []corev1.PersistentVolumeClaim, cindex int, target gtypes.Disk) (string, *corev1.PersistentVolumeClaim) {
	pvcName := name + "-pvc-" + fmt.Sprint(cindex) + "-" + strings.ReplaceAll(string(target.Path.AsSSDField0), "/", "-")
	for _, value := range arr {
		if value.ObjectMeta.Name == pvcName {
			return pvcName, &value
		}
	}
	return pvcName, nil
}
