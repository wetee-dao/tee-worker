package types

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

type Node struct {
	ID   string `json:"id"`
	Type uint8  `json:"type"` // 0: worker, 1: dsecret
}

func (n *Node) PeerID() peer.ID {
	pk, err := PublicKeyFromLibp2pHex(n.ID)
	if err != nil {
		fmt.Println("Node types.PublicKeyFromHex error:", err)
		return peer.ID("")
	}
	peerID, err := peer.IDFromPublicKey(pk)
	if err != nil {
		fmt.Println("Node peer.IDFromPublicKey error:", err)
		return peer.ID("")
	}
	return peer.ID(peerID)
}
