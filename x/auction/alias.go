//
// Copyright 2020 Wireline, Inc.
//

package auction

import (
	"github.com/wirelineio/wns/x/auction/internal/keeper"
	"github.com/wirelineio/wns/x/auction/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	AuctionBurnModuleAccountName = types.AuctionBurnModuleAccountName
	AuctionStatusCompleted       = types.AuctionStatusCompleted
)

var (
	DefaultParamspace = keeper.DefaultParamspace
	NewKeeper         = keeper.NewKeeper
	NewQuerier        = keeper.NewQuerier
	ModuleCdc         = types.ModuleCdc
	RegisterCodec     = types.RegisterCodec

	RegisterInvariants = keeper.RegisterInvariants

	NewMsgCreateAuction = types.NewMsgCreateAuction
)

type (
	ID                  = types.ID
	Auction             = types.Auction
	Keeper              = keeper.Keeper
	AuctionUsageKeeper  = types.AuctionUsageKeeper
	AuctionClientKeeper = keeper.AuctionClientKeeper
	Params              = types.Params
)
