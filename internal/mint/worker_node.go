package mint

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/pallet/dsecret"
	"github.com/wetee-dao/go-sdk/pallet/types"
	"wetee.app/worker/util"
)

// RegisterNode register node
// 注册节点
func (c *Minter) RegisterNode(signer *chain.Signer, pubkey []byte) error {
	var bt [32]byte
	copy(bt[:], pubkey)

	call := dsecret.MakeRegisterNodeCall(bt)
	return c.ChainClient.SignAndSubmit(signer, call, true)
}

// 获取全网当前程序的代码版本
// Get CodeSignature
func (c *Minter) GetCodeSignature() ([]byte, error) {
	return dsecret.GetCodeSignatureLatest(c.ChainClient.Api.RPC.State)
}

// 获取全网当前程序的签名人
// Get GetCodeSigner
func (c *Minter) GetGetCodeSigner() ([]byte, error) {
	return dsecret.GetCodeSignerLatest(c.ChainClient.Api.RPC.State)
}

// GetNodeList get node list
// 获取节点列表
func (c *Minter) GetNodeList() ([][32]byte, error) {
	ret, err := c.ChainClient.QueryMapAll("DSecret", "Nodes")
	if err != nil {
		return nil, err
	}

	nodes := make([][32]byte, 0)
	for _, elem := range ret {
		for _, change := range elem.Changes {
			n := [32]byte{}
			if err := codec.Decode(change.StorageData, &n); err != nil {
				util.LogError("codec.Decode", err)
				continue
			}
			nodes = append(nodes, n)
		}
	}

	return nodes, nil
}

// GetWorkerList get worker list
// 获取矿工列表
func (c *Minter) GetWorkerList() ([]*types.K8sCluster, error) {
	ret, err := c.ChainClient.QueryMapAll("Worker", "K8sClusters")
	if err != nil {
		return nil, err
	}

	// 获取节点列表
	nodes := make([]*types.K8sCluster, 0)
	for _, elem := range ret {
		for _, change := range elem.Changes {
			n := &types.K8sCluster{}
			if err := codec.Decode(change.StorageData, n); err != nil {
				util.LogError("codec.Decode", err)
				continue
			}
			nodes = append(nodes, n)
		}
	}

	return nodes, nil
}

func (c *Minter) GetBootPeers() ([]types.P2PAddr, error) {
	return dsecret.GetBootPeersLatest(c.ChainClient.Api.RPC.State)
}
