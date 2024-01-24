package dao

import (
	"log"
	"testing"

	"github.com/wetee-dao/go-sdk/gen/types"
)

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

func TestSetSecrets(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()
	if err := SetSecrets(types.WorkId{
		Wtype: types.WorkType{IsAPP: true},
		Id:    1,
	}, &Secrets{
		Files: map[string]string{
			"test": "test",
		},
		Env: map[string]string{
			"test": "test",
		},
	}); err != nil {
		log.Println(err)
		t.Fail()
	}
}

func TestGetSecrets(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()

	if _, err := GetSecrets(types.WorkId{
		Wtype: types.WorkType{IsAPP: true},
		Id:    1,
	}); err != nil {
		log.Println(err)
		t.Fail()
	}
}
