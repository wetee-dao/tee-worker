package mint

import "testing"

func TestGetMintKey(t *testing.T) {
	k, err := GetMintKey()
	if err != nil {
		t.Error(err)
	}
	t.Log(k)
}
