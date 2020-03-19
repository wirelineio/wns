//
// Copyright 2019 Wireline, Inc.
//

package nameservice

import (
	"github.com/wirelineio/wns/x/nameservice/internal/keeper"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	DefaultParamspace = keeper.DefaultParamspace
	NewKeeper         = keeper.NewKeeper
	NewRecordKeeper   = keeper.NewRecordKeeper
	NewQuerier        = keeper.NewQuerier
	ModuleCdc         = types.ModuleCdc
	RegisterCodec     = types.RegisterCodec

	RegisterInvariants = keeper.RegisterInvariants

	GetBlockChangesetIndexKey = keeper.GetBlockChangesetIndexKey
	GetRecordIndexKey         = keeper.GetRecordIndexKey
	GetNameRecordIndexKey     = keeper.GetNameRecordIndexKey

	HasRecord = keeper.HasRecord
	GetRecord = keeper.GetRecord
)

type (
	Keeper       = keeper.Keeper
	RecordKeeper = keeper.RecordKeeper

	MsgSetRecord = types.MsgSetRecord

	ID         = types.ID
	Record     = types.Record
	RecordObj  = types.RecordObj
	NameRecord = types.NameRecord

	BlockChangeset = types.BlockChangeset
)
