package store

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

var DB *model.DB

// DBInit 初始化数据库
func DBInit() (*model.DB, error) {
	var err error
	DB, err = model.NewDB()
	return DB, err
}

// DBClose 关闭数据库
func DBClose() {
	DB.Close()
}
