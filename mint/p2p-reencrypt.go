package mint

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/mint/proof"
	types "wetee.app/worker/type"
)

// ReencryptSecretRequest 函数用于生成重新加密的请求，并处理返回结果
func (m *Minter) ReencryptSecretRequest(secretId string, rdrPk *types.PubKey) (*types.ReencryptSecret, error) {
	req := types.ReencryptSecretRequest{
		SecretId: secretId,
		RdrPk:    rdrPk,
	}

	bt, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal reencrypt secret request: %w", err)
	}

	// Generate a random UUID string as the message ID
	msgId := uuid.NewV4().String()

	// Call the SendMessageToSecret method to send a message
	err = m.SendMessageToSecret(context.Background(), &types.Message{
		MsgID:   msgId,
		Type:    "reencrypt_secret_remote_request",
		Payload: bt,
	})

	// If an error occurs while sending the message, return an error
	if err != nil {
		return nil, err
	}

	// Lock the mutex to ensure thread safety
	m.mu.Lock()
	// Initialize a channel for the message ID
	m.preRecerve[msgId] = make(chan interface{})
	// Unlock the mutex
	m.mu.Unlock()

	// Initialize a variable of type Result
	var data *types.Result
	// Select statement to wait for data on the channel
	select {
	// If there is data on the channel, assign it to the data variable
	case d := <-m.preRecerve[msgId]:
		data = d.(*types.Result)
	// If no data is received within 30 seconds, return a timeout error
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("timeout receiving from channel")
	}

	// Lock the mutex to ensure thread safety
	m.mu.Lock()
	delete(m.preRecerve, msgId)
	// Unlock the mutex
	m.mu.Unlock()

	// If there is an error in the result, return an error
	if data.Error != "" {
		return nil, errors.New(data.Error)
	}

	// Unmarshal the data into a ReencryptSecret struct
	var reencryptSecret types.ReencryptSecret
	err = json.Unmarshal(data.Result, &reencryptSecret)

	return &reencryptSecret, err
}

// ReencryptSecretReply 函数处理重新加密的秘密回复
func (m *Minter) ReencryptSecretReply(data []byte, err string, msgID string, OrgId string) error {
	// 检查消息ID是否存在
	if _, ok := m.preRecerve[msgID]; !ok {
		return nil
	}

	m.preRecerve[msgID] <- &types.Result{
		Error:  err,
		Result: data,
	}

	return nil
}

// LaunchFromDsecret
func (m *Minter) LaunchFromDsecret(wid *gtypes.WorkId, libosReport *types.TeeParam) (*types.ReencryptSecret, error) {
	signer, _ := m.PrivateKey.ToSigner()

	// 获取 TEE 根证书
	// get root dcap report
	report, t, err := proof.GetRemoteReport(signer, nil)
	if err != nil {
		fmt.Println("GetRootDcapReport => ", err)
		return nil, err
	}

	// 构造集群可信证明
	// make cluster dcap report
	clusterReport := types.TeeParam{
		Report:  report,
		Time:    t,
		TeeType: 0,
		Address: signer.SS58Address(42),
		Data:    nil,
	}

	// 构造启动请求
	// make launch request
	req := types.LaunchRequest{
		Libos:   libosReport,
		Cluster: &clusterReport,
	}

	bt, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal reencrypt secret request: %w", err)
	}

	// Generate a random UUID string as the message ID
	msgId := uuid.NewV4().String()

	// Call the SendMessageToSecret method to send a message
	err = m.SendMessageToSecret(context.Background(), &types.Message{
		MsgID:   msgId,
		Type:    "work_launch_request",
		Payload: bt,
	})

	// Lock the mutex to ensure thread safety
	m.mu.Lock()
	// Initialize a channel for the message ID
	m.preRecerve[msgId] = make(chan interface{})
	// Unlock the mutex
	m.mu.Unlock()

	// Initialize a variable of type Result
	var data *types.Result
	// Select statement to wait for data on the channel
	select {
	// If there is data on the channel, assign it to the data variable
	case d := <-m.preRecerve[msgId]:
		data = d.(*types.Result)
	// If no data is received within 30 seconds, return a timeout error
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("timeout receiving from channel")
	}

	// Lock the mutex to ensure thread safety
	m.mu.Lock()
	delete(m.preRecerve, msgId)
	// Unlock the mutex
	m.mu.Unlock()

	// If there is an error in the result, return an error
	if data.Error != "" {
		return nil, errors.New(data.Error)
	}

	// Unmarshal the data into a ReencryptSecret struct
	var reencryptSecret types.ReencryptSecret
	err = json.Unmarshal(data.Result, &reencryptSecret)

	return &reencryptSecret, err
}
