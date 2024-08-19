package types

// TeeTrigger TEE 触发器
type TeeTrigger struct {
	// TEE 证明
	Tee TeeParam
	// 集群 ID
	ClusterId uint64
	// 调用 ID
	Callids []string
}
