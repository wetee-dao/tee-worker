package mint

import (
	"context"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	inkutil "github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"wetee.app/worker/internal/mint/proof"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/internal/util"
)

func (m *Minter) AddPendingDeployTx(podId uint64, key []byte) error {
	pod, err := store.GetPod(podId)
	if err != nil || pod == nil {
		return errors.New("POD not found")
	}

	// get namespace name
	address := AccountToSpace(pod.Owner[:])

	// Get logs and use compute resource
	now := time.Now()
	ctx := context.Background()
	logs, crs, err := m.GetMetric(&ctx, address, *pod, now, 1, false)
	if err != nil {
		return errors.Wrap(err, "GetLogAndCr")
	}

	ss58 := model.SS58Encode(key, 42)
	account, _ := types.NewAccountID(key)

	util.LogWithBlue("LOADED POD", podId, "address =>", ss58)

	// make pod start proof
	call, index, err := proof.MakeWorkProof(*pod, inkutil.NewSome(*account), logs, crs, now)
	if err != nil {
		return errors.Wrap(err, "MakeWorkProof")
	}

	indexCall := model.ToIndexCall(call, index)
	return store.AddPendingCall(indexCall)
}
