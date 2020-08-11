//
// Copyright 2020 Wireline, Inc.
//

package keeper

import (
	"encoding/json"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wirelineio/wns/x/auction/internal/types"
)

// query endpoints supported by the auction Querier
const (
	ListAuctions    = "list"
	GetAuction      = "get"
	QueryByOwner    = "query-by-owner"
	QueryParameters = "parameters"
	Balance         = "balance"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case ListAuctions:
			return listAuctions(ctx, path[1:], req, keeper)
		case GetAuction:
			return getAuction(ctx, path[1:], req, keeper)
		case QueryByOwner:
			return queryAuctionsByOwner(ctx, path[1:], req, keeper)
		case QueryParameters:
			return queryParameters(ctx, path[1:], req, keeper)
		case Balance:
			return queryBalance(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auction query endpoint")
		}
	}
}

// nolint: unparam
func listAuctions(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	auctions := keeper.ListAuctions(ctx)

	bz, err2 := json.MarshalIndent(auctions, "", "  ")
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func getAuction(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {

	id := types.ID(strings.Join(path, "/"))
	if !keeper.HasAuction(ctx, id) {
		return nil, sdk.ErrUnknownRequest("Auction not found.")
	}

	auction := keeper.GetAuction(ctx, id)

	bz, err2 := json.MarshalIndent(auction, "", "  ")
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func queryAuctionsByOwner(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	auctions := keeper.QueryAuctionsByOwner(ctx, path[0])

	bz, err2 := json.MarshalIndent(auctions, "", "  ")
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
	balances := keeper.GetAuctionModuleBalances(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, balances)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}
