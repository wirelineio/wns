//
// Copyright 2019 Wireline, Inc.
//

package bond

import (
	"github.com/wirelineio/wns/x/bond/internal/keeper"
	"github.com/wirelineio/wns/x/bond/internal/types"
)

const (
	ModuleName                  = types.ModuleName
	RecordRentModuleAccountName = types.RecordRentModuleAccountName
	RouterKey                   = types.RouterKey
	StoreKey                    = types.StoreKey
)

var (
	DefaultParamspace = keeper.DefaultParamspace
	NewKeeper         = keeper.NewKeeper
	NewQuerier        = keeper.NewQuerier
	ModuleCdc         = types.ModuleCdc
	RegisterCodec     = types.RegisterCodec
)

type (
	ID            = types.ID
	Keeper        = keeper.Keeper
	MsgCreateBond = types.MsgCreateBond
)
