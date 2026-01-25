package proof

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"golang.org/x/crypto/blake2b"
	"wetee.app/worker/internal/store"
)

// TestListLogsById tests the ListLogsById function
func TestListLogsById(t *testing.T) {
	store.DBInit()
	defer store.DBClose()

	// Test case 1: Valid input
	id := uint64(time.Now().Unix())
	start := ""
	size := 2

	_, bt, err := GetWorkLogHash([]string{"log1", "xlog2"})
	if err != nil {
		t.Errorf("GetWorkLogHash Expected no error, got %v", err)
	}

	name := fmt.Sprint(id)
	model.AddToList(LogBucket, name, bt)

	logs, lastKey, err := ListLogsById(id, start, size, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("GetWorkLogHash Expected 1, got %v", len(logs))
	}

	// Test case 2: Error case
	size = 1
	logs, _, err = ListLogsById(id, lastKey, size, false)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if logs != nil {
		t.Errorf("Expected nil log list, got %+v", logs)
	}
}

func TestGetWorkLogHash(t *testing.T) {
	store.DBInit()
	defer store.DBClose()

	logs := []string{"log1", "log2"}

	hash2, _, err := GetWorkLogHash(logs)
	if err != nil {
		t.Errorf("GetWorkLogHash Expected no error, got %v", err)
	}

	pf := WorkLogProof{
		Time: uint64(time.Now().Unix()),
		Logs: logs,
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
