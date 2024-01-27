package mint

import (
	"encoding/hex"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/edgelesssys/ego/enclave"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"wetee.app/worker/util"
)

// 获取挖矿密钥
// GetKey get mint key
func GetMintKey() (*signature.KeyringPair, error) {
	k, _, err := enclave.GetProductSealKey()
	if err != nil {
		util.LogWithRed("GetKey error", err)
		return nil, err
	}
	util.LogWithRed("GetKey", k)
	// LogWithRed("GetKeyInfo", info)

	var mss [32]byte
	copy(mss[:], k)

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
