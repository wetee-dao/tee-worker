package store

import (
	"fmt"
	"strconv"
)

const UserBucket = "user"

func SetClusterId(id uint64) error {
	key := []byte("clusterId")
	val := []byte(fmt.Sprint(id))
	return SealSave(UserBucket, key, val)
}

func GetClusterId() (uint64, error) {
	val, err := SealGet(UserBucket, []byte("clusterId"))
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(val), 10, 64)
}

func SetMintId(id []byte) error {
	key := []byte("MinterId")
	val := id
	return SealSave(UserBucket, key, val)
}

func GetMintId() ([]byte, error) {
	val, err := SealGet(UserBucket, []byte("MinterId"))
	if err != nil {
		return nil, err
	}
	return val, nil
}
