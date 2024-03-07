package mint

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/edgelesssys/ego/enclave"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"wetee.app/worker/store"
	"wetee.app/worker/util"
)

// 获取挖矿密钥
// GetKey get mint key
func GetMintKey() (*signature.KeyringPair, error) {
	key, err := store.GetMintId()

	var mss [32]byte
	if err != nil {
		// 前16位
		k, _, err := enclave.GetProductSealKey()
		if err != nil {
			util.LogWithRed("GetKey error", err)
			return nil, err
		}
		copy(mss[:], k)

		// 后16位随机数
		token := make([]byte, 16)
		rand.Read(token)
		copy(mss[16:], token)
	} else {
		copy(mss[:], key)
	}

	fmt.Println("GetKey", mss)
	store.SetMintId(mss[:])

	uri := hex.EncodeToString(mss[:])
	scheme := sr25519.Scheme{}
	kr, err := subkey.DeriveKeyPair(scheme, uri)
	if err != nil {
		return nil, err
	}

	return &signature.KeyringPair{
		URI:       uri,
		Address:   kr.SS58Address(42),
		PublicKey: kr.Public(),
	}, nil
}
