package hack

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	contracts "github.com/wetee-dao/tee-dsecret/pkg/chains/revive"
	"github.com/wetee-dao/tee-dsecret/pkg/chains/revive/cloud"
	"wetee.app/worker/internal/util"
)

const TestChainUrl = "wss://xiaobai.asyou.me:30001/ws"

func TestAddPod(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
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
	cloudIns.ExecCreatePod(
		[]byte("Test"),
		cloud.PodType{CPU: &trueValue},
		cloud.TEEType{SGX: &trueValue},
		[]cloud.Container{
			{
				Image:   []byte("wetee/ego-hello:2025-07-28-06-58"),
				Command: cloud.Command{NONE: &trueValue},
				Port:    []cloud.Service{{Http: &value80}},
				Cpu:     1000,
				Mem:     800,
				Env: []cloud.Env{
					{
						Env: &struct {
							F0 []byte
							F1 []byte
						}{
							F0: []byte("K"),
							F1: []byte("V"),
						},
					},
					{
						File: &struct {
							F0 []byte
							F1 []byte
						}{
							F0: []byte("K2"),
							F1: []byte("K3"),
						},
					},
					{
						Encrypt: &struct {
							F0 []byte
							F1 uint64
						}{
							F0: []byte("K2"),
							F1: 0,
						},
					},
				},
			},
		},
		0,
		1,
		0,
		0,
		chain.ExecParams{
			Signer:    &pk,
			PayAmount: types.NewU128(*big.NewInt(0)),
		},
	)
}

func TestUpdatePod(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
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
	var container_id uint64 = 0
	err = cloudIns.ExecEditContainer(0,
		[]cloud.ContainerInput{
			{
				Etype: cloud.EditType{UPDATE: &container_id},
				Container: cloud.Container{
					Image:   []byte("wetee/ego-hello:2025-07-28-06-58"),
					Command: cloud.Command{NONE: &trueValue},
					Port:    []cloud.Service{{Http: &value80}},
					Cpu:     1000,
					Mem:     800,
					Env: []cloud.Env{
						{
							Env: &struct {
								F0 []byte
								F1 []byte
							}{
								F0: []byte("K"),
								F1: []byte("V"),
							},
						},
						{
							File: &struct {
								F0 []byte
								F1 []byte
							}{
								F0: []byte("K2"),
								F1: []byte("K3"),
							},
						},
						{
							Encrypt: &struct {
								F0 []byte
								F1 uint64
							}{
								F0: []byte("K2"),
								F1: 0,
							},
						},
					},
				},
			},
		},
		chain.ExecParams{
			Signer:    &pk,
			PayAmount: types.NewU128(*big.NewInt(0)),
		},
	)

	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		panic(err)
	}
}

func TestRestartPod(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
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

	cloudIns.ExecRestartPod(0, chain.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}

func TestDeletePod(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
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

	cloudIns.ExecStopPod(0, chain.ExecParams{
		Signer:    &pk,
		PayAmount: types.NewU128(*big.NewInt(0)),
	})
}

func TestQueryPod(t *testing.T) {
	client, err := chain.InitClient([]string{TestChainUrl}, true)
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

	pods, _, err := cloudIns.QueryPodsByIds([]uint64{0}, chain.DefaultParamWithOrigin(pk.AccountID()))
	fmt.Println((*pods)[0].F5, err)
}
