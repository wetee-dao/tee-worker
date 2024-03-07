package store

const WorkerBucket = "worker"

func SetRootUser(address string) error {
	key := []byte("rootUser")
	val := []byte(address)
	return SealSave(WorkerBucket, key, val)
}

func GetRootUser() (string, error) {
	val, err := SealGet(WorkerBucket, []byte("rootUser"))
	if err != nil {
		return "", err
	}
	return string(val), err
}

func SetChainUrl(id string) error {
	key := []byte("ChainUrl")
	val := []byte(id)
	return SealSave(WorkerBucket, key, val)
}

func GetChainUrl() (string, error) {
	val, err := SealGet(WorkerBucket, []byte("ChainUrl"))
	if err != nil {
		return "", err
	}
	return string(val), nil
}
