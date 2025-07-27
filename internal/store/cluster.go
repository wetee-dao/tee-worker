package store

import (
	"fmt"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const UserBucket = "user"

func SetClusterId(id uint64) error {
	key := "clusterId"
	val := fmt.Append(nil, id)
	return model.SetKey(UserBucket, key, val)
}

func GetClusterId() (uint64, error) {
	val, err := model.GetKey(UserBucket, "clusterId")
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(val), 10, 64)
}
