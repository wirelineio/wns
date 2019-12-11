//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/wirelineio/wns/x/bond/internal/types"
)

// prefixIDToBondIndex is the prefix for ID -> Bond index in the KVStore.
// Note: This is the primary index in the system.
// Note: Golang doesn't support const arrays.
var prefixIDToBondIndex = []byte{0x00}

// prefixOwnerToBondsIndex is the prefix for the Owner -> [Bond] index in the KVStore.
var prefixOwnerToBondsIndex = []byte{0x01}

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	AccountKeeper auth.AccountKeeper
	CoinKeeper    bank.Keeper
	SupplyKeeper  supply.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the bond Keeper
func NewKeeper(accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, supplyKeeper supply.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		AccountKeeper: accountKeeper,
		CoinKeeper:    coinKeeper,
		SupplyKeeper:  supplyKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
	}
}

// CreateBond - saves a bond to the store.
func (k Keeper) CreateBond(ctx sdk.Context, bond types.Bond) {
	store := ctx.KVStore(k.storeKey)

	// Bond ID -> Bond index.
	store.Set(append(prefixIDToBondIndex, []byte(bond.ID)...), k.cdc.MustMarshalBinaryBare(bond))

	// Owner -> [Bond] index.
	var key = append(prefixOwnerToBondsIndex, []byte(bond.Owner)...)
	key = append(key, []byte(bond.ID)...)
	store.Set(key, []byte{})
}

// Clear - Deletes all entries and indexes.
// NOTE: FOR LOCAL TESTING PURPOSES ONLY!
func (k Keeper) Clear(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	// Note: Clear everything, entries and indexes.
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		store.Delete(itr.Key())
	}
}
