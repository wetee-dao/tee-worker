package store

import (
	"log"
	"testing"

	"github.com/nutsdb/nutsdb"
	"github.com/wetee-dao/go-sdk/pallet/types"
)

func Test(t *testing.T) {
	if err := DBInit("bin/testdb"); err != nil {
		log.Fatal(err)
		t.Fail()
	}
	DBClose()
}

func TestCheckBucket(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()
	if err := checkBucket("b", nutsdb.DataStructureBTree); err != nil {
		t.Fail()
	}
}

func TestSealSave(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()
	if err := SealSave("b", []byte("key"), []byte("value")); err != nil {
		log.Fatal(err)
		t.Fail()
	}

	val, err := SealGet("b", []byte("key"))
	if err != nil || string(val) != "value" {
		t.Fail()
	}
}

func TestSealGet(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()
	if _, err := SealGet("b", []byte("key")); err != nil {
		log.Fatal(err)
		t.Fail()
	}
}

func TestSealAppID(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()

	_, err := SealAppID(types.WorkId{
		Wtype: types.WorkType{IsAPP: true},
		Id:    1,
	})
	if err != nil {
		log.Println(err)
		t.Fail()
		return
	}
}

func TestUnSealAppID(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()

	wid := types.WorkId{
		Wtype: types.WorkType{IsAPP: true},
		Id:    1,
	}
	key, err := SealAppID(wid)
	if err != nil {
		log.Println(err)
		t.Fail()
		return
	}

	wid2, err := UnSealAppID(key)
	if err != nil {
		log.Println(err)
		t.Fail()
		return
	}
	if wid2 != wid {
		log.Println(wid2, wid)
		t.Fail()
		return
	}
}
