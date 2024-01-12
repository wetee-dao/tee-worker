package chain

import (
	"errors"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"wetee.app/worker/internal/mint/chain/gen/weteeworker"

	gtypes "wetee.app/worker/internal/mint/chain/gen/types"
)

// Worker
type Worker struct {
	Client *ChainClient
	Signer *signature.KeyringPair
}

// 集群注册
// ClusterRegister
func (w *Worker) ClusterRegister(ip []uint8) error {
	runtimeCall := weteeworker.MakeClusterRegisterCall(
		[]byte("test"),
		[]gtypes.Ip{
			{
				Ipv4: gtypes.OptionTUint32{IsNone: false, IsSome: true, AsSomeField0: 100},
				Ipv6: gtypes.OptionTU128{IsNone: false, IsSome: true, AsSomeField0: types.NewU128(*big.NewInt(100))},
			},
		},
		100,
		1,
	)

	return w.Client.SignAndSubmit(w.Signer, runtimeCall)
}

func (w *Worker) ClusterMortgage() error {
	runtimeCall := weteeworker.MakeClusterMortgageCall(
		1,
		100,
		100,
		100,
		types.UCompact(*big.NewInt(100)),
	)

	return w.Client.SignAndSubmit(w.Signer, runtimeCall)
}

func (w *Worker) ClusterWithdrawal() error {
	runtimeCall := weteeworker.MakeClusterWithdrawalCall(
		gtypes.WorkId{
			Wtype: gtypes.WorkType{IsAPP: true, IsTASK: false},
		},
		types.NewU128(*big.NewInt(100)),
	)

	return w.Client.SignAndSubmit(w.Signer, runtimeCall)
}

func (w *Worker) ClusterUnmortgage() error {
	runtimeCall := weteeworker.MakeClusterUnmortgageCall(
		1,
		100,
	)

	return w.Client.SignAndSubmit(w.Signer, runtimeCall)
}

func (w *Worker) ClusterStop() error {
	runtimeCall := weteeworker.MakeClusterStopCall(
		1,
	)

	return w.Client.SignAndSubmit(w.Signer, runtimeCall)
}

func (w *Worker) Getk8sClusterAccounts(publey []byte) (uint64, error) {
	if len(publey) != 32 {
		return 0, errors.New("publey length error")
	}

	var mss [32]byte
	copy(mss[:], publey)

	res, ok, err := weteeworker.GetK8sClusterAccountsLatest(w.Client.Api.RPC.State, mss)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("GetK8sClusterAccountsLatest => not ok")
	}
	return res, nil
}

func (w *Worker) GetClusterContracts() error {
	// res, ok, err := weteeworker.GetClusterContractsLatest()
	return nil
}
