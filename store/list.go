package store

import (
	"github.com/nutsdb/nutsdb"
)

func AddToList(bucket string, key []byte, val []byte) error {
	val, err := SealWithProductKey(val, nil)
	if err != nil {
		return err
	}
	err = checkBucket(bucket, nutsdb.DataStructureList)
	if err != nil {
		return err
	}

	return DB.Update(
		func(tx *nutsdb.Tx) error {
			err := tx.LPush(bucket, key, val)
			return err
		},
	)
}

func GetList(bucket string, key []byte, page int, size int) ([][]byte, error) {
	err := checkBucket(bucket, nutsdb.DataStructureList)
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
			clist, err2 := tx.LRange(bucket, key, start, end)
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
