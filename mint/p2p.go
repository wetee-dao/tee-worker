package mint

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"wetee.app/worker/peer"
	"wetee.app/worker/types"
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

	// 检查节点代码是否和 wetee 上要求的版本一致

	// 获取节点列表
	nodesFromChain, err := c.GetNodeList()
	if err != nil {
		fmt.Println("Get node list error:", err)
		return err
	}

	// 获取阈值参数
	nodes := []*types.Node{}
	for _, n := range nodesFromChain {
		nodes = append(nodes, &types.Node{
			ID: hex.EncodeToString(n.Pubkey[:]),
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	boots := make([]string, 0, len(bootPeers))
	for _, b := range bootPeers {
		n := &types.Node{
			ID: hex.EncodeToString(b.Id[:]),
		}
		d := util.GetUrlFromIp1(b.Ip)
		url := "/ip4/" + d + "/tcp/" + fmt.Sprint(b.Port) + "/p2p/" + string(n.PeerID())
		boots = append(boots, url)
	}

	// 启动 P2P 网络
	port := util.GetEnvInt("P2P_PORT", 8881)
	peer, err := peer.NewP2PNetwork(ctx, c.PrivateKey, boots, uint32(port), uint32(port))
	if err != nil {
		fmt.Println("Start P2P peer error:", err)
		return err
	}

	c.P2Peer = peer

	return nil
}
