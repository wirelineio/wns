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

// SaveBond - saves a bond to the store.
func (k Keeper) SaveBond(ctx sdk.Context, bond types.Bond) {
	store := ctx.KVStore(k.storeKey)

	// Bond ID -> Bond index.
	store.Set(append(prefixIDToBondIndex, []byte(bond.ID)...), k.cdc.MustMarshalBinaryBare(bond))

	// Owner -> [Bond] index.
	var key = append(prefixOwnerToBondsIndex, []byte(bond.Owner)...)
	key = append(key, []byte(bond.ID)...)
	store.Set(key, []byte{})
}

// HasBond - checks if a bond by the given ID exists.
func (k Keeper) HasBond(ctx sdk.Context, id types.ID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(append(prefixIDToBondIndex, []byte(id)...))
}

// GetBond - gets a record from the store.
func (k Keeper) GetBond(ctx sdk.Context, id types.ID) types.Bond {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(append(prefixIDToBondIndex, []byte(id)...))
	var obj types.Bond
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return obj
}

// ListBonds - get all bonds.
func (k Keeper) ListBonds(ctx sdk.Context) []types.Bond {
	var bonds []types.Bond

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, prefixIDToBondIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Bond
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			bonds = append(bonds, obj)
		}
	}

	return bonds
}

// ListBondsByOwner - list bonds by owner.
func (k Keeper) ListBondsByOwner(ctx sdk.Context, ownerAddress string) []types.Bond {
	var bonds []types.Bond

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, prefixIDToBondIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Bond
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)

			if obj.Owner == ownerAddress {
				bonds = append(bonds, obj)
			}
		}
	}

	return bonds
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
