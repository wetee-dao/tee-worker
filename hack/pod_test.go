package hack

import (
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/contracts/cloud"
	"wetee.app/worker/internal/util"
)

func TestAddPod(t *testing.T) {
	client, err := chain.ClientInit("ws://127.0.0.1:9944", true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	cloudIns, err := cloud.InitCloudContract(client, contracts.GetCloudAddress())
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	trueValue := true
	var value80 uint16 = 80
	cloudIns.CallCreatePod(
		[]byte("Test"),
		cloud.PodType{CpuService: &trueValue},
		cloud.TEEType{SGX: &trueValue},
		[]cloud.Container{
			{
				Image:   []byte("nginx"),
				Command: cloud.Command{NONE: &trueValue},
				Port:    []cloud.Service{{Http: &value80}},
				Cr:      cloud.CR{Cpu: 1000, Mem: 800},
			},
		},
		0,
		1,
		0,
		chain.CallParams{
			Signer:    &pk,
			PayAmount: types.NewU128(*big.NewInt(0)),
		},
	)
}

func TestUpdatePod(t *testing.T) {
	client, err := chain.ClientInit("ws://127.0.0.1:9944", true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	cloudIns, err := cloud.InitCloudContract(client, contracts.GetCloudAddress())
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	trueValue := true
	var value80 uint16 = 80
	cloudIns.CallUpdateContainer(0, 0,
		cloud.Container{
			Image:   []byte("nginx"),
			Command: cloud.Command{NONE: &trueValue},
			Port:    []cloud.Service{{Http: &value80}},
			Cr:      cloud.CR{Cpu: 1000, Mem: 800},
		},
		chain.CallParams{
			Signer:    &pk,
			PayAmount: types.NewU128(*big.NewInt(0)),
		},
	)
}

func TestDeletePod(t *testing.T) {
	client, err := chain.ClientInit("ws://127.0.0.1:9944", true)
	if err != nil {
		panic(err)
	}

	pk, err := chain.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	cloudIns, err := cloud.InitCloudContract(client, contracts.GetCloudAddress())
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}

	cloudIns.CallStopPod(2, chain.CallParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}
