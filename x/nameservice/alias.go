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
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper       = keeper.Keeper
	MsgSetRecord = types.MsgSetRecord
)
