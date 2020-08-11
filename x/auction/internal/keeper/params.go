//
// Copyright 2020 Wireline, Inc.
//

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/wirelineio/wns/x/auction/internal/types"
)

// Default parameter namespace.
const (
	DefaultParamspace = types.ModuleName
)

// ParamKeyTable - ParamTable for auction module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// MaxAuctionAmount - get the max auction amount.
func (k Keeper) MaxAuctionAmount(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyMaxAuctionAmount, &res)
	return
}

// GetParams - Get all parameteras as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.MaxAuctionAmount(ctx),
	)
}

// SetParams - set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
