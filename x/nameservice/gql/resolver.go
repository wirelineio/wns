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
	panic("not implemented")
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
