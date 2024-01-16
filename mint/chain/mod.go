package chain

import (
	"errors"

	gtypes "wetee.app/worker/mint/chain/gen/types"
)

func GetAccount(client *ChainClient, workID gtypes.WorkId) ([]byte, error) {
	if workID.Wtype.IsAPP {
		app := &App{
			Client: client,
			Signer: nil,
		}
		return app.GetAccount(workID.Id)
	}
	if workID.Wtype.IsTASK {
		task := &Task{
			Client: client,
			Signer: nil,
		}
		return task.GetAccount(workID.Id)
	}
	return nil, errors.New("unknow work type")
}
