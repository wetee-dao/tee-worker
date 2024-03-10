package mint

import (
	"testing"

	"wetee.app/worker/store"
)

func TestGetMintKey(t *testing.T) {
	store.DBInit("bin/testdb")
	defer store.DBClose()

	k, err := GetMintKey()
	if err != nil {
		t.Error(err)
	}
	t.Log(k)
}
