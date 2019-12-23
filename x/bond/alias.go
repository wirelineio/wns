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

const (
	MicroWire = types.MicroWire
	MilliWire = types.MilliWire
)

var (
	DefaultParamspace = keeper.DefaultParamspace
	NewKeeper         = keeper.NewKeeper
	NewQuerier        = keeper.NewQuerier
	ModuleCdc         = types.ModuleCdc
	RegisterCodec     = types.RegisterCodec

	RegisterInvariants = keeper.RegisterInvariants
)

type (
	ID            = types.ID
	Bond          = types.Bond
	Keeper        = keeper.Keeper
	MsgCreateBond = types.MsgCreateBond
)
