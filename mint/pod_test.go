package mint

import "testing"

func TestAccountToSpace(t *testing.T) {
	addr := AccountToSpace([]byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxx"))
	if addr != "8787878787878787878787878787878787878787878787878787878" {
		t.Error("AccountToSpace failed")
	}
}
