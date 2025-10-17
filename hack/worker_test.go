package hack

// import (
// 	"math/big"

// 	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
// 	chain "github.com/wetee-dao/ink.go"
// 	"github.com/wetee-dao/ink.go/util"
// 	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts"
// 	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/subnet"
// )

// func TestUpdateWorker() {
// 	client, err := chain.InitClient([]string{TestChainUrl}, true)
// 	if err != nil {
// 		panic(err)
// 	}

// 	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
// 	if err != nil {
// 		util.LogWithPurple("Sr25519PairFromSecret", err)
// 		panic(err)
// 	}

// 	_call := chain.ExecParams{
// 		Signer:    &pk,
// 		PayAmount: types.NewU128(*big.NewInt(0)),
// 	}

// 	subnetContract, err := subnet.InitSubnetContract(client, contracts.GetSubnetAddress())
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = subnetContract.ExecWorkerUpdate(0, []byte("w0"), subnet.Ip{}, 10000, _call)
// }
