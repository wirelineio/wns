//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// Default parameter namespace.
const (
	DefaultParamspace = types.ModuleName
)

// ParamKeyTable - ParamTable for nameservice module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// RecordAnnualRent - get the record annual rent.
func (k Keeper) RecordAnnualRent(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyRecordAnnualRent, &res)
	return
}

// GetParams - Get all parameteras as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.RecordAnnualRent(ctx),
	)
}

// SetParams - set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
