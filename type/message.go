package types

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

type Result struct {
	// 错误信息
	Error string `json:"error"`
	// 结果
	Result []byte `json:"result"`
}
