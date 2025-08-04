package hack

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/cloud"
)

func TestInitSecret(t *testing.T) {
	client, err := ink.InitClient([]string{"ws://127.0.0.1:9944"}, true)
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

	err = cloudIns.ExecInitSecret([]byte("TEST"), ink.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
	fmt.Println(err)

	list, _, err := cloudIns.QueryUserSecrets(pk.H160Address(), util.NewNone[uint64](), 100, ink.DefaultParamWithOrigin(pk.AccountID()))
	fmt.Println(list)
	fmt.Println(err)
}

func TestQuerySecret(t *testing.T) {
	client, err := ink.InitClient([]string{"ws://127.0.0.1:9944"}, true)
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

	list, _, err := cloudIns.QueryUserSecrets(pk.H160Address(), util.NewNone[uint64](), 100, ink.DefaultParamWithOrigin(pk.AccountID()))
	fmt.Println(list)
	fmt.Println(err)
}
