package store

import (
	"bytes"
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/model/protoio"
)

const TEEBucket = "tee-report"

// SetPodTEEReport 将指定的工作的 DCAP 报告设置为特定的值
func SetPendingTEEReport(podId uint64, val *model.TeeCall) error {
	key := "pod_" + fmt.Sprint(podId)

	return model.SetProtoMessage(TEEBucket, key, val)
}

// GetPodTEEReport 根据工作 ID 获取对应的 DCAP 报告
func GetPendingTEEReport(podId uint64) (*model.TeeCall, error) {
	key := "pod_" + fmt.Sprint(podId)

	bt, err := model.GetKey(TEEBucket, key)
	if err != nil {
		return nil, err
	}

	tx := new(model.TeeCall)
	err = protoio.ReadMessage(bytes.NewBuffer(bt), tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetPodTEEReportList 获取所有待处理的 DCAP 报告
func GetPendingTEEReportList() ([]*model.TeeCall, error) {
	calls, _, err := model.GetProtoMessageList[model.TeeCall](TEEBucket, "pod_")
	return calls, err
}
