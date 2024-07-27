package types

import (
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/peer"
)

type Node struct {
	ID string `json:"id"`
}

func (n *Node) PeerID() peer.ID {
	pk, err := PublicKeyFromHex(n.ID)
	if err != nil {
		fmt.Println("Node types.PublicKeyFromHex error:", err)
		os.Exit(1)
	}
	peerID, err := peer.IDFromPublicKey(pk)
	if err != nil {
		fmt.Println("Node peer.IDFromPublicKey error:", err)
		os.Exit(1)
	}
	return peer.ID(peerID)
}
