//
// Copyright 2020 Wireline, Inc.
//

package gql

import (
	"context"
	"errors"
	"strconv"

	baseGql "github.com/wirelineio/wns/gql"
)

// BigUInt represents a 64-bit unsigned integer.
type BigUInt uint64

type accountResolver struct{ *Resolver }

type coinResolver struct{ *Resolver }

type mutationResolver struct{ *Resolver }

func (r *accountResolver) Number(ctx context.Context, obj *baseGql.Account) (string, error) {
	val := uint64(obj.Number)
	return strconv.FormatUint(val, 10), nil
}

func (r *accountResolver) Sequence(ctx context.Context, obj *baseGql.Account) (string, error) {
	val := uint64(obj.Sequence)
	return strconv.FormatUint(val, 10), nil
}

func (r *coinResolver) Quantity(ctx context.Context, obj *baseGql.Coin) (string, error) {
	val := uint64(obj.Quantity)
	return strconv.FormatUint(val, 10), nil
}

// Account resolver.
func (r *Resolver) Account() baseGql.AccountResolver {
	return &accountResolver{r}
}

// Coin resolver.
func (r *Resolver) Coin() baseGql.CoinResolver {
	return &coinResolver{r}
}

// Mutation is the entry point to tx execution.
func (r *Resolver) Mutation() baseGql.MutationResolver {
	return &mutationResolver{r}
}

func (r *mutationResolver) InsertRecord(ctx context.Context, attributes []*baseGql.KeyValueInput) (*baseGql.Record, error) {
	// Only supported by mock server.
	return nil, errors.New("Not supported")
}

func (r *mutationResolver) Submit(ctx context.Context, tx string) (*string, error) {
	// Only supported by a full-node.
	return nil, errors.New("Not supported")
}

func (r *queryResolver) GetAccounts(ctx context.Context, addresses []string) ([]*baseGql.Account, error) {
	// Only supported by a full-node.
	return nil, errors.New("Not supported")
}

func (r *queryResolver) GetAccount(ctx context.Context, address string) (*baseGql.Account, error) {
	// Only supported by a full-node.
	return nil, errors.New("Not supported")
}

func (r *queryResolver) GetBondsByIds(ctx context.Context, ids []string) ([]*baseGql.Bond, error) {
	// Only supported by a full-node.
	return nil, errors.New("Not supported")
}

func (r *queryResolver) GetBond(ctx context.Context, id string) (*baseGql.Bond, error) {
	// Only supported by a full-node.
	return nil, errors.New("Not supported")
}

func (r *queryResolver) QueryBonds(ctx context.Context, attributes []*baseGql.KeyValueInput) ([]*baseGql.Bond, error) {
	// Only supported by a full-node.
	return nil, errors.New("Not supported")
}
