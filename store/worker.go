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
