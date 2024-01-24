package dao

import "testing"

func TestSetClusterId(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()

	if err := SetClusterId(1); err != nil {
		t.Error(err)
	}
}

func TestGetClusterId(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()

	if id, err := GetClusterId(); err != nil {
		t.Error(err)
	} else {
		t.Log(id)
	}
}
