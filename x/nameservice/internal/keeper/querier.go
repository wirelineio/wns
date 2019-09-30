//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"encoding/json"
	"strings"

	"github.com/wirelineio/wns/x/nameservice/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	ListResources = "list"
	GetResource   = "get"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case ListResources:
			return listResources(ctx, path[1:], req, keeper)
		case GetResource:
			return getResource(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

// nolint: unparam
func listResources(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	records := keeper.ListResources(ctx)

	bz, err2 := json.MarshalIndent(records, "", "  ")
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}

// nolint: unparam
func getResource(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {

	id := types.ID(strings.Join(path, "/"))
	if !keeper.HasResource(ctx, id) {
		return nil, sdk.ErrUnknownRequest("Record not found.")
	}

	record := keeper.GetResource(ctx, id)

	bz, err2 := json.MarshalIndent(record, "", "  ")
	if err2 != nil {
		panic("Could not marshal result to JSON.")
	}

	return bz, nil
}
