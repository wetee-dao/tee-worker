package types

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/go-sdk/core"
)

type TeeTrigger struct {
	Tee       TeeParam
	ClusterId uint64
	Callids   []types.U128
	Sig       []byte
}

func (t *TeeTrigger) Sign(signer *core.Signer) error {
	msg := fmt.Sprint(t.ClusterId, t.Callids)
	sig, err := signer.Sign([]byte(msg))
	if err != nil {
		return err
	}
	t.Sig = sig
	return nil
}
