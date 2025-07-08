package libos

import (
	"testing"

	"github.com/wetee-dao/tee-dsecret/pkg/chains/pallets/generated/types"
	"wetee.app/worker/internal/store"
)

func TestLoading(t *testing.T) {
	store.DBInit()
	defer store.DBClose()

	workId := types.WorkId{
		Id: 1,
		Wtype: types.WorkType{
			IsAPP: true,
		},
	}
	_, err := store.SealAppID(workId)
	if err != nil {
		t.Error(err)
	}
}
