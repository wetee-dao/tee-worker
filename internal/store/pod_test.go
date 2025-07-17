package store

import (
	"testing"
	"time"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func TestAddPod(t *testing.T) {
	DBInit()

	SetPod(model.Pod{
		PodId: uint64(time.Now().Unix()),
	})

	pods, _ := GetPods()
	if len(pods) != 1 {
		t.Errorf("add pod error")
	}

	DelPod(*pods[0])

	pods, _ = GetPods()
	if len(pods) != 0 {
		t.Errorf("del pod error")
	}
}
