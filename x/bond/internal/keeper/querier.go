//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"encoding/json"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wirelineio/wns/x/bond/internal/types"
)

// query endpoints supported by the bond Querier
const (
	ListBonds    = "list"
	GetBond      = "get"
	QueryByOwner = "query-by-owner"
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
