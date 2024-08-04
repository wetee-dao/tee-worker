package mint

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"

	"wetee.app/worker/internal/peer"
	types "wetee.app/worker/type"
	"wetee.app/worker/util"
)

func (c *Minter) StartP2P() error {
	if c.P2Peer != nil {
		return nil
	}
	if c.ChainClient == nil {
		return errors.New("ChainClient is nil")
	}
	if c.PrivateKey == nil {
		return errors.New("PrivateKey is nil")
	}

	// Get boot peers from chain
	bootPeers, err := c.GetBootPeers()
	if err != nil {
		fmt.Println("Get node list error:", err)
		return err
	}
	if len(bootPeers) == 0 {
		fmt.Println("No boot peers found")
		return errors.New("No boot peers found")
	}

	// 获取节点列表
	nodesFromChain, err := c.GetNodeList()
	if err != nil {
		fmt.Println("Get node list error:", err)
		return err
	}
	workersFromChain, err := c.GetWorkerList()
	if err != nil {
		fmt.Println("Get worker list error:", err)
		return err
	}

	// 获取机密节点和矿工节点
	nodes := []*types.Node{}
	for _, n := range nodesFromChain {
		var gopub ed25519.PublicKey = n.Pubkey[:]
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// 启动 P2P 网络
	port := util.GetEnvInt("P2P_PORT", 8881)
	peer, err := peer.NewP2PNetwork(ctx, c.PrivateKey, boots, nodes, uint32(port), uint32(port))
	if err != nil {
		fmt.Println("Start P2P peer error:", err)
		return err
	}

	c.P2Peer = peer

	return nil
}
