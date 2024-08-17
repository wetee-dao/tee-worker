package mint

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
	"wetee.app/worker/internal/peer"
	types "wetee.app/worker/type"
	"wetee.app/worker/util"
)

// StartP2P starts a P2P network for a blockchain node
func (m *Minter) StartP2P() error {
	// If P2Peer is not nil, indicate that it has been started, and no operation is required
	if m.P2Peer != nil {
		return nil
	}

	// If ChainClient is nil, return an error
	if m.ChainClient == nil {
		return errors.New("ChainClient is nil")
	}

	// If PrivateKey is nil, return an error
	if m.PrivateKey == nil {
		return errors.New("PrivateKey is nil")
	}

	// Get boot peers from the chain
	bootPeers, err := m.GetBootPeers()
	if err != nil {
		fmt.Println("Get node list error:", err)
		return err
	}

	// If no boot peers are found, return an error
	if len(bootPeers) == 0 {
		fmt.Println("No boot peers found")
		return errors.New("No boot peers found")
	}

	// Get the node list from the chain
	nodesFromChain, err := m.GetNodeList()
	if err != nil {
		fmt.Println("Get node list error:", err)
		return err
	}

	// Get the worker list from the chain
	workersFromChain, err := m.GetWorkerList()
	if err != nil {
		fmt.Println("Get worker list error:", err)
		return err
	}

	// Get secret nodes and miner nodes
	nodes := []*types.Node{}
	for _, n := range nodesFromChain {
		var gopub ed25519.PublicKey = n[:]
		pub, _ := types.PubKeyFromStdPubKey(gopub)
		nodes = append(nodes, &types.Node{
			ID:   pub.String(),
			Type: 1,
		})
	}
	for _, w := range workersFromChain {
		var gopub ed25519.PublicKey = w.Account[:]
		pub, _ := types.PubKeyFromStdPubKey(gopub)
		nodes = append(nodes, &types.Node{
			ID: pub.String(),
		})
	}

	boots := make([]string, 0, len(bootPeers))
	for _, b := range bootPeers {
		var gopub ed25519.PublicKey = b.Id[:]
		pub, _ := types.PubKeyFromStdPubKey(gopub)
		n := &types.Node{
			ID: pub.String(),
		}
		d := util.GetUrlFromIp1(b.Ip)
		url := d + "/tcp/" + fmt.Sprint(b.Port) + "/p2p/" + n.PeerID().String()
		boots = append(boots, url)
	}

	// Start the P2P network
	// 启动 P2P
	port := util.GetEnvInt("P2P_PORT", 8881)
	peer, err := peer.NewP2PNetwork(context.Background(), m.PrivateKey, boots, nodes, uint32(port), uint32(port))
	if err != nil {
		fmt.Println("Start P2P peer error:", err)
		return err
	}

	m.P2Peer = peer
	m.Nodes = nodes

	peer.AddHandler("worker", m.HandleWorker)
	return nil
}

// HandleWorker handles incoming messages and branches based on the message type
func (m *Minter) HandleWorker(msg *types.Message) error {
	switch msg.Type {
	/// -------------------- Proof -----------------------
	case "upload_cluster_proof_reply":
		err := m.UploadClusterProofreply(msg.Payload, msg.Error, msg.MsgID, msg.OrgId)
		return err
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

// UploadClusterProof is a function used to upload cluster verification information
func (m *Minter) UploadClusterProof(param *types.TeeParam) ([]byte, error) {
	// Use the json package to serialize the param object into JSON format
	bt, _ := json.Marshal(param)

	// Generate a random UUID string as the message ID
	msgId := uuid.NewV4().String()

	// Call the SendMessageToSecret method to send a message
	err := m.SendMessageToSecret(context.Background(), &types.Message{
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
	// Re-initialize the channel for the message ID
	m.preRecerve[msgId] = make(chan interface{})
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

// SendMessageToSecret Sends a message to a randomly selected node of type 1 within the context, while adding OrgId information
func (m *Minter) SendMessageToSecret(ctx context.Context, message *types.Message) error {
	//If the P2Peer object is nil, return an error indicating that it has not been started
	if m.P2Peer == nil {
		return errors.New("P2Peer is not start")
	}

	// Create a slice to store nodes of type 1
	nodes := make([]*types.Node, 0, len(m.Nodes))
	for _, n := range m.Nodes {
		if n.Type == 1 {
			nodes = append(nodes, n)
		}
	}
	// 随机选择一个节点
	// Randomly select an index
	randomIndex := rand.Intn(len(nodes))

	// 发送消息添加 OrgId
	// Set the OrgId field of the message object
	message.OrgId = m.P2Peer.ID().String()

	// Calling the Send method of the P2Peer object, the message is sent to the selected node
	return m.P2Peer.Send(ctx, nodes[randomIndex], "worker", message)
}
