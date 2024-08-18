package types

import (
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
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

// ReencryptSecretRequest 函数处理重新加密的秘密请求
type ReencryptSecretRequest struct {
	// 密文ID
	SecretId string `json:"secret_id,omitempty"`
	// 密文接收者公钥
	RdrPk *PubKey `json:"rdr_pk,omitempty"`
}

// ReencryptSecret 函数处理重新加密的结果
type ReencryptSecret struct {
	// 密文解码数据，需配合私钥使用
	XncCmt []byte `json:"xnc_cmt,omitempty"`
	// 密文
	EncScrt [][]byte `json:"enc_scrt,omitempty"`
}

// LaunchRequest 函数处理启动请求
type LaunchRequest struct {
	// libos tee report
	Libos *TeeParam
	// cluster tee report
	Cluster *TeeParam
	// worker tee report
	WorkID *gtypes.WorkId
}
