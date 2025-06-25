package store

import (
	"fmt"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const UserBucket = "user"

func SetClusterId(id uint64) error {
	key := "clusterId"
	val := []byte(fmt.Sprint(id))
	return model.SetKey(UserBucket, key, val)
}

func GetClusterId() (uint64, error) {
	val, err := model.GetKey(UserBucket, "clusterId")
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(val), 10, 64)
}

func SetMintId(id []byte) error {
	key := "MinterId"
	val := id
	return model.SetKey(UserBucket, key, val)
}

func GetMintId() ([]byte, error) {
	val, err := model.GetKey(UserBucket, "MinterId")
	if err != nil {
		return nil, err
	}
	return val, nil
}
