package mint

import (
	"github.com/pkg/errors"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"wetee.app/worker/internal/store"
)

func (m *Minter) AddPendingDeployTx(call *model.TeeCall) error {
	pod, err := store.GetPod(call.GetPodStart().Id)
	if err != nil || pod == nil {
		return errors.New("POD not found")
	}

	return store.AddPendingCall(call)
}
