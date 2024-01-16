package chain

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/pkg/errors"
	"wetee.app/worker/mint/chain/gen/types"
	"wetee.app/worker/mint/chain/gen/weteeapp"
)

// Worker
type App struct {
	Client *ChainClient
	Signer *signature.KeyringPair
}

func (w *App) GetApp(publey []byte, id uint64) (*types.TeeApp, error) {
	if len(publey) != 32 {
		return nil, errors.New("publey length error")
	}

	var mss [32]byte
	copy(mss[:], publey)

	res, ok, err := weteeapp.GetTEEAppsLatest(w.Client.Api.RPC.State, mss, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetK8sClusterAccountsLatest => not ok")
	}
	return &res, nil
}

func (w *App) GetAccount(id uint64) ([]byte, error) {
	res, ok, err := weteeapp.GetAppIdAccountsLatest(w.Client.Api.RPC.State, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("GetAppIdAccountsLatest => not ok")
	}
	return res[:], nil
}
