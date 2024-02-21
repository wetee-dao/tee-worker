package store

import (
	"fmt"

	"github.com/wetee-dao/go-sdk/gen/types"
	"wetee.app/worker/util"
)

const DcapBucket = "dcap"

func SetWorkDcapReport(WorkID types.WorkId, val []byte) error {
	key := util.GetWorkTypeStr(WorkID) + "-" + fmt.Sprint(WorkID.Id) + "_dcap_report"

	return SealSave(DcapBucket, []byte(key), val)
}

func GetWorkDcapReport(WorkID types.WorkId) (string, error) {
	key := util.GetWorkTypeStr(WorkID) + "-" + fmt.Sprint(WorkID.Id) + "_dcap_report"

	val, err := SealGet(DcapBucket, []byte(key))
	if err != nil {
		return "", err
	}
	return string(val), err
}

func SetRootDcapReport(val []byte) error {
	key := []byte("rootDcapReport")

	return SealSave(DcapBucket, key, val)
}

func GetRootDcapReport() ([]byte, error) {
	val, err := SealGet(DcapBucket, []byte("rootDcapReport"))
	if err != nil {
		return nil, err
	}

	return val, err
}
