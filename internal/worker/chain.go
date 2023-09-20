package worker

import (
	"context"
	"fmt"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/author"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	secretv1 "wetee.app/worker/api/v1"
)

func WorkerInit(mgr manager.Manager) {
	time.Sleep(time.Second * 3)
	c := mgr.GetClient()
	ctx := context.Background()
	pods := &appsv1.DeploymentList{}

	c.List(ctx, pods, client.InNamespace("worker-system"))
	err := c.Create(ctx, &secretv1.Tee{
		ObjectMeta: metav1.ObjectMeta{Namespace: "worker-system", Name: "ttt"},
	})
	fmt.Println("c.Create(ctx, &secretv1.Tee{})", err)

	for i := 0; i <= 10; i++ {
		time.Sleep(time.Second)
		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxx len(pods.Items) ", len(pods.Items))
	}
}

func connect() {
	// Query the system events and extract information from them. This example runs until exited via Ctrl-C

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
	fmt.Println(genesisHash.Hex())

	from := signature.TestKeyringPairAlice

	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")

	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompactFromUInt(1000000000000000000))

	ext := types.NewExtrinsic(c)

	era := types.ExtrinsicEra{IsMortalEra: false}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()

	key, err := types.CreateStorageKey(meta, "System", "Account", from.PublicKey)

	var sub *author.ExtrinsicStatusSubscription

	var accountInfo types.AccountInfo
	_, err = api.RPC.State.GetStorageLatest(key, &accountInfo)

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
				fmt.Println("accountInfo.Data.Free===>", accountInfo.Data.Free)
				return
			}
		case <-timeout:
			fmt.Println("timeout")
			return
		}
	}
}
