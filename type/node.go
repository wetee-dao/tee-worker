package types

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

type Node struct {
	// publicKey hex
	ID   string `json:"id"`
	Type uint8  `json:"type"` // 0: worker, 1: dsecret
}

// PeerID 函数返回节点的 peer.ID
func (n *Node) PeerID() peer.ID {
	// 从 n.ID 中获取公钥，作为 PublicKeyFromLibp2pHex 的输入
	pk, err := PublicKeyFromLibp2pHex(n.ID)
	// 如果获取公钥时发生错误，打印错误信息并返回空的 peer.ID
	if err != nil {
		fmt.Println("Node types.PublicKeyFromHex error:", err)
		return peer.ID("")
	}
	// 使用获取到的公钥生成 peer.ID，这里作为 IDFromPublicKey 的输入
	peerID, err := peer.IDFromPublicKey(pk)
	// 如果生成 peer.ID 时发生错误，打印错误信息并返回空的 peer.ID
	if err != nil {
		fmt.Println("Node peer.IDFromPublicKey error:", err)
		return peer.ID("")
	}
	// 返回生成的 peer.ID
	return peer.ID(peerID)
}
