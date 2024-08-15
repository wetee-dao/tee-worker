package mint

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/wetee-dao/go-sdk/core"
	"github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/go-sdk/pallet/weteedsecret"
	"github.com/wetee-dao/go-sdk/pallet/weteeworker"
	"wetee.app/worker/util"
)

// RegisterNode register node
// 注册节点
func (c *Minter) RegisterNode(signer *core.Signer, pubkey []byte) error {
	var bt [32]byte
	copy(bt[:], pubkey)

	call := weteedsecret.MakeRegisterNodeCall(bt)
	return c.ChainClient.SignAndSubmit(signer, call, true)
}

// GetNodeList get node list
// 获取节点列表
func (c *Minter) GetNodeList() ([][32]byte, error) {
	ret, err := c.ChainClient.QueryMapAll("WeTEEDsecret", "Nodes")
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

// 获取全网当前程序的代码版本
// Get CodeMrenclave
func (c *Minter) GetCodeMrenclave() ([]byte, error) {
	return weteedsecret.GetCodeMrenclaveLatest(c.ChainClient.Api.RPC.State)
}

// 获取全网当前程序的签名人
// Get CodeMrsigner
func (c *Minter) GetCodeMrsigner() ([]byte, error) {
	return weteedsecret.GetCodeMrsignerLatest(c.ChainClient.Api.RPC.State)
}

func (c *Minter) GetWorkerList() ([]*types.K8sCluster, error) {
	ret, err := c.ChainClient.QueryMapAll("WeTEEWorker", "K8sClusters")
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
	return weteeworker.GetBootPeersLatest(c.ChainClient.Api.RPC.State)
}
