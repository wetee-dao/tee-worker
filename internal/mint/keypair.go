package mint

import (
	"fmt"
	"os"

	"github.com/cometbft/cometbft/p2p"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// 获取挖矿密钥
// GetKey get mint key
func GetMintKey() (*chain.Signer, *model.PrivKey) {
	// init sidechain node key
	nodeKey, err := p2p.LoadNodeKey("./chain_data/config/node_key.json")
	if err != nil {
		fmt.Println("failed to load node key:", err)
		os.Exit(1)
	}

	privateKey, err := model.PrivateKeyFromOed25519(nodeKey.PrivKey.Bytes())
	if err != nil {
		fmt.Println("Marshal PKG_PK error:", err)
		os.Exit(1)
	}

	kr, err := privateKey.ToSigner()
	if err != nil {
		fmt.Println("ToSigner error:", err)
		os.Exit(1)
	}
	return kr, privateKey
}
