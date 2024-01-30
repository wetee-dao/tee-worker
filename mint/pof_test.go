package mint

import "testing"

func TestAccountToAddress(t *testing.T) {
	addr := AccountToAddress([]byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxx"))
	if addr != "8787878787878787878787878787878787878787878787878787878" {
		t.Error("AccountToAddress failed")
	}
}
