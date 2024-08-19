package store

import (
	"encoding/json"
	"strings"
)

const RuningBucket = "runing"

type RuningCache struct {
	NameSpace string
	Status    string
	DeleteAt  int64
}

// SetRuning 将运行时缓存数据存储到指定的 bucket 和键中
func SetRuning(data map[string]RuningCache) error {
	// 将字符串 "runing_cache" 转换为字节切片作为键
	key := []byte("runing_cache")
	// 将 map 数据序列化为 JSON 格式的字节切片
	val, err := json.Marshal(data)
	if err != nil {
		// 如果序列化过程中出现错误，返回错误
		return err
	}

	// 将键值对保存到名为 SecretBucket 的 bucket 中
	return SealSave(SecretBucket, key, val)
}

// GetRuning 函数从 SecretBucket 中获取键为 "runing_cache" 的值，然后将其解封为 map[string]RuningCache 类型
func GetRuning() (map[string]RuningCache, error) {
	// 设置键为字节切片形式的 "runing_cache"
	key := []byte("runing_cache")
	// 调用 SealGet 函数从 SecretBucket 中获取键对应的值
	val, err := SealGet(SecretBucket, key)
	// 处理获取值时发生的错误
	if err != nil {
		// 如果错误是因为键不存在，则创建一个空 map[string]RuningCache 并返回
		if strings.Contains(err.Error(), "key not found") {
			return map[string]RuningCache{}, nil
		}
		return nil, err
	}
	// 解析获取到的值为 JSON 格式，将其内容填充进一个空的 map[string]RuningCache 中
	s := map[string]RuningCache{}
	err = json.Unmarshal(val, &s)
	// 处理解析 JSON 时发生的错误
	if err != nil {
		// 返回 nil 和解析 JSON 时发生的错误信息
		return nil, err
	}

	return s, nil
}
