package chain

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/pkg/errors"
	"wetee.app/worker/mint/chain/gen/weteetask"
)

// Worker
type Task struct {
	Client *ChainClient
	Signer *signature.KeyringPair
}

func (w *Task) GetAccount(id uint64) ([]byte, error) {
	res, ok, err := weteetask.GetTaskIdAccountsLatest(w.Client.Api.RPC.State, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetAppIdAccountsLatest => not ok")
	}
	return res[:], nil
}
