package db

import (
	"fmt"
	"strconv"

	"github.com/nutsdb/nutsdb"
)

const UserBucket = "user"

func SetClusterId(id uint64) error {
	key := []byte("rootUser")
	val := []byte(fmt.Sprint(id))
	return DB.Update(
		func(tx *nutsdb.Tx) error {
			if !tx.ExistBucket(nutsdb.DataStructureBTree, UserBucket) {
				if err := tx.NewBucket(nutsdb.DataStructureBTree, UserBucket); err != nil {
					return err
				}
				tx.SubmitBucket()
			}
			if err := tx.Put(UserBucket, key, val, 0); err != nil {
				return err
			}
			return nil
		},
	)
}

func GetClusterId() (uint64, error) {
	var id string
	err := DB.View(
		func(tx *nutsdb.Tx) error {
			val, err := tx.Get(UserBucket, []byte("rootUser"))
			if err != nil {
				return err
			}
			id = string(val)
			return nil
		},
	)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(id, 10, 64)
}
