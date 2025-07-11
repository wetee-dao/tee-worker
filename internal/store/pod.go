package store

import (
	"fmt"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const PodBucket = "pod"

// RuningCache 运行时缓存数据结构
type RuningCache struct {
	NameSpace string
	Status    string
	DeleteAt  int64
}

// SetRuning 将运行时缓存数据存储到指定的 bucket 和键中
func SetPod(pod model.Pod) error {
	// 将字符串 "runing_cache" 转换为字节切片作为键
	key := "runing_pod_" + fmt.Sprint(pod.PodId)

	// 将键值对保存到名为 SecretBucket 的 bucket 中
	return model.SetJson(PodBucket, key, &pod)
}

func DelPod(pod model.Pod) error {
	key := "runing_pod_" + fmt.Sprint(pod.PodId)
	return model.DeleteKey(PodBucket, key)
}

// GetRuning 函数从 SecretBucket 中获取键为 "runing_cache" 的值，然后将其解封为 map[string]RuningCache 类型
func GetPods() ([]*model.Pod, error) {
	// 设置键为字节切片形式的 "runing_cache"
	key := "runing_pod"
	// 调用 SealGet 函数从 SecretBucket 中获取键对应的值
	return model.GetJsonList[model.Pod](PodBucket, key)
}
