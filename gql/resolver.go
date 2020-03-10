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
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"github.com/wirelineio/wns/x/bond"
	"github.com/wirelineio/wns/x/nameservice"
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

	var records = r.keeper.MatchRecords(sdkContext, func(record *nameservice.Record) bool {
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
	rpcContext := &rpctypes.Context{}

	nodeInfo, syncInfo, validatorInfo, err := getStatusInfo(rpcContext)
	if err != nil {
		return nil, err
	}

	numPeers, peers, err := getNetInfo(rpcContext)
	if err != nil {
		return nil, err
	}

	diskUsage, err := getDiskUsage()
	if err != nil {
		return nil, err
	}

	return &Status{
		Version:   NamserviceVersion,
		Node:      *nodeInfo,
		Sync:      *syncInfo,
		Validator: validatorInfo,
		NumPeers:  numPeers,
		Peers:     peers,
		DiskUsage: diskUsage,
	}, nil
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

	accNum := BigUInt(account.GetAccountNumber())
	seq := BigUInt(account.GetSequence())

	return &Account{
		Address:  address,
		Number:   accNum,
		Sequence: seq,
		PubKey:   pubKey,
		Balance:  getGQLCoins(account.GetCoins()),
	}, nil
}

func (r *queryResolver) GetRecord(ctx context.Context, id string) (*Record, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	dbID := nameservice.ID(id)
	if r.keeper.HasRecord(sdkContext, dbID) {
		record := r.keeper.GetRecord(sdkContext, dbID)
		if !record.Deleted {
			return getGQLRecord(ctx, r, &record)
		}
	}

	return nil, nil
}

func (r *queryResolver) GetBondsByIds(ctx context.Context, ids []string) ([]*Bond, error) {
	bonds := make([]*Bond, len(ids))
	for index, id := range ids {
		bondObj, err := r.GetBond(ctx, id)
		if err != nil {
			return nil, err
		}

		bonds[index] = bondObj
	}

	return bonds, nil
}

func (r *queryResolver) GetBond(ctx context.Context, id string) (*Bond, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	dbID := bond.ID(id)
	if r.keeper.BondKeeper.HasBond(sdkContext, dbID) {
		bondObj := r.keeper.BondKeeper.GetBond(sdkContext, dbID)
		return getGQLBond(ctx, r, &bondObj)
	}

	return nil, nil
}

func (r *queryResolver) QueryBonds(ctx context.Context, attributes []*KeyValueInput) ([]*Bond, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})
	gqlResponse := []*Bond{}

	var bonds = r.keeper.BondKeeper.MatchBonds(sdkContext, func(bondObj *bond.Bond) bool {
		return matchBondOnAttributes(bondObj, attributes)
	})

	for _, bondObj := range bonds {
		gqlBond, err := getGQLBond(ctx, r, bondObj)
		if err != nil {
			return nil, err
		}

		gqlResponse = append(gqlResponse, gqlBond)
	}

	return gqlResponse, nil
}
