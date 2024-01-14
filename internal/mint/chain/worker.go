package chain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"wetee.app/worker/internal/mint/chain/gen/weteeworker"

	gtypes "wetee.app/worker/internal/mint/chain/gen/types"
)

// Worker
type Worker struct {
	Client *ChainClient
	Signer *signature.KeyringPair
}

// 集群注册
// Cluster register
func (w *Worker) ClusterRegister(name string, ip []uint8, port uint32, level uint8) error {
	runtimeCall := weteeworker.MakeClusterRegisterCall(
		[]byte(name),
		[]gtypes.Ip{
			{
				Ipv4: gtypes.OptionTUint32{IsNone: false, IsSome: true, AsSomeField0: 100},
				Ipv6: gtypes.OptionTU128{IsNone: false, IsSome: true, AsSomeField0: types.NewU128(*big.NewInt(100))},
			},
		},
		port,
		level,
	)

	return w.Client.SignAndSubmit(w.Signer, runtimeCall)
}

// 集群抵押
// Cluster mortgage
func (w *Worker) ClusterMortgage(id uint64, cpu uint16, mem uint16, disk uint16, deposit uint64) error {
	d := big.NewInt(0)
	d.SetUint64(deposit)
	runtimeCall := weteeworker.MakeClusterMortgageCall(
		id,
		cpu,
		mem,
		disk,
		types.UCompact(*d),
	)

	return w.Client.SignAndSubmit(w.Signer, runtimeCall)
}

func (w *Worker) ClusterWithdrawal(id uint64, val int64) error {
	runtimeCall := weteeworker.MakeClusterWithdrawalCall(
		gtypes.WorkId{
			Wtype: gtypes.WorkType{IsAPP: true, IsTASK: false},
			Id:    id,
		},
		types.NewU128(*big.NewInt(val)),
	)

	return w.Client.SignAndSubmit(w.Signer, runtimeCall)
}

func (w *Worker) ClusterUnmortgage(clusterID uint64, id uint64) error {
	runtimeCall := weteeworker.MakeClusterUnmortgageCall(
		clusterID,
		id,
	)

	return w.Client.SignAndSubmit(w.Signer, runtimeCall)
}

func (w *Worker) ClusterStop(clusterID uint64) error {
	runtimeCall := weteeworker.MakeClusterStopCall(
		clusterID,
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

func (w *Worker) GetClusterContracts(clusterID uint64) ([]gtypes.ContractState, error) {
	set, err := w.Client.QueryDoubleMapAll("WeteeWorker", "ClusterContracts", clusterID)
	if err != nil {
		return nil, err
	}

	var list []gtypes.ContractState = make([]gtypes.ContractState, 0, len(set))
	for _, elem := range set {
		for _, change := range elem.Changes {
			var cs gtypes.ContractState
			if err := codec.Decode(change.StorageData, &cs); err != nil {
				fmt.Println(err)
				continue
			}
			list = append(list, cs)
		}
	}

	fmt.Println(err)
	return list, nil
}
