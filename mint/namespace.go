package mint

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Minter) checkNameSpace(ctx context.Context, address string) error {
	k8s := m.K8sClient.CoreV1()
	nameSpaces, err := k8s.Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	var found bool = false
	for _, namespace := range nameSpaces.Items {
		if namespace.Name == address {
			found = true
			break
		}
	}
	if !found {
		_, err = k8s.Namespaces().Create(ctx, &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: address,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
