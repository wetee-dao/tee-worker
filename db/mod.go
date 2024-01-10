package db

import (
	"github.com/nutsdb/nutsdb"
)

var DB *nutsdb.DB

func DBInit() error {
	var err error
	DB, err = nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir("nutsdb"),
	)

	return err
}

func DBClose() {
	DB.Close()
}
