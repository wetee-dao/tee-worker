package store

import (
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const PodBucket = "pod"
const PodKey = "run_"

// SetPod
func SetPod(pod model.Pod) error {
	key := PodKey + fmt.Sprint(pod.PodId)
	return model.SetJson(PodBucket, key, &pod)
}

// Delete ruing pod
func DelPod(pod model.Pod) error {
	key := PodKey + fmt.Sprint(pod.PodId)
	return model.DeleteKey(PodBucket, key)
}

// GetPods
func GetPods() ([]*model.Pod, error) {
	return model.GetJsonList[model.Pod](PodBucket, PodKey)
}
