package gql

import (
	"context"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/wirelineio/wns/x/nameservice"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

// Resolver is the GQL query resolver.
type Resolver struct {
	baseApp       *bam.BaseApp
	codec         *codec.Codec
	keeper        nameservice.Keeper
	accountKeeper auth.AccountKeeper
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Submit(ctx context.Context, tx string) (*string, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

// GetStatus returns the registry status.
func (r *queryResolver) GetStatus(ctx context.Context) (*Status, error) {
	return &Status{Version: NamserviceVersion}, nil
}

func (r *queryResolver) GetAccounts(ctx context.Context, addresses []string) ([]*Account, error) {
	panic("not implemented")
}
func (r *queryResolver) GetRecordsByIds(ctx context.Context, ids []string) ([]*Record, error) {
	panic("not implemented")
}
func (r *queryResolver) GetRecordsByAttributes(ctx context.Context, attributes []*KeyValueInput) ([]*Record, error) {
	panic("not implemented")
}
func (r *queryResolver) GetBotsByAttributes(ctx context.Context, attributes []*KeyValueInput) ([]*Bot, error) {
	panic("not implemented")
}
