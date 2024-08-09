package store

import (
	"fmt"

	"github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/util"
)

const DcapBucket = "dcap"

func SetWorkDcapReport(WorkID types.WorkId, val []byte) error {
	key := util.GetWorkTypeStr(WorkID) + "-" + fmt.Sprint(WorkID.Id) + "_dcap_report"

	return SealSave(DcapBucket, []byte(key), val)
}

func GetWorkDcapReport(WorkID types.WorkId) ([]byte, error) {
	key := util.GetWorkTypeStr(WorkID) + "-" + fmt.Sprint(WorkID.Id) + "_dcap_report"

	val, err := SealGet(DcapBucket, []byte(key))
	if err != nil {
		return nil, err
	}

	return val, err
}

func SetWorkDeploy(WorkID types.WorkId, val []byte) error {
	key := util.GetWorkTypeStr(WorkID) + "-" + fmt.Sprint(WorkID.Id) + "_dcap_report"

	return SealSave(DcapBucket, []byte(key), val)
}

func GetWorkDeploy(WorkID types.WorkId) ([]byte, error) {
	key := util.GetWorkTypeStr(WorkID) + "-" + fmt.Sprint(WorkID.Id) + "_dcap_report"

	val, err := SealGet(DcapBucket, []byte(key))
	if err != nil {
		return nil, err
	}

	return val, err
}
