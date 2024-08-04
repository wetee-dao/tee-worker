package secret

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"wetee.app/worker/internal/store"

	"github.com/wetee-dao/go-sdk/pallet/types"
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

	err = store.SetSecrets(workId, &store.Secrets{
		Env: map[string]string{
			"": "",
		},
	})
	if err != nil {
		t.Error(err)
	}

	sigKey, err := sr25519.Scheme{}.Generate()
	if err != nil {
		t.Error(err)
	}

	param := &store.LoadParam{
		Address:   sigKey.SS58Address(42),
		Time:      fmt.Sprint(time.Now().Unix()),
		Signature: "NONE",
	}

	// 签名
	// Sign
	sig, err := sigKey.Sign([]byte(param.Time))
	if err != nil {
		t.Error(err)
	}
	param.Signature = hex.EncodeToString(sig)

	// _, err = loading(wid, param)
	// if err != nil {
	// 	t.Error(err)
	// }
}
