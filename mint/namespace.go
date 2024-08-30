package mint

import (
	"context"
	"encoding/base32"
	"encoding/hex"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// checkNameSpace check if the namespace exists, if not, create it
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

// Account To Hex Address
// 将用户公钥转换为hex地址
func AccountToSpace(user []byte) string {
	address := base32.HexEncoding.EncodeToString(user[:])
	address = strings.ReplaceAll(strings.ToLower(address), "=", "")
	return strings.TrimRight(address, "000000000000000000")
}

// Hex Address To Account
// 将hex地址转换为用户公钥
func HexStringToSpace(address string) string {
	address = strings.ReplaceAll(address, "0x", "")
	user, _ := hex.DecodeString(address)
	return AccountToSpace(user)
}
