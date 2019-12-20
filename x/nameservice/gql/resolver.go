//
// Copyright 2019 Wireline, Inc.
//

package gql

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wirelineio/wns/x/nameservice"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// WnsTypeProtocol => Protocol.
const WnsTypeProtocol = "wrn:protocol"

// WnsTypeBot => Bot.
const WnsTypeBot = "wrn:bot"

// WnsTypePad => Pad.
const WnsTypePad = "wrn:pad"

// WrnTypeReference => Reference.
const WrnTypeReference = "wrn:reference"

// BigUInt represents a 64-bit unsigned integer.
type BigUInt uint64

// Resolver is the GQL query resolver.
type Resolver struct {
	baseApp       *bam.BaseApp
	codec         *codec.Codec
	keeper        nameservice.Keeper
	accountKeeper auth.AccountKeeper
}

// Account resolver.
func (r *Resolver) Account() AccountResolver {
	return &accountResolver{r}
}

// Coin resolver.
func (r *Resolver) Coin() CoinResolver {
	return &coinResolver{r}
}

// Mutation is the entry point to tx execution.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query is the entry point to query execution.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type accountResolver struct{ *Resolver }

func (r *accountResolver) Number(ctx context.Context, obj *Account) (string, error) {
	val := uint64(obj.Number)
	return strconv.FormatUint(val, 10), nil
}

func (r *accountResolver) Sequence(ctx context.Context, obj *Account) (string, error) {
	val := uint64(obj.Sequence)
	return strconv.FormatUint(val, 10), nil
}

type coinResolver struct{ *Resolver }

func (r *coinResolver) Quantity(ctx context.Context, obj *Coin) (string, error) {
	val := uint64(obj.Quantity)
	return strconv.FormatUint(val, 10), nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) InsertRecord(ctx context.Context, attributes []*KeyValueInput) (*Record, error) {
	// Only supported by mock server.
	return nil, errors.New("not implemented")
}

func (r *mutationResolver) Submit(ctx context.Context, tx string) (*string, error) {
	stdTx, err := decodeStdTx(r.codec, tx)
	if err != nil {
		return nil, err
	}

	res, err := broadcastTx(r, stdTx)
	if err != nil {
		return nil, err
	}

	txHash := res.Hash.String()

	return &txHash, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetRecordsByIds(ctx context.Context, ids []string) ([]*Record, error) {
	records := make([]*Record, len(ids))
	for index, id := range ids {
		record, err := r.GetRecord(ctx, id)
		if err != nil {
			return nil, err
		}

		records[index] = record
	}

	return records, nil
}

// QueryRecords filters records by K=V conditions.
func (r *queryResolver) QueryRecords(ctx context.Context, attributes []*KeyValueInput) ([]*Record, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})
	gqlResponse := []*Record{}

	var records = r.keeper.MatchRecords(sdkContext, func(record *types.Record) bool {
		return matchOnAttributes(record, attributes)
	})

	if requestedLatestVersionsOnly(attributes) {
		records = getLatestVersions(records)
	}

	for _, record := range records {
		gqlRecord, err := getGQLRecord(ctx, r, record)
		if err != nil {
			return nil, err
		}

		gqlResponse = append(gqlResponse, gqlRecord)
	}

	return gqlResponse, nil
}

// ResolveRecords resolves records by ref/WRN, with semver range support.
func (r *queryResolver) ResolveRecords(ctx context.Context, refs []string) ([]*Record, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})
	gqlResponse := []*Record{}

	for _, ref := range refs {
		record := r.keeper.ResolveWRN(sdkContext, ref)
		gqlRecord, err := getGQLRecord(ctx, r, record)
		if err != nil {
			return nil, err
		}

		gqlResponse = append(gqlResponse, gqlRecord)
	}

	return gqlResponse, nil
}

func (r *queryResolver) GetStatus(ctx context.Context) (*Status, error) {
	return &Status{Version: NamserviceVersion}, nil
}

func (r *queryResolver) GetAccounts(ctx context.Context, addresses []string) ([]*Account, error) {
	accounts := make([]*Account, len(addresses))
	for index, address := range addresses {
		account, err := r.GetAccount(ctx, address)
		if err != nil {
			return nil, err
		}

		accounts[index] = account
	}

	return accounts, nil
}

func (r *queryResolver) GetAccount(ctx context.Context, address string) (*Account, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	account := r.accountKeeper.GetAccount(sdkContext, addr)
	if account == nil {
		return nil, nil
	}

	var pubKey *string
	if account.GetPubKey() != nil {
		pubKeyStr := base64.StdEncoding.EncodeToString(account.GetPubKey().Bytes())
		pubKey = &pubKeyStr
	}

	coins := []sdk.Coin(account.GetCoins())
	gqlCoins := make([]Coin, len(coins))

	for index, coin := range account.GetCoins() {
		amount := coin.Amount.Int64()
		if amount < 0 {
			return nil, errors.New("amount cannot be negative")
		}

		gqlCoins[index] = Coin{
			Type:     coin.Denom,
			Quantity: BigUInt(amount),
		}
	}

	accNum := BigUInt(account.GetAccountNumber())
	seq := BigUInt(account.GetSequence())

	return &Account{
		Address:  address,
		Number:   accNum,
		Sequence: seq,
		PubKey:   pubKey,
		Balance:  gqlCoins,
	}, nil
}

func (r *queryResolver) GetRecord(ctx context.Context, id string) (*Record, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	dbID := types.ID(id)
	if r.keeper.HasRecord(sdkContext, dbID) {
		record := r.keeper.GetRecord(sdkContext, dbID)
		if !record.Deleted {
			return getGQLRecord(ctx, r, &record)
		}
	}

	return nil, nil
}
