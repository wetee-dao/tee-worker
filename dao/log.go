package dao

import "github.com/nutsdb/nutsdb"

const LogBucket = "log"

func Addlog(key []byte, val []byte) error {
	val, err := SealWithProductKey(val, nil)
	if err != nil {
		return err
	}
	err = checkBucket(LogBucket, nutsdb.DataStructureList)
	if err != nil {
		return err
	}

	return DB.Update(
		func(tx *nutsdb.Tx) error {
			err := tx.RPush(LogBucket, key, val)
			return err
		},
	)
}

const CrBucket = "cr"

func AddCr(key []byte, val []byte) error {
	val, err := SealWithProductKey(val, nil)
	if err != nil {
		return err
	}
	err = checkBucket(CrBucket, nutsdb.DataStructureList)
	if err != nil {
		return err
	}

	return DB.Update(
		func(tx *nutsdb.Tx) error {
			err := tx.RPush(CrBucket, key, val)
			return err
		},
	)
}
