package store

import (
	"fmt"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const CacheBucket = "cache"

func SetCacheId(id string, value int64) error {
	key := id
	val := []byte(fmt.Sprint(value))
	return model.SetKey(CacheBucket, key, val)
}

func GetCacheId(id string) (int64, error) {
	val, err := model.GetKey(CacheBucket, id)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(val), 10, 64)
}
