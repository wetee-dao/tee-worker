package libos

import (
	"testing"

	"wetee.app/worker/internal/store"

	"wetee.app/dsecret/type/pallet/types"
)

func TestLoading(t *testing.T) {
	store.DBInit("bin/testdb")
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
