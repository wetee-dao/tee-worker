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
	key := "key" + now

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

func TestDeleteList(t *testing.T) {
	DBInit("bin/testdb")
	defer DBClose()

	now := fmt.Sprint(time.Now().Unix() + 1)
	b1 := "log"
	key := "key" + now

	if err := AddToList(b1, key, []byte("val")); err != nil {
		t.Fatal(err)
	}

	if err := AddToList(b1, key, []byte("val")); err != nil {
		t.Fatal(err)
	}

	if err := AddToList(b1, key, []byte("val")); err != nil {
		t.Fatal(err)
	}

	if err := AddToList(b1, key, []byte("val")); err != nil {
		t.Fatal(err)
	}

	res, err := GetList(b1, key, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
	if len(res) != 4 {
		t.Fatal("expected 4")
	}

	if err := DeleteList(b1, key); err != nil {
		t.Fatal(err)
	}

	res2, err := GetList(b1, key, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(res2) != 0 {
		fmt.Println(res2)
		t.Fatal("Expected 0")
	}
}
