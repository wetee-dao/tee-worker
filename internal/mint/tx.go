package mint

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"wetee.app/worker/internal/mint/proof"
	"wetee.app/worker/internal/store"
)

func (m *Minter) AddPendingDeployTx(podId uint64) error {
	pod, err := store.GetPod(podId)
	if err != nil {
		return errors.Wrap(err, "AddPendingDeployTx GetPod")
	}

	// get namespace name
	address := AccountToSpace(pod.Owner[:])

	// Get logs and use compute resource
	now := time.Now()
	ctx := context.Background()
	logs, crs, err := m.GetMetric(&ctx, address, *pod, now, 1, false)
	if err != nil {
		return errors.Wrap(err, "AddPendingDeployTx GetLogAndCr")
	}

	// make pod start proof
	call, index, err := proof.MakeWorkProof(*pod, logs, crs, now)
	if err != nil {
		return errors.Wrap(err, "AddPendingDeployTx MakeWorkProof")
	}
	indexCall := model.ToIndexCall(call, index)
	return store.AddPendingCall(indexCall)
}
