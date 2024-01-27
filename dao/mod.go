package dao

import (
	"flag"

	"github.com/edgelesssys/ego/ecrypto"
	"github.com/nutsdb/nutsdb"
)

var DB *nutsdb.DB

func DBInit(path string) error {
	var err error
	DB, err = nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(path),
	)

	return err
}

func DBClose() {
	DB.Close()
}

func SealSave(bucket string, key []byte, val []byte) error {
	val, err := SealWithProductKey(val, nil)
	if err != nil {
		return err
	}

	err = checkBucket(bucket, nutsdb.DataStructureBTree)
	if err != nil {
		return err
	}

	return DB.Update(
		func(tx *nutsdb.Tx) error {
			err := tx.Put(bucket, key, val, 0)
			return err
		},
	)
}

func SealGet(bucket string, key []byte) ([]byte, error) {
	var data []byte = []byte{}
	err := checkBucket(bucket, nutsdb.DataStructureBTree)
	if err != nil {
		return nil, err
	}
	err = DB.View(
		func(tx *nutsdb.Tx) error {
			val, err := tx.Get(bucket, key)
			if err != nil {
				return err
			}

			if flag.Lookup("test.v") != nil {
				data = val
			} else {
				val, err = ecrypto.Unseal(val, nil)
				if err != nil {
					return err
				}
			}

			data = val
			return nil
		},
	)
	return data, err
}

func checkBucket(bucket string, ds uint16) error {
	return DB.Update(
		func(tx *nutsdb.Tx) error {
			if !tx.ExistBucket(ds, bucket) {
				err := tx.NewBucket(ds, bucket)
				return err
			}
			return nil
		},
	)
}

func SealWithProductKey(val []byte, additionalData []byte) ([]byte, error) {
	if flag.Lookup("test.v") == nil {
		return ecrypto.SealWithProductKey(val, additionalData)
	}
	return val, nil
}

func Unseal(ciphertext []byte, additionalData []byte) ([]byte, error) {
	if flag.Lookup("test.v") == nil {
		return ecrypto.Unseal(ciphertext, additionalData)
	}
	return ciphertext, nil
}
