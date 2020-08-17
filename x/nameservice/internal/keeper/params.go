//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"time"

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

func (k Keeper) NameAuctionsEnabled(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyNameAuctions, &res)
	return
}

func (k Keeper) NameAuctionCommitsDuration(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyCommitsDuration, &res)
	return
}

func (k Keeper) NameAuctionRevealsDuration(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyRevealsDuration, &res)
	return
}

func (k Keeper) NameAuctionCommitFee(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyCommitFee, &res)
	return
}

func (k Keeper) NameAuctionRevealFee(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyRevealFee, &res)
	return
}

func (k Keeper) NameAuctionMinimumBid(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyMinimumBid, &res)
	return
}

// GetParams - Get all parameteras as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.RecordRent(ctx),
		k.RecordExpiryTime(ctx),

		k.NameAuctionsEnabled(ctx),
		k.NameAuctionCommitsDuration(ctx),
		k.NameAuctionRevealsDuration(ctx),
		k.NameAuctionCommitFee(ctx),
		k.NameAuctionRevealFee(ctx),
		k.NameAuctionMinimumBid(ctx),
	)
}

// SetParams - set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
