package util

import "github.com/wetee-dao/go-sdk/gen/types"

// work type to string
func GetWorkTypeStr(work types.WorkId) string {
	if work.Wtype.IsAPP {
		return "app"
	}

	if work.Wtype.IsTASK {
		return "task"
	}

	return "unknown"
}

// string to work type
func GetWorkType(ty string) types.WorkType {
	if ty == "app" {
		return types.WorkType{IsAPP: true}
	}
	if ty == "task" {
		return types.WorkType{IsTASK: true}
	}
	return types.WorkType{}
}
