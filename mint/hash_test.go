package mint

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"golang.org/x/crypto/blake2b"
	"wetee.app/worker/dao"
)

func TestGetWorkLogHash(t *testing.T) {
	dao.DBInit("bin/testdb")
	defer dao.DBClose()

	var blockNumber uint64 = 0
	logs := []string{
		"2021-01-01 12:00:00",
	}
	logHash, err := getWorkLogHash("test", logs, blockNumber)
	if err != nil {
		t.Error(err)
	}
	hashHex := hex.EncodeToString(logHash)

	pf := WorkLogProof{
		BlockNumber: blockNumber,
		Time:        uint64(time.Now().Unix()),
		Logs:        logs,
	}
	bt, _ := json.Marshal(&pf)

	h := blake2b.Sum256(bt)
	hashString := hex.EncodeToString(h[:])

	if hashString != hashHex {
		t.Error("hash not equal")
	}
}

func TestGetWorkCrHash(t *testing.T) {
	dao.DBInit("bin/testdb")
	defer dao.DBClose()

	cr := map[string][]int64{
		"test": {1, 2, 3},
	}
	blockNumber := uint64(0)
	crHash, _, err := getWorkCrHash("test", cr, blockNumber)
	if err != nil {
		t.Error(err)
	}
	hashHex := hex.EncodeToString(crHash)

	pf := WorkCrProof{
		Cr:          cr,
		BlockNumber: blockNumber,
		Time:        uint64(time.Now().Unix()),
	}

	bt, _ := json.Marshal(&pf)
	h := blake2b.Sum256(bt)

	hashString := hex.EncodeToString(h[:])

	if hashString != hashHex {
		t.Error("hash not equal")
	}
}
