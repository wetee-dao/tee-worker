package store

import "github.com/wetee-dao/tee-dsecret/pkg/model"

const WorkerBucket = "worker"

func SetRootUser(address string) error {
	key := "rootUser"
	val := []byte(address)
	return model.SetKey(WorkerBucket, key, val)
}

func GetRootUser() (string, error) {
	val, err := model.GetKey(WorkerBucket, "rootUser")
	if err != nil {
		return "", err
	}
	return string(val), err
}

func SetChainUrl(id string) error {
	key := "ChainUrl"
	val := []byte(id)
	return model.SetKey(WorkerBucket, key, val)
}

func GetChainUrl() (string, error) {
	val, err := model.GetKey(WorkerBucket, "ChainUrl")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
