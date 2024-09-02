package proof

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"golang.org/x/crypto/blake2b"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/util"
)

// TestListLogsById tests the ListLogsById function
func TestListLogsById(t *testing.T) {
	store.DBInit("bin/testdb")
	defer store.DBClose()

	// Test case 1: Valid input
	id := gtypes.WorkId{Id: uint64(time.Now().Unix())}
	page := 1
	size := 2

	_, bt, err := GetWorkLogHash([]string{"log1", "xlog2"}, 1)
	if err != nil {
		t.Errorf("GetWorkLogHash Expected no error, got %v", err)
	}

	name := util.GetWorkTypeStr(id) + "-" + fmt.Sprint(id.Id)
	store.AddToList(LogBucket, name, bt)

	logs, err := ListLogsById(id, page, size, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("GetWorkLogHash Expected 1, got %v", len(logs))
	}

	// Test case 2: Error case
	page = 2
	size = 1
	logs, err = ListLogsById(id, page, size, false)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if logs != nil {
		t.Errorf("Expected nil log list, got %+v", logs)
	}
}

func TestGetWorkLogHash(t *testing.T) {
	store.DBInit("bin/testdb")
	defer store.DBClose()

	var blockNumber uint64 = 1
	logs := []string{"log1", "log2"}

	hash2, _, err := GetWorkLogHash(logs, blockNumber)
	if err != nil {
		t.Errorf("GetWorkLogHash Expected no error, got %v", err)
	}

	pf := WorkLogProof{
		BlockNumber: blockNumber,
		Time:        uint64(time.Now().Unix()),
		Logs:        logs,
	}
	bt, err := json.Marshal(&pf)
	hash := blake2b.Sum256(bt)

	if err != nil {
		t.Errorf("GetWorkLogHash Expected no error, got %v", err)
	}

	if hex.EncodeToString(hash[:]) != hex.EncodeToString(hash2) {
		t.Errorf("GetWorkLogHash Expected %v, got %v", hash, hash2)
	}
}
