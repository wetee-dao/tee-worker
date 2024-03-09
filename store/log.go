package store

import (
	"github.com/nutsdb/nutsdb"
)

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
			err := tx.LPush(LogBucket, key, val)
			return err
		},
	)
}

func GetLogList(key []byte, page int, size int) ([][]byte, error) {
	err := checkBucket(LogBucket, nutsdb.DataStructureList)
	if err != nil {
		return nil, err
	}

	list := make([][]byte, 0, size)
	err = DB.View(
		func(tx *nutsdb.Tx) error {
			var start = 0
			var end = size
			if page > 1 {
				start = (page - 1) * size
				end = start + size
			}
			clist, err2 := tx.LRange(LogBucket, key, start, end)
			for _, v := range clist {
				item, err := Unseal(v, nil)
				if err != nil {
					return err
				}
				list = append(list, item)
			}
			return err2
		},
	)
	return list, err
}

const CrBucket = "cr"

func GetMetricList(key []byte, page int, size int) ([][]byte, error) {
	err := checkBucket(CrBucket, nutsdb.DataStructureList)
	if err != nil {
		return nil, err
	}

	list := make([][]byte, 0, size)
	err = DB.View(
		func(tx *nutsdb.Tx) error {
			var start = 0
			var end = size
			if page > 1 {
				start = (page - 1) * size
				end = start + size
			}

			clist, err2 := tx.LRange(CrBucket, key, start, end)
			for _, v := range clist {
				item, err := Unseal(v, nil)
				if err != nil {
					return err
				}
				list = append(list, item)
			}
			return err2
		},
	)
	return list, err
}

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
			err := tx.LPush(CrBucket, key, val)
			return err
		},
	)
}
