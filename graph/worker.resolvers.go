package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vektah/gqlparser/v2/gqlerror"
	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/gen/balances"
	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	"wetee.app/worker/graph/model"
	"wetee.app/worker/mint"
	"wetee.app/worker/store"
	"wetee.app/worker/util"
)

// ClusterRegister is the resolver for the cluster_register field.
func (r *mutationResolver) ClusterRegister(ctx context.Context, name string, ip string, port int, level int) (string, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	ipstrs := strings.Split(ip, ".")
	if len(ipstrs) != 4 {
		return "", gqlerror.Errorf("Ip address format error")
	}
	iparr := []uint8{}
	for _, ipstr := range ipstrs {
		i, err := strconv.Atoi(ipstr)
		if err != nil {
			return "", gqlerror.Errorf("Ip address int transfer error")
		}
		iparr = append(iparr, uint8(i))
	}

	err := worker.ClusterRegister(name, iparr, uint32(port), uint8(level), false)
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterMortgage is the resolver for the cluster_mortgage field.
func (r *mutationResolver) ClusterMortgage(ctx context.Context, cpu int, mem int, disk int, deposit int64) (string, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	id, err := store.GetClusterId()
	if err != nil {
		return "", gqlerror.Errorf("Cant get cluster id:" + err.Error())
	}
	err = worker.ClusterMortgage(id, uint32(cpu), uint32(mem), uint32(disk), uint64(deposit), false)
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterUnmortgage is the resolver for the cluster_unmortgage field.
func (r *mutationResolver) ClusterUnmortgage(ctx context.Context, id int64) (string, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	clusterID, err := store.GetClusterId()
	if err != nil {
		return "", gqlerror.Errorf("Cant get cluster id:" + err.Error())
	}

	err = worker.ClusterUnmortgage(clusterID, uint64(id), false)
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterWithdrawal is the resolver for the cluster_withdrawal field.
func (r *mutationResolver) ClusterWithdrawal(ctx context.Context, id int64, ty model.WorkType, val int64) (string, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	err := worker.ClusterWithdrawal(gtypes.WorkId{
		Wtype: gtypes.WorkType{IsAPP: ty == model.WorkTypeApp, IsTASK: ty == model.WorkTypeTask},
		Id:    uint64(id),
	}, val, false)
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterStop is the resolver for the cluster_stop field.
func (r *mutationResolver) ClusterStop(ctx context.Context) (string, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	clusterID, err := store.GetClusterId()
	if err != nil {
		return "", gqlerror.Errorf("Cant get cluster id:" + err.Error())
	}

	err = worker.ClusterStop(clusterID, false)
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// StartForTest is the resolver for the start_for_test field.
func (r *mutationResolver) StartForTest(ctx context.Context) (bool, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return false, gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	// 1 unit of transfer
	bal, ok := new(big.Int).SetString("50000000000000000", 10)
	if !ok {
		panic(fmt.Errorf("failed to convert balance"))
	}

	minter, _ := types.NewMultiAddressFromAccountID(mint.Signer.PublicKey)
	minterWrap := gtypes.MultiAddress{
		IsId:       true,
		AsIdField0: minter.AsID,
	}
	c := balances.MakeTransferCall(minterWrap, types.NewUCompact(bal))
	err := client.SignAndSubmit(&signature.TestKeyringPairAlice, c, false)
	if err != nil {
		return false, gqlerror.Errorf("Chain call error:" + err.Error())
	}

	err = worker.ClusterRegister("", []uint8{127, 0, 0, 1}, uint32(80), uint8(1), false)
	if err != nil {
		return false, gqlerror.Errorf("Chain ClusterRegister error:" + err.Error())
	}

	time.Sleep(7 * time.Second)

	clusterId, err := worker.Getk8sClusterAccounts(mint.Signer.PublicKey)
	if err != nil {
		return false, gqlerror.Errorf("Getk8sClusterAccounts:" + err.Error())
	}
	fmt.Println("ClusterId => ", clusterId)
	store.SetClusterId(clusterId)

	id, err := store.GetClusterId()
	if err != nil {
		return false, gqlerror.Errorf("Cant get cluster id:" + err.Error())
	}

	err = worker.ClusterMortgage(id, uint32(1000000), uint32(1000000), uint32(1000000), uint64(100000000000), false)
	if err != nil {
		return false, gqlerror.Errorf("Chain ClusterMortgage error:" + err.Error())
	}

	return true, nil
}

// WorkerInfo is the resolver for the workerInfo field.
func (r *queryResolver) WorkerInfo(ctx context.Context) (*model.WorkerInfo, error) {
	root, err := store.GetRootUser()
	if err != nil {
		root = ""
	}
	var maddress = ""
	minter, err := mint.GetMintKey()
	if err == nil {
		maddress = minter.Address
	}

	report, err := store.GetRootDcapReport()
	if err != nil {
		report = nil
	}

	return &model.WorkerInfo{
		RootAddress: root,
		MintAddress: maddress,
		Report:      hex.EncodeToString(report),
	}, nil
}

// Worker is the resolver for the worker field.
func (r *queryResolver) Worker(ctx context.Context) ([]*model.Contract, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return nil, gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	clusterID, err := store.GetClusterId()
	if err != nil {
		return nil, gqlerror.Errorf("Cant get cluster id:" + err.Error())
	}

	contracts, err := worker.GetClusterContracts(clusterID, nil)
	if err != nil {
		return nil, gqlerror.Errorf("GetClusterContracts:" + err.Error())
	}

	list := make([]*model.Contract, 0, len(contracts))
	for _, contract := range contracts {
		list = append(list, &model.Contract{
			StartNumber: fmt.Sprint(contract.ContractState.StartNumber),
			User:        hex.EncodeToString(contract.ContractState.User[:]),
			WorkID:      util.GetWorkTypeStr(contract.ContractState.WorkId) + "-" + fmt.Sprint(contract.ContractState.WorkId.Id),
		})
	}
	return list, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
