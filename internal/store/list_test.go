package store

import (
	"fmt"
	"testing"
	"time"
)

func TestAddToList(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()

	now := fmt.Sprint(time.Now().Unix())
	b1 := "log"
	b2 := "cr"
	key := []byte("key" + now)

	if err := AddToList(b1, key, []byte("val")); err != nil {
		t.Fatal(err)
	}

	res, err := GetList(b1, key, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Fatal("no result")
	}

	if err := AddToList(b2, key, []byte("val")); err != nil {
		t.Fatal(err)
	}

	res2, err := GetList(b1, key, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(res2) == 2 {
		t.Fatal("Expected 1")
	}

}
