package mint

import (
	"testing"

	"github.com/vedhavyas/go-subkey/v2"
)

func TestAccountToSpace(t *testing.T) {
	account := "5EYCAe5k8AykupEcF5NFAdyUJa2TiYCrRG7NzJy9qiXkHxdb"
	_, pubkeyBytes, _ := subkey.SS58Decode(account)
	addr := AccountToSpace(pubkeyBytes)

	if addr != "dlnm8r3nclq6apb4c5ng00000000000008" {
		t.Error("AccountToSpace failed")
	}
}
