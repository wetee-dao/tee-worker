package types

type Message struct {
	MsgID   string `json:"msg_id"`
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}
