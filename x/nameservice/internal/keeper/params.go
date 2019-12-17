//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
	"time"
)

// Default parameter namespace.
const (
	DefaultParamspace = types.ModuleName
)

// ParamKeyTable - ParamTable for nameservice module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// RecordRent - get the record periodic rent.
func (k Keeper) RecordRent(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyRecordRent, &res)
	return
}

// RecordExpiryTime - get the record expiry duration.
func (k Keeper) RecordExpiryTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyRecordExpiryTime, &res)
	return
}

// GetParams - Get all parameteras as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.RecordRent(ctx),
		k.RecordExpiryTime(ctx),
	)
}

// SetParams - set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
