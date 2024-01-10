package chain

import (
	"fmt"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"

	gtypes "wetee.app/worker/internal/mint/chain/gen/types"
	"wetee.app/worker/util"
)

// 区块链链接
// chain client
type ChainClient struct {
	Api     *gsrpc.SubstrateAPI
	Meta    *types.Metadata
	Hash    types.Hash
	Runtime *types.RuntimeVersion
}

// 初始化区块连链接
// init chain client
func ClientInit() (*ChainClient, error) {
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		return nil, err
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}

	runtime, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	return &ChainClient{api, meta, genesisHash, runtime}, nil
}

// 获取区块高度
// get block number
func (c *ChainClient) GetBlockNumber() (types.BlockNumber, error) {
	hash, err := c.Api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return 0, err
	}
	return hash.Number, nil
}

// 获取账户信息
// get account info
func (c *ChainClient) GetAccount(address *signature.KeyringPair) (*types.AccountInfo, error) {
	key, err := types.CreateStorageKey(c.Meta, "System", "Account", address.PublicKey)
	if err != nil {
		panic(err)
	}
	var accountInfo types.AccountInfo
	_, err = c.Api.RPC.State.GetStorageLatest(key, &accountInfo)
	return &accountInfo, err
}

// 签名并提交交易
// sign and submit transaction
func (c *ChainClient) SignAndSubmit(signer *signature.KeyringPair, runtimeCall gtypes.RuntimeCall) error {
	accountInfo, err := c.GetAccount(signer)
	if err != nil {
		return err
	}
	call, err := (runtimeCall).AsCall()
	if err != nil {
		return err
	}

	ext := types.NewExtrinsic(call)
	era := types.ExtrinsicEra{IsMortalEra: false}
	nonce := uint32(accountInfo.Nonce)

	o := types.SignatureOptions{
		BlockHash:          c.Hash,
		Era:                era,
		GenesisHash:        c.Hash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        c.Runtime.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: c.Runtime.TransactionVersion,
	}

	err = ext.Sign(*signer, o)
	if err != nil {
		return err
	}

	sub, err := c.Api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return err
	}

	defer sub.Unsubscribe()
	timeout := time.After(20 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			fmt.Printf("%#v\n", status)

			if status.IsInBlock {
				util.LogWithRed("SubmitAndWatchExtrinsic", " => IsInBlock")
			}
			if status.IsFinalized {
				util.LogWithRed("SubmitAndWatchExtrinsic", " => IsFinalized")
				return nil
			}
		case <-timeout:
			fmt.Println("timeout")
			return nil
		}
	}
}

// func ChainConnect() {
// 	// Create our API with a default connection to the local node
// 	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	meta, err := api.RPC.State.GetMetadataLatest()
// 	if err != nil {
// 		panic(err)
// 	}

// 	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(genesisHash.Hex())

// 	from := signature.TestKeyringPairAlice

// 	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
// 	if err != nil {
// 		panic(err)
// 	}
// 	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompactFromUInt(1000000000000000000))
// 	if err != nil {
// 		panic(err)
// 	}
// 	ext := types.NewExtrinsic(c)

// 	era := types.ExtrinsicEra{IsMortalEra: false}

// 	rv, err := api.RPC.State.GetRuntimeVersionLatest()
// 	if err != nil {
// 		panic(err)
// 	}
// 	key, err := types.CreateStorageKey(meta, "System", "Account", from.PublicKey)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var sub *author.ExtrinsicStatusSubscription

// 	var accountInfo types.AccountInfo
// 	_, err = api.RPC.State.GetStorageLatest(key, &accountInfo)
// 	if err != nil {
// 		panic(err)
// 	}
// 	nonce := uint32(accountInfo.Nonce)

// 	fmt.Println("accountInfo.Data.Free===>", accountInfo.Data.Free)
// 	o := types.SignatureOptions{
// 		// BlockHash:   blockHash,
// 		BlockHash:          genesisHash, // BlockHash needs to == GenesisHash if era is immortal. // TODO: add an error?
// 		Era:                era,
// 		GenesisHash:        genesisHash,
// 		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
// 		SpecVersion:        rv.SpecVersion,
// 		Tip:                types.NewUCompactFromUInt(0),
// 		TransactionVersion: rv.TransactionVersion,
// 	}

// 	err = ext.Sign(from, o)
// 	if err != nil {
// 		panic(err)
// 	}
// 	sub, err = api.RPC.Author.SubmitAndWatchExtrinsic(ext)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer sub.Unsubscribe()
// 	timeout := time.After(20 * time.Second)
// 	for {
// 		select {
// 		case status := <-sub.Chan():
// 			fmt.Printf("%#v\n", status)

// 			if status.IsInBlock {
// 				fmt.Println("IsInBlock")
// 				// return
// 			}
// 			if status.IsFinalized {
// 				fmt.Println("IsFinalized")
// 				var accountInfo types.AccountInfo
// 				api.RPC.State.GetStorageLatest(key, &accountInfo)
// 				fmt.Println("accountInfo.Data.Free ===> ", accountInfo.Data.Free)
// 				return
// 			}
// 		case <-timeout:
// 			fmt.Println("timeout")
// 			return
// 		}
// 	}
// }
