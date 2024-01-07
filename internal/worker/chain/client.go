package chain

import (
	"fmt"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/author"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type ChainClient struct {
	api  *gsrpc.SubstrateAPI
	meta *types.Metadata
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
	fmt.Println(genesisHash.Hex())

	return &ChainClient{api, meta}, nil
}

func (c *ChainClient) GetBlockHash() (string, error) {
	genesisHash, err := c.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return "", err
	}
	return genesisHash.Hex(), nil
}

func ChainConnect() {
	// Create our API with a default connection to the local node
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}
	fmt.Println(genesisHash.Hex())

	from := signature.TestKeyringPairAlice

	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}
	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompactFromUInt(1000000000000000000))
	if err != nil {
		panic(err)
	}
	ext := types.NewExtrinsic(c)

	era := types.ExtrinsicEra{IsMortalEra: false}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}
	key, err := types.CreateStorageKey(meta, "System", "Account", from.PublicKey)
	if err != nil {
		panic(err)
	}
	var sub *author.ExtrinsicStatusSubscription

	var accountInfo types.AccountInfo
	_, err = api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		panic(err)
	}
	nonce := uint32(accountInfo.Nonce)

	fmt.Println("accountInfo.Data.Free===>", accountInfo.Data.Free)
	o := types.SignatureOptions{
		// BlockHash:   blockHash,
		BlockHash:          genesisHash, // BlockHash needs to == GenesisHash if era is immortal. // TODO: add an error?
		Era:                era,
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	err = ext.Sign(from, o)
	if err != nil {
		panic(err)
	}
	sub, err = api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		panic(err)
	}

	defer sub.Unsubscribe()
	timeout := time.After(20 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			fmt.Printf("%#v\n", status)

			if status.IsInBlock {
				fmt.Println("IsInBlock")
				// return
			}
			if status.IsFinalized {
				fmt.Println("IsFinalized")
				var accountInfo types.AccountInfo
				api.RPC.State.GetStorageLatest(key, &accountInfo)
				fmt.Println("accountInfo.Data.Free ===> ", accountInfo.Data.Free)
				return
			}
		case <-timeout:
			fmt.Println("timeout")
			return
		}
	}
}
