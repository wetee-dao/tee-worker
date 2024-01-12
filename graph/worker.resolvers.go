package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"wetee.app/worker/internal/mint"
	"wetee.app/worker/internal/mint/chain"
)

// ClusterRegister is the resolver for the cluster_register field.
func (r *mutationResolver) ClusterRegister(ctx context.Context, input string) (string, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	err := worker.ClusterRegister([]uint8{127, 0, 0, 1})
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterMortgage is the resolver for the cluster_mortgage field.
func (r *mutationResolver) ClusterMortgage(ctx context.Context, input string) (string, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	err := worker.ClusterMortgage()
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterUnmortgage is the resolver for the cluster_unmortgage field.
func (r *mutationResolver) ClusterUnmortgage(ctx context.Context, input string) (string, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	err := worker.ClusterUnmortgage()
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterWithdrawal is the resolver for the cluster_withdrawal field.
func (r *mutationResolver) ClusterWithdrawal(ctx context.Context, input string) (string, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	err := worker.ClusterWithdrawal()
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterStop is the resolver for the cluster_stop field.
func (r *mutationResolver) ClusterStop(ctx context.Context, input string) (string, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
	}

	err := worker.ClusterStop()
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// Worker is the resolver for the worker field.
func (r *queryResolver) Worker(ctx context.Context) ([]string, error) {
	return []string{"127.0.0.1:8080"}, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
