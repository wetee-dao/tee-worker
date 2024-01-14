package chain

import (
	"fmt"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/xxhash"
	"github.com/pkg/errors"

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
		return nil, err
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

// query map data list
func (c *ChainClient) QueryMapAll(pallet string, method string) ([]types.StorageChangeSet, error) {
	key := createPrefixedKey(pallet, method)

	keys, err := c.Api.RPC.State.GetKeysLatest(key)
	if err != nil {
		return []types.StorageChangeSet{}, errors.Wrap(err, "[GetKeysLatest]")
	}

	set, err := c.Api.RPC.State.QueryStorageAtLatest(keys)
	if err != nil {
		return []types.StorageChangeSet{}, errors.Wrap(err, "[QueryStorageAtLatest]")
	}

	return set, nil
}

// query double map data list
func (c *ChainClient) QueryDoubleMapAll(pallet string, method string, keyarg interface{}) ([]types.StorageChangeSet, error) {
	arg, err := codec.Encode(keyarg)
	if err != nil {
		return []types.StorageChangeSet{}, err
	}

	// create key prefix
	key := createPrefixedKey(pallet, method)

	// get entry metadata
	// 获取储存元数据
	entryMeta, err := c.Meta.FindStorageEntryMetadata(pallet, method)
	if err != nil {
		return nil, err
	}

	// check if it's a map
	// 判断是否为map
	if !entryMeta.IsMap() {
		return nil, errors.New(pallet + "." + method + "is not map")
	}

	// get map hashers
	// 获取储存的 hasher 函数
	hashers, err := entryMeta.Hashers()
	if err != nil {
		return []types.StorageChangeSet{}, errors.Wrap(err, "[Hashers]")
	}

	// write key
	_, err = hashers[0].Write(arg)
	if err != nil {
		return nil, fmt.Errorf("unable to hash args[%d]: %s Error: %v", 0, arg, err)
	}
	// append hash to key
	key = append(key, hashers[0].Sum(nil)...)

	// query key
	keys, err := c.Api.RPC.State.GetKeysLatest(key)
	if err != nil {
		return []types.StorageChangeSet{}, errors.Wrap(err, "[GetKeysLatest]")
	}

	// get all data
	set, err := c.Api.RPC.State.QueryStorageAtLatest(keys)
	if err != nil {
		return []types.StorageChangeSet{}, errors.Wrap(err, "[QueryStorageAtLatest]")
	}

	return set, nil
}

func createPrefixedKey(pallet, method string) []byte {
	return append(xxhash.New128([]byte(pallet)).Sum(nil), xxhash.New128([]byte(method)).Sum(nil)...)
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
