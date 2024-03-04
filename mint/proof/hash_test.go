package proof

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"golang.org/x/crypto/blake2b"
	"wetee.app/worker/store"
)

func TestGetWorkLogHash(t *testing.T) {
	store.DBInit("bin/testdb")
	defer store.DBClose()

	var blockNumber uint64 = 0
	logs := []string{
		"2021-01-01 12:00:00",
	}
	logHash, err := GetWorkLogHash("test", logs, blockNumber)
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
	store.DBInit("bin/testdb")
	defer store.DBClose()

	cr := map[string][]int64{
		"test": {1, 2, 3},
	}
	blockNumber := uint64(0)
	crHash, _, err := GetWorkCrHash("test", cr, blockNumber)
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
