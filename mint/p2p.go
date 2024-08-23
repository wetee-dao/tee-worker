package mint

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"math/rand"

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
	// 获取 node list
	nodesFromChain, err := m.GetNodeList()
	if err != nil {
		fmt.Println("Get node list error:", err)
		return err
	}

	// Get the worker list from the chain
	// 获取 worker list
	workersFromChain, err := m.GetWorkerList()
	if err != nil {
		fmt.Println("Get worker list error:", err)
		return err
	}

	// Get secret nodes and miner nodes
	// 获取 secret nodes 和 miner nodes
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
		d := util.GetUrlFromIp(b.Ip)
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

// SendMessageToSecret Sends a message to a randomly selected node of type 1 within the context, while adding OrgId information
func (m *Minter) SendMessageToSecret(ctx context.Context, message *types.Message) error {
	//If the P2Peer object is nil, return an error indicating that it has not been started
	if m.P2Peer == nil {
		return errors.New("P2Peer is not start")
	}

	// Create a slice to store nodes of type 1
	// 创建一个切片，存储节点类型为 1 的节点
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
