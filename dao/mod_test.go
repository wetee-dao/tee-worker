package dao

import (
	"log"
	"testing"

	"github.com/nutsdb/nutsdb"
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
