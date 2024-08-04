package peer

import (
	"github.com/libp2p/go-libp2p/core/control"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	types "wetee.app/worker/type"
)

func newConnectionGater(nodes []*types.Node) *ChainConnectionGater {
	return &ChainConnectionGater{nodes}
}

// 节点验证器
type ChainConnectionGater struct {
	Nodes []*types.Node
}

// 过滤非网络节点，减轻网络压力
func (g *ChainConnectionGater) chainRoutingTableFilter(dht interface{}, p peer.ID) bool {
	if len(g.Nodes) > 0 {
		for _, n := range g.Nodes {
			if n.PeerID() == p {
				return true
			}
		}
	}
	return false
}

// 在这里实现节点连接前的身份验证逻辑
// 例如,可以要求节点提供数字签名或者预共享密钥来验证身份
func (g *ChainConnectionGater) InterceptPeerDial(p peer.ID) bool {
	if len(g.Nodes) > 0 {
		for _, n := range g.Nodes {
			if n.PeerID() == p {
				return true
			}
		}
	}
	return false
}

// 在这里实现节点身份验证完成后的进一步检查
// 例如,可以根据节点的信任度或者其他自定义指标来决定是否接受连接
func (g *ChainConnectionGater) InterceptAddrDial(p peer.ID, info ma.Multiaddr) (allow bool) {
	return true
}

// 在这里实现节点地址验证逻辑
// 例如,可以根据节点的地址类型或者地址的IP地址段来决定是否接受连接
func (g *ChainConnectionGater) InterceptAccept(addrs network.ConnMultiaddrs) bool {
	return true
}

// 这个方法在节点成功建立一个安全的连接时被调用。
// 开发者可以在这个方法中实现对连接安全性的进一步检查,比如验证对方节点的身份、检查加密算法的强度等。
func (g *ChainConnectionGater) InterceptSecured(network.Direction, peer.ID, network.ConnMultiaddrs) (allow bool) {
	return true
}

// 这个方法在节点成功将一个原始连接升级为一个多路复用的连接时被调用。
// 开发者可以在这个方法中实现对连接升级过程的检查和验证,比如确保升级后的连接仍然满足安全性要求。
func (g *ChainConnectionGater) InterceptUpgraded(network.Conn) (allow bool, reason control.DisconnectReason) {
	return true, control.DisconnectReason(0)
}
