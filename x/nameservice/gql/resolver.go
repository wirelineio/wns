package gql

import (
	"context"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/wirelineio/wns/x/nameservice"
)

// BigUInt represents a 64-bit unsigned integer.
type BigUInt uint64

// Resolver is the GQL query resolver.
type Resolver struct {
	baseApp       *bam.BaseApp
	codec         *codec.Codec
	keeper        nameservice.Keeper
	accountKeeper auth.AccountKeeper
}

func (r *Resolver) Account() AccountResolver {
	return &accountResolver{r}
}
func (r *Resolver) Coin() CoinResolver {
	return &coinResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type accountResolver struct{ *Resolver }

func (r *accountResolver) Number(ctx context.Context, obj *Account) (string, error) {
	panic("not implemented")
}
func (r *accountResolver) Sequence(ctx context.Context, obj *Account) (string, error) {
	panic("not implemented")
}

type coinResolver struct{ *Resolver }

func (r *coinResolver) Quantity(ctx context.Context, obj *Coin) (string, error) {
	panic("not implemented")
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) InsertRecord(ctx context.Context, attributes []*KeyValueInput) (*Record, error) {
	panic("not implemented")
}
func (r *mutationResolver) Submit(ctx context.Context, tx string) (*string, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetStatus(ctx context.Context) (*Status, error) {
	panic("not implemented")
}
func (r *queryResolver) GetAccounts(ctx context.Context, addresses []string) ([]*Account, error) {
	panic("not implemented")
}
func (r *queryResolver) GetRecordsByIds(ctx context.Context, ids []string) ([]*Record, error) {
	panic("not implemented")
}
func (r *queryResolver) QueryRecords(ctx context.Context, attributes []*KeyValueInput) ([]*Record, error) {
	panic("not implemented")
}
