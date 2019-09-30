//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	CoinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// PutResource - saves a record to the store.
func (k Keeper) PutResource(ctx sdk.Context, record types.Record) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(record.ID), k.cdc.MustMarshalBinaryBare(record.ToRecordObj()))
}

// HasResource - checks if a record by the given ID exists.
func (k Keeper) HasResource(ctx sdk.Context, id types.ID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(id))
}

// GetResource - gets a record from the store.
func (k Keeper) GetResource(ctx sdk.Context, id types.ID) types.Record {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get([]byte(id))
	var obj types.RecordObj
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return obj.ToRecord()
}

// ListResources - get all records.
func (k Keeper) ListResources(ctx sdk.Context) []types.Record {
	var records []types.Record

	store := ctx.KVStore(k.storeKey)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.RecordObj
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			records = append(records, obj.ToRecord())
		}
	}

	return records
}

// MatchResources - get all matching records.
func (k Keeper) MatchResources(ctx sdk.Context, matchFn func(*types.Record) bool) []types.Record {
	var records []types.Record

	store := ctx.KVStore(k.storeKey)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.RecordObj
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			record := obj.ToRecord()
			if matchFn(&record) {
				records = append(records, record)
			}
		}
	}

	return records
}

// DeleteResource - deletes a record from the store.
func (k Keeper) DeleteResource(ctx sdk.Context, id types.ID) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(id))
}

// ClearResources - Deletes all records.
// NOTE: FOR LOCAL TESTING PURPOSES ONLY!
func (k Keeper) ClearResources(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		store.Delete(itr.Key())
	}
}
