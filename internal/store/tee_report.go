package store

import (
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const TEEBucket = "tee-report"

// SetPodTEEReport 将指定的工作的 DCAP 报告设置为特定的值
func SetPendingTEEReport(podId uint64, val *model.TeeParam) error {
	key := "pod_" + fmt.Sprint(podId)

	return model.SetJson(TEEBucket, key, val)
}

// GetPodTEEReport 根据工作 ID 获取对应的 DCAP 报告
func GetPendingTEEReport(podId uint64) (*model.TeeParam, error) {
	key := "pod_" + fmt.Sprint(podId)

	return model.GetJson[model.TeeParam](TEEBucket, key)
}

// GetPodTEEReportList 获取所有待处理的 DCAP 报告
func GetPendingTEEReportList() ([]*model.TeeParam, error) {
	return model.GetJsonList[model.TeeParam](TEEBucket, "pod_")
}
