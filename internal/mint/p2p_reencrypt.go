package mint

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"wetee.app/worker/internal/store"
)

// ReencryptSecretRequest 函数用于生成重新加密的请求，并处理返回结果
func (m *Minter) ReencryptSecretRequest(secretId string, rdrPk *model.PubKey) (*model.ReencryptSecret, error) {
	req := model.ReencryptSecretRequest{
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
	err = m.SendMessageToSecret(context.Background(), &store.Message{
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
	m.preRecerve[msgId] = make(chan any)
	// Unlock the mutex
	m.mu.Unlock()

	// Initialize a variable of type Result
	var data *store.Result
	// Select statement to wait for data on the channel
	select {
	// If there is data on the channel, assign it to the data variable
	case d := <-m.preRecerve[msgId]:
		data = d.(*store.Result)
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
	var reencryptSecret model.ReencryptSecret
	err = json.Unmarshal(data.Result, &reencryptSecret)

	return &reencryptSecret, err
}

// ReencryptSecretReply 函数处理重新加密的秘密回复
func (m *Minter) ReencryptSecretReply(data []byte, err string, msgID string, OrgId string) error {
	// 检查消息ID是否存在
	if _, ok := m.preRecerve[msgID]; !ok {
		return nil
	}

	m.preRecerve[msgID] <- &store.Result{
		Error:  err,
		Result: data,
	}

	return nil
}

// LaunchFromDsecret 函数处理重新加密的秘密回复
func (m *Minter) LaunchFromDsecret(pid uint64, libosReport *model.TeeCall) (*model.ReencryptSecret, error) {
	signer, _ := m.PrivateKey.ToSigner()

	// 构造集群可信证明
	// make cluster dcap report
	clusterReport := model.TeeCall{
		TeeType: 0,
		Caller:  signer.PublicKey,
	}

	err := model.IssueReport(signer, &clusterReport)
	if err != nil {
		fmt.Println("GetRootDcapReport => ", err)
		return nil, err
	}

	// 构造启动请求
	// make launch request
	req := store.LaunchRequest{
		WorkID:  fmt.Sprint(pid),
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
	err = m.SendMessageToSecret(context.Background(), &store.Message{
		MsgID:   msgId,
		Type:    "work_launch_request",
		Payload: bt,
	})
	if err != nil {
		return nil, err
	}

	// Lock the mutex to ensure thread safety
	m.mu.Lock()
	// Initialize a channel for the message ID
	m.preRecerve[msgId] = make(chan any)
	// Unlock the mutex
	m.mu.Unlock()

	// Initialize a variable of type Result
	var data *store.Result
	// Select statement to wait for data on the channel
	select {
	// If there is data on the channel, assign it to the data variable
	case d := <-m.preRecerve[msgId]:
		data = d.(*store.Result)
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
	var reencryptSecret model.ReencryptSecret
	err = json.Unmarshal(data.Result, &reencryptSecret)

	return &reencryptSecret, err
}

// LaunchFromDsecret 函数处理重新加密的秘密回复
func (m *Minter) WorkLaunchReply(data []byte, err string, msgID string, OrgId string) error {
	// 检查消息ID是否存在
	if _, ok := m.preRecerve[msgID]; !ok {
		return nil
	}

	m.preRecerve[msgID] <- &store.Result{
		Error:  err,
		Result: data,
	}

	return nil
}
