//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"encoding/json"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wirelineio/wns/x/bond/internal/types"
)

// query endpoints supported by the bond Querier
const (
	ListBonds       = "list"
	GetBond         = "get"
	QueryByOwner    = "query-by-owner"
	QueryParameters = "parameters"
	Balance         = "balance"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case ListBonds:
			return listBonds(ctx, path[1:], req, keeper)
		case GetBond:
			return getBond(ctx, path[1:], req, keeper)
		case QueryByOwner:
			return queryBondsByOwner(ctx, path[1:], req, keeper)
		case QueryParameters:
			return queryParameters(ctx, path[1:], req, keeper)
		case Balance:
			return queryBalance(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown bond query endpoint")
		}
	}
}

// nolint: unparam
func listBonds(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	bonds := keeper.ListBonds(ctx)

	bz, err2 := json.MarshalIndent(bonds, "", "  ")
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func getBond(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {

	id := types.ID(strings.Join(path, "/"))
	if !keeper.HasBond(ctx, id) {
		return nil, sdk.ErrUnknownRequest("Bond not found.")
	}

	bond := keeper.GetBond(ctx, id)

	bz, err2 := json.MarshalIndent(bond, "", "  ")
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func queryBondsByOwner(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	bonds := keeper.QueryBondsByOwner(ctx, path[0])

	bz, err2 := json.MarshalIndent(bonds, "", "  ")
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

func queryParameters(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	params := keeper.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}

func queryBalance(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	balances := map[string]sdk.Coins{}
	accountNames := []string{types.ModuleName}

	for _, accountName := range accountNames {
		moduleAddress := keeper.SupplyKeeper.GetModuleAddress(accountName)
		moduleAccount := keeper.AccountKeeper.GetAccount(ctx, moduleAddress)
		if moduleAccount != nil {
			balances[accountName] = moduleAccount.GetCoins()
		}
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, balances)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}
