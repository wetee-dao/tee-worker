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

// TestListMonitoringsById tests the ListMonitoringsById function
func TestListMonitoringsById(t *testing.T) {
	store.DBInit("bin/testdb")
	defer store.DBClose()

	id := gtypes.WorkId{Id: uint64(time.Now().Unix() + 1)}
	page := 1
	size := 10

	var use map[string][]int64 = map[string][]int64{"x": {1, 2, 3}}

	_, _, bt, err := GetWorkCrHash(use, 1)
	if err != nil {
		t.Errorf("GetGetWorkCrHash Expected no error, got %v", err)
	}

	name := util.GetWorkTypeStr(id) + "-" + fmt.Sprint(id.Id)
	store.AddToList(CrBucket, []byte(CrBucket+name), bt)

	crs, err := ListMonitoringsById(id, page, size)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(crs) != 1 {
		t.Errorf("GetGetWorkCrHash Expected 1, got %v", len(crs))
	}

	// page = 2
	// size = 1
	// crs, err = ListMonitoringsById(id, page, size)
	// if err == nil {
	// 	t.Error("Expected error, got nil")
	// }
	// if crs != nil {
	// 	t.Errorf("Expected nil log list, got %+v", crs)
	// }
}

func TestGetGetWorkCrHash(t *testing.T) {
	store.DBInit("bin/testdb")
	defer store.DBClose()

	var blockNumber uint64 = 1
	var use map[string][]int64 = map[string][]int64{"x": {1, 2, 3}}

	hash2, _, _, err := GetWorkCrHash(use, blockNumber)
	if err != nil {
		t.Errorf("GetGetWorkCrHash Expected no error, got %v", err)
	}

	pf := WorkCrProof{
		BlockNumber: blockNumber,
		Time:        uint64(time.Now().Unix()),
		Cr:          use,
	}
	bt, err := json.Marshal(&pf)
	hash := blake2b.Sum256(bt)

	if err != nil {
		t.Errorf("GetGetWorkCrHash Expected no error, got %v", err)
	}

	if hex.EncodeToString(hash[:]) != hex.EncodeToString(hash2) {
		t.Errorf("GetGetWorkCrHash Expected %v, got %v", hash, hash2)
	}
}
