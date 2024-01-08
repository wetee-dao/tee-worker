package chain

import (
	"github.com/edgelesssys/ego/enclave"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"wetee.app/worker/internal/util"
)

func GetMintKey() (subkey.KeyPair, error) {
	k, info, err := enclave.GetProductSealKey()
	if err != nil {
		util.LogWithRed("GetKey error", err)
		return nil, err
	}
	util.LogWithRed("GetKey", k)
	util.LogWithRed("GetKeyInfo", info)

	var mss [32]byte
	copy(mss[:], k)
	return sr25519.Scheme{}.FromSeed(mss[:])
}
