package util

import "github.com/wetee-dao/go-sdk/gen/types"

// work type to string
func GetWorkTypeStr(work types.WorkId) string {
	if work.Wtype.IsAPP {
		return "a"
	}

	if work.Wtype.IsTASK {
		return "t"
	}

	if work.Wtype.IsGPU {
		return "g"
	}

	return "unknown"
}

// string to work type
func GetWorkType(ty string) types.WorkType {
	if ty == "app" || ty == "APP" || ty == "a" {
		return types.WorkType{IsAPP: true}
	}
	if ty == "task" || ty == "TASK" || ty == "t" {
		return types.WorkType{IsTASK: true}
	}
	if ty == "gpu" || ty == "GPU" || ty == "g" {
		return types.WorkType{IsGPU: true}
	}
	return types.WorkType{}
}
