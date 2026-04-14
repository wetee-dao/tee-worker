package hack

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	contracts "github.com/wetee-dao/tee-dsecret/pkg/chains/revives"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/revives/cloud"
)

func TestInitDisk(t *testing.T) {
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

	err = cloudIns.ExecCreateDisk([]byte("TEST"), 10, ink.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
	fmt.Println(err)

	list, _, err := cloudIns.QueryUserSecrets(pk.H160Address(), util.NewNone[uint64](), 100, ink.DefaultParamWithOrigin(pk.AccountID()))
	fmt.Println(list)
	fmt.Println(err)

}

func TestQueryDisk(t *testing.T) {
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

	list, _, err := cloudIns.QueryUserDisks(pk.H160Address(), util.NewNone[uint64](), 100, ink.DefaultParamWithOrigin(pk.AccountID()))
	fmt.Println(list)
	fmt.Println(err)
}
