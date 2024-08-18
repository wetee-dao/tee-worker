package mint

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"wetee.app/worker/mint/proof"
	types "wetee.app/worker/type"
)

// UploadClusterProof is a function used to upload cluster verification information
func (m *Minter) UploadClusterProof() ([]byte, error) {
	signer, _ := m.PrivateKey.ToSigner()

	// 获取 TEE 根证书
	report, t, err := proof.GetRemoteReport(signer, nil)
	if err != nil {
		fmt.Println("GetRootDcapReport => ", err)
		return nil, err
	}

	// 上传 TEE 证书
	param := types.TeeParam{
		Report:  report,
		Time:    t,
		TeeType: 0,
		Address: signer.SS58Address(42),
		Data:    nil,
	}

	// Use the json package to serialize the param object into JSON format
	bt, _ := json.Marshal(param)

	// Generate a random UUID string as the message ID
	msgId := uuid.NewV4().String()

	// Call the SendMessageToSecret method to send a message
	err = m.SendMessageToSecret(context.Background(), &types.Message{
		MsgID:   msgId,
		Type:    "upload_cluster_proof",
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

	// Return the data in the result
	return data.Result, nil
}

// UploadClusterProofreply handles the reply to the upload cluster proof request
func (m *Minter) UploadClusterProofreply(data []byte, err string, msgID string, OrgId string) error {
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
