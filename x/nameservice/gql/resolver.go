package gql

import (
	"context"
	"encoding/base64"
	"errors"
	"reflect"
	"strconv"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/mitchellh/mapstructure"
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
	panic("not implemented")
}
func (r *mutationResolver) Submit(ctx context.Context, tx string) (*string, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetRecordsByIds(ctx context.Context, ids []string) ([]*Record, error) {
	records := make([]*Record, len(ids))
	for index, id := range ids {
		record, err := r.GetResource(ctx, id)
		if err != nil {
			return nil, err
		}

		records[index] = record
	}

	return records, nil
}

func (r *queryResolver) QueryRecords(ctx context.Context, attributes []*KeyValueInput) ([]*Record, error) {
	panic("not implemented")
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

func (r *queryResolver) GetResource(ctx context.Context, id string) (*Record, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	dbID := types.ID(id)
	if r.keeper.HasResource(sdkContext, dbID) {
		record := r.keeper.GetResource(sdkContext, dbID)
		return getGQLRecord(record)
	}

	return nil, nil
}

func getGQLRecord(record types.Record) (*Record, error) {

	attributes, err := getAttributes(&record)
	if err != nil {
		return nil, err
	}

	extension, err := getExtension(&record)
	if err != nil {
		return nil, err
	}

	return &Record{
		ID:         string(record.ID),
		Type:       record.Type(),
		Name:       record.Name(),
		Version:    record.Version(),
		Owners:     record.GetOwners(),
		Attributes: attributes,
		Extension:  extension,
	}, nil
}

func getAttributes(r *types.Record) (attributes []*KeyValue, err error) {
	attributes, err = mapToKeyValuePairs(r.Attributes)
	return
}

func getExtension(r *types.Record) (ext Extension, err error) {
	switch r.Type() {
	case WnsTypeProtocol:
		var protocol Protocol
		err := mapstructure.Decode(r.Extension, &protocol)
		return protocol, err
	case WnsTypeBot:
		var bot Bot
		err := mapstructure.Decode(r.Extension, &bot)
		return bot, err
	case WnsTypePad:
		var pad Pad
		err := mapstructure.Decode(r.Extension, &pad)
		return pad, err
	default:
		var unknown UnknownExtension
		err := mapstructure.Decode(r.Extension, &unknown)
		return unknown, err
	}
}

func mapToKeyValuePairs(attrs map[string]interface{}) ([]*KeyValue, error) {
	kvPairs := []*KeyValue{}

	trueVal := true
	falseVal := false

	for key, value := range attrs {

		kvPair := &KeyValue{
			Key: key,
		}

		switch val := value.(type) {
		case nil:
			kvPair.Value.Null = &trueVal
		case int:
			kvPair.Value.Int = &val
		case float64:
			kvPair.Value.Float = &val
		case string:
			kvPair.Value.String = &val
		case bool:
			kvPair.Value.Boolean = &val
		}

		if kvPair.Value.Null == nil {
			kvPair.Value.Null = &falseVal
		}

		valueType := reflect.ValueOf(value)
		if valueType.Kind() == reflect.Slice {
			// TODO(ashwin): Handle arrays.
		}

		kvPairs = append(kvPairs, kvPair)
	}

	return kvPairs, nil
}
