//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
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

	// Track bond usage in other cosmos-sdk modules (more like a usage tracker).
	UsageKeepers []types.BondUsageKeeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	paramstore params.Subspace
}

// NewKeeper creates new instances of the bond Keeper
func NewKeeper(accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, supplyKeeper supply.Keeper,
	usageKeepers []types.BondUsageKeeper, storeKey sdk.StoreKey, cdc *codec.Codec, paramstore params.Subspace) Keeper {
	return Keeper{
		AccountKeeper: accountKeeper,
		CoinKeeper:    coinKeeper,
		SupplyKeeper:  supplyKeeper,
		UsageKeepers:  usageKeepers,
		storeKey:      storeKey,
		cdc:           cdc,
		paramstore:    paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// Generates Bond ID -> Bond index key.
func getBondIndexKey(id types.ID) []byte {
	return append(prefixIDToBondIndex, []byte(id)...)
}

// Generates Owner -> Bonds index key.
func getOwnerToBondsIndexKey(owner string, bondID types.ID) []byte {
	return append(append(prefixOwnerToBondsIndex, []byte(owner)...), []byte(bondID)...)
}

// SaveBond - saves a bond to the store.
func (k Keeper) SaveBond(ctx sdk.Context, bond types.Bond) {
	store := ctx.KVStore(k.storeKey)

	// Bond ID -> Bond index.
	store.Set(getBondIndexKey(bond.ID), k.cdc.MustMarshalBinaryBare(bond))

	// Owner -> [Bond] index.
	store.Set(getOwnerToBondsIndexKey(bond.Owner, bond.ID), []byte{})
}

// HasBond - checks if a bond by the given ID exists.
func (k Keeper) HasBond(ctx sdk.Context, id types.ID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(getBondIndexKey(id))
}

// DeleteBond - deletes the bond.
func (k Keeper) DeleteBond(ctx sdk.Context, bond types.Bond) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(getBondIndexKey(bond.ID))
	store.Delete(getOwnerToBondsIndexKey(bond.Owner, bond.ID))
}

// GetBond - gets a record from the store.
func (k Keeper) GetBond(ctx sdk.Context, id types.ID) types.Bond {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(getBondIndexKey(id))
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

// QueryBondsByOwner - query bonds by owner.
func (k Keeper) QueryBondsByOwner(ctx sdk.Context, ownerAddress string) []types.Bond {
	var bonds []types.Bond

	ownerPrefix := append(prefixOwnerToBondsIndex, []byte(ownerAddress)...)
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, ownerPrefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bondID := itr.Key()[len(ownerPrefix):]
		bz := store.Get(append(prefixIDToBondIndex, bondID...))
		if bz != nil {
			var obj types.Bond
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			bonds = append(bonds, obj)
		}
	}

	return bonds
}

// MatchBonds - get all matching bonds.
func (k Keeper) MatchBonds(ctx sdk.Context, matchFn func(*types.Bond) bool) []*types.Bond {
	var bonds []*types.Bond

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, prefixIDToBondIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Bond
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			if matchFn(&obj) {
				bonds = append(bonds, &obj)
			}
		}
	}

	return bonds
}
