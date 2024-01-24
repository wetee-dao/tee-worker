package dao

import "testing"

func TestSetRootUser(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()

	if err := SetRootUser("1"); err != nil {
		t.Error(err)
	}
}

func TestGetRootUser(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()

	if err := SetRootUser("1"); err != nil {
		t.Error(err)
	}

	if id, err := GetRootUser(); err != nil {
		t.Error(err)
	} else {
		t.Log(id)
		if id != "1" {
			t.Error("id is not 1")
		}
	}
}
