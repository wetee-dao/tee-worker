package mint

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/wetee-dao/go-sdk/core"
	"go.dedis.ch/kyber/v3/suites"
	"wetee.app/worker/internal/store"
	types "wetee.app/worker/type"
)

// 获取挖矿密钥
// GetKey get mint key
func GetMintKey() (*core.Signer, *types.PrivKey, error) {
	key, err := store.GetMintId()
	var mss []byte
	var privateKey *types.PrivKey
	if err != nil {
		suite := suites.MustFind("Ed25519")
		privateKey, _, err = types.GenerateKeyPair(suite, rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		bt, err := hex.DecodeString(privateKey.String())
		if err != nil {
			return nil, nil, err
		}
		mss = bt
	} else {
		keyString := hex.EncodeToString(key)
		privateKey, err = types.PrivateKeyFromLibp2pHex(keyString)
		if err != nil {
			fmt.Println("Marshal PKG_PK error:", err)
			return nil, nil, err
		}
		mss = key
	}

	store.SetMintId(mss)

	bt, err := privateKey.Raw()
	if err != nil {
		return nil, nil, err
	}

	var ed25519Key ed25519.PrivateKey = bt
	kr, err := core.Ed25519PairFromPk(ed25519Key, 42)
	if err != nil {
		return nil, nil, err
	}

	return &kr, privateKey, nil
}
