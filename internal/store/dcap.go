package store

import (
	"fmt"

	"github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/util"
)

const DcapBucket = "dcap"

// SetWorkDcapReport 将指定的工作的 DCAP 报告设置为特定的值
func SetWorkDcapReport(WorkID types.WorkId, val []byte) error {
	// 获取基于 WorkID 的键
	key := util.GetWorkTypeStr(WorkID) + "-" + fmt.Sprint(WorkID.Id) + "_dcap_report"

	// 将值保存到指定的键
	return SealSave(DcapBucket, []byte(key), val)
}

// GetWorkDcapReport 根据工作 ID 获取对应的 DCAP 报告
func GetWorkDcapReport(WorkID types.WorkId) ([]byte, error) {
	// 生成用于获取 DCAP 报告的键
	key := util.GetWorkTypeStr(WorkID) + "-" + fmt.Sprint(WorkID.Id) + "_dcap_report"

	// 使用生成的键在 DcapBucket 中获取相应的报告数据
	val, err := SealGet(DcapBucket, []byte(key))
	if err != nil {
		// 如果在获取过程中发生错误，返回错误信息
		return nil, err
	}

	// 返回获取到的报告数据和 nil 错误信息
	return val, nil
}
