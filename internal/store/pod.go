package store

import (
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const PodBucket = "pod"
const PodKey = "run_"

// SetPod 将运行时缓存数据存储到指定的 bucket 和键中
func SetPod(pod model.Pod) error {
	key := PodKey + fmt.Sprint(pod.PodId)
	return model.SetJson(PodBucket, key, &pod)
}

// Delete ruing pod
func DelPod(pod model.Pod) error {
	key := PodKey + fmt.Sprint(pod.PodId)
	return model.DeleteKey(PodBucket, key)
}

// GetPods 函数从 SecretBucket 中获取键为 "run_" 的值
func GetPods() ([]*model.Pod, error) {
	return model.GetJsonList[model.Pod](PodBucket, PodKey)
}
