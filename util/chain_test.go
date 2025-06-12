package util

import (
	"testing"

	"wetee.app/dsecret/type/pallet/types"
)

func TestGetWorkTypeStr(t *testing.T) {
	expect := "app"
	if expect != GetWorkTypeStr(types.WorkId{
		Wtype: types.WorkType{IsAPP: true, IsTASK: false},
		Id:    1,
	}) {
		t.Error("GetWorkTypeStr error")
	}
}

func TestGetWorkType(t *testing.T) {
	expect := types.WorkType{IsAPP: true}
	if expect != GetWorkType("app") {
		t.Error("GetWorkTypeStr error")
	}
}
