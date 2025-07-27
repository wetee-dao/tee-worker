package store

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// P2P 请求 Message 消息体
type Message struct {
	// 消息ID
	MsgID string `json:"msg_id"`
	// 来源ID
	OrgId   string `json:"org_id,omitempty"`
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
	// 错误信息
	Error string `json:"error"`
}

// Result 函数处理结果
type Result struct {
	// 错误信息
	Error string `json:"error"`
	// 结果
	Result []byte `json:"result"`
}

// LaunchRequest 函数处理启动请求
type LaunchRequest struct {
	// libos tee report
	Libos *model.TeeCall
	// cluster tee report
	Cluster *model.TeeCall
	// worker tee report
	WorkID string
}

// 公开的环境变量
type Envs struct {
	Envs  map[string]string
	Files map[string][]byte
}

// 环境变量包装
type EnvWrap struct {
	Id  uint64
	Pub Envs
	Sec model.ReencryptSecret
}
