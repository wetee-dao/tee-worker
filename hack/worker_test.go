package hack

import (
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	contracts "github.com/wetee-dao/tee-dsecret/pkg/chains/ink"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/ink/cloud"
)

func TestUpdatePodCode(t *testing.T) {
	client, err := ink.InitClient([]string{TestChainUrl}, true)
	if err != nil {
		panic(err)
	}

	pk, err := ink.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	cloudIns, err := cloud.InitCloudContract(client, contracts.GetCloudAddress())
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	err = cloudIns.ExecCharge(ink.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(100000000000)),
	})
	if err != nil {
		panic(err)
	}
}
