package store

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const (
	dbPath = "./chain_data/wetee"
)

var DB *model.DB

func DBInit() (*model.DB, error) {
	var err error
	DB, err = model.NewDB()
	return DB, err
}

func DBClose() {
	DB.Close()
}
