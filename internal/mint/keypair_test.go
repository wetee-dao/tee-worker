package mint

import (
	"testing"

	"wetee.app/worker/internal/store"
)

func TestGetMintKey(t *testing.T) {
	store.DBInit()
	defer store.DBClose()

	k, _ := GetMintKey()
	t.Log(k)
}
