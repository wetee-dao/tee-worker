package store

import (
	"encoding/json"
	"strings"
)

const RuningBucket = "runing"

type RuningCache struct {
	NameSpace string
	Status    string
	DeleteAt  int64
}

func SetRuning(data map[string]RuningCache) error {
	key := []byte("runing_cache")
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return SealSave(SecretBucket, key, val)
}

func GetRuning() (map[string]RuningCache, error) {
	key := []byte("runing_cache")
	val, err := SealGet(SecretBucket, key)
	if err != nil {
		if strings.Contains(err.Error(), "key not found") {
			return map[string]RuningCache{}, nil
		}
		return nil, err
	}
	s := map[string]RuningCache{}
	err = json.Unmarshal(val, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
