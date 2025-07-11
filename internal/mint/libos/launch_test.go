package libos

import (
	"testing"

	"wetee.app/worker/internal/store"
)

func TestLoading(t *testing.T) {
	store.DBInit()
	defer store.DBClose()

	_, err := store.SealAppID(1)
	if err != nil {
		t.Error(err)
	}
}
