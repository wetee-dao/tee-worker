package store

import (
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const TEEBucket = "tee-report"

// SetPodTEEReport 将指定的工作的 DCAP 报告设置为特定的值
func SetPendingTEEReport(podId uint64, val *model.TeeCall) error {
	key := "pod_" + fmt.Sprint(podId)

	return model.SetJson(TEEBucket, key, val)
}

// GetPodTEEReport 根据工作 ID 获取对应的 DCAP 报告
func GetPendingTEEReport(podId uint64) (*model.TeeCall, error) {
	key := "pod_" + fmt.Sprint(podId)

	return model.GetJson[model.TeeCall](TEEBucket, key)
}

// GetPodTEEReportList 获取所有待处理的 DCAP 报告
func GetPendingTEEReportList() ([]*model.TeeCall, error) {
	return model.GetJsonList[model.TeeCall](TEEBucket, "pod_")
}
