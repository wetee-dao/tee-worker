package db

import (
	"github.com/edgelesssys/ego/ecrypto"
	"github.com/nutsdb/nutsdb"
)

var DB *nutsdb.DB

func DBInit() error {
	var err error
	DB, err = nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir("/opt/nutsdb"),
	)

	return err
}

func DBClose() {
	DB.Close()
}

func SealSave(bucket string, key []byte, val []byte) error {
	val, err := ecrypto.SealWithProductKey(val, nil)
	if err != nil {
		return err
	}
	return DB.Update(
		func(tx *nutsdb.Tx) error {
			if !tx.ExistBucket(nutsdb.DataStructureBTree, bucket) {
				if err := tx.NewBucket(nutsdb.DataStructureBTree, bucket); err != nil {
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

func SealGet(bucket string, key []byte) ([]byte, error) {
	var data []byte = []byte{}
	err := DB.View(
		func(tx *nutsdb.Tx) error {
			val, err := tx.Get(bucket, key)
			if err != nil {
				return err
			}

			val, err = ecrypto.Unseal(val, nil)
			if err != nil {
				return err
			}
			data = val
			return nil
		},
	)
	return data, err
}
