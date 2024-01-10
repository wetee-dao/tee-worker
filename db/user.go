package db

import (
	"github.com/nutsdb/nutsdb"
)

const UserBucket = "user"

func SetRootUser(address string) error {
	key := []byte("rootUser")
	val := []byte(address)
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

func GetRootUser() (string, error) {
	var address string
	err := DB.View(
		func(tx *nutsdb.Tx) error {
			val, err := tx.Get(UserBucket, []byte("rootUser"))
			if err != nil {
				return err
			}
			address = string(val)
			return nil
		},
	)
	return address, err
}
