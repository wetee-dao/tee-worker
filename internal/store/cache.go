package store

import (
	"fmt"
	"strconv"
)

const CacheBucket = "cache"

func SetCacheId(id string, value int64) error {
	key := []byte(id)
	val := []byte(fmt.Sprint(value))
	return SealSave(CacheBucket, key, val)
}

func GetCacheId(id string) (int64, error) {
	val, err := SealGet(CacheBucket, []byte(id))
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(val), 10, 64)
}
