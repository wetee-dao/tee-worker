package mint

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/vedhavyas/go-subkey/sr25519"

	"github.com/wetee-dao/go-sdk/gen/types"
	"wetee.app/worker/dao"
)

func TestLoading(t *testing.T) {
	dao.DBInit("bin/testdb")
	defer dao.DBClose()

	workId := types.WorkId{
		Id: 1,
		Wtype: types.WorkType{
			IsAPP: true,
		},
	}
	wid, err := dao.SealAppID(workId)
	if err != nil {
		t.Error(err)
	}

	err = dao.SetSecrets(workId, &dao.Secrets{
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

	param := &dao.LoadParam{
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

	_, err = loading(wid, param)
	if err != nil {
		t.Error(err)
	}
}
