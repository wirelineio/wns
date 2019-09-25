package gql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/mitchellh/mapstructure"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/rpc/core"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
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
	stdTx, err := decodeStdTx(tx)
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
		record, err := r.GetResource(ctx, id)
		if err != nil {
			return nil, err
		}

		records[index] = record
	}

	return records, nil
}

func (r *queryResolver) QueryRecords(ctx context.Context, attributes []*KeyValueInput) ([]*Record, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	records := r.keeper.ListResources(sdkContext)
	gqlResponse := []*Record{}

	for _, record := range records {
		gqlRecord, err := getGQLRecord(record)
		if err != nil {
			return nil, err
		}

		if matchesOnAttributes(&record, attributes) {
			gqlResponse = append(gqlResponse, gqlRecord)
		}
	}

	return gqlResponse, nil
}

func matchesOnAttributes(record *types.Record, attributes []*KeyValueInput) bool {
	recAttrs := record.Attributes

	for _, attr := range attributes {
		recAttrVal, recAttrFound := recAttrs[attr.Key]
		if !recAttrFound {
			return false
		}

		if attr.Value.Int != nil {
			recAttrValInt, ok := recAttrVal.(int)
			if !ok || *attr.Value.Int != recAttrValInt {
				return false
			}
		}

		if attr.Value.Float != nil {
			recAttrValFloat, ok := recAttrVal.(float64)
			if !ok || *attr.Value.Float != recAttrValFloat {
				return false
			}
		}

		if attr.Value.String != nil {
			recAttrValString, ok := recAttrVal.(string)
			if !ok || *attr.Value.String != recAttrValString {
				return false
			}
		}

		if attr.Value.Boolean != nil {
			recAttrValBool, ok := recAttrVal.(bool)
			if !ok || *attr.Value.Boolean != recAttrValBool {
				return false
			}
		}

		if attr.Value.Reference != nil {
			obj, ok := recAttrVal.(map[string]interface{})
			if !ok || obj["type"].(string) != WrnTypeReference {
				return false
			}
			recAttrValRefID := obj["id"].(string)
			if recAttrValRefID != attr.Value.Reference.ID {
				return false
			}
		}

		// TODO(ashwin): Handle arrays.
	}

	return true
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
		case interface{}:
			if obj, ok := value.(map[string]interface{}); ok {
				if obj["type"].(string) == WrnTypeReference {
					kvPair.Value.Reference = &Reference{
						ID: obj["id"].(string),
					}
				}
			}
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

func decodeStdTx(tx string) (*auth.StdTx, error) {
	bytes, err := base64.StdEncoding.DecodeString(tx)
	if err != nil {
		return nil, err
	}

	// Note: json.Unmarshal doesn't known which Msg struct to use, so we do it "manually".
	// See https://stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(bytes, &objmap)
	if err != nil {
		return nil, err
	}

	var msgsArr []*json.RawMessage
	err = json.Unmarshal(*objmap["msg"], &msgsArr)
	if err != nil {
		return nil, err
	}

	var firstMsg map[string]*json.RawMessage
	err = json.Unmarshal(*msgsArr[0], &firstMsg)
	if err != nil {
		return nil, err
	}

	var messageType string
	err = json.Unmarshal(*firstMsg["type"], &messageType)
	if err != nil {
		return nil, err
	}

	var msgs []sdk.Msg

	switch messageType {
	case "nameservice/SetRecord":
		{
			var setMsg nameservice.MsgSetRecord
			err = json.Unmarshal(*firstMsg["value"], &setMsg)
			if err != nil {
				return nil, err
			}
			msgs = []sdk.Msg{setMsg}
		}
	}

	var fee auth.StdFee
	err = json.Unmarshal(*objmap["fee"], &fee)
	if err != nil {
		return nil, err
	}

	var sigs []*json.RawMessage
	err = json.Unmarshal(*objmap["signatures"], &sigs)
	if err != nil {
		return nil, err
	}

	var sig map[string]*json.RawMessage
	err = json.Unmarshal(*sigs[0], &sig)
	if err != nil {
		return nil, err
	}

	var pubKeyStr string
	err = json.Unmarshal(*sig["pub_key"], &pubKeyStr)
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, err
	}

	pubKey, err := cryptoAmino.PubKeyFromBytes(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	var signature []byte
	err = json.Unmarshal(*sig["signature"], &signature)
	if err != nil {
		return nil, err
	}

	var memo string
	err = json.Unmarshal(*objmap["memo"], &memo)
	if err != nil {
		return nil, err
	}

	stdTx := auth.StdTx{
		Msgs: msgs,
		Fee:  fee,
		Signatures: []auth.StdSignature{auth.StdSignature{
			PubKey:    pubKey,
			Signature: signature,
		}},
		Memo: memo,
	}

	return &stdTx, nil
}

func broadcastTx(r *mutationResolver, stdTx *auth.StdTx) (*ctypes.ResultBroadcastTxCommit, error) {
	txBytes, err := r.Resolver.codec.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		return nil, err
	}

	ctx := &rpctypes.Context{}
	res, err := core.BroadcastTxCommit(ctx, txBytes)
	if err != nil {
		return nil, err
	}

	if res.CheckTx.IsErr() {
		return nil, errors.New(res.CheckTx.String())
	}

	if res.DeliverTx.IsErr() {
		return nil, errors.New(res.DeliverTx.String())
	}

	return res, nil
}
