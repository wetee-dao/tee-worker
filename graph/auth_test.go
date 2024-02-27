package graph

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"wetee.app/worker/graph/model"
	"wetee.app/worker/store"
)

func TestDecodeToken(t *testing.T) {
	store.DBInit("bin/testdb")
	defer store.DBClose()

	prikey, err := sr25519.Scheme{}.Generate()
	if err != nil {
		t.Error(err)
	}

	account := &model.LoginContent{
		Address:   prikey.SS58Address(42),
		Timestamp: time.Now().Unix(),
	}
	bt, err := json.Marshal(account)
	if err != nil {
		t.Error(err)
	}

	sig, err := prikey.Sign([]byte("<Bytes>" + string(bt) + "</Bytes>"))
	if err != nil {
		t.Error(err)
	}

	accountstr := subkey.EncodeHex(bt)
	sigstr := subkey.EncodeHex(sig)

	token := accountstr + "||" + sigstr

	user := decodeToken(token)
	if user.Address != account.Address {
		t.Errorf("expected %s got %s", account.Address, user.Address)
	}
}
