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

// PutRecord - saves a record to the store.
func (k Keeper) PutRecord(ctx sdk.Context, record types.Record) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(record.ID), k.cdc.MustMarshalBinaryBare(record.ToRecordObj()))
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(ctx sdk.Context, id types.ID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(id))
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(ctx sdk.Context, id types.ID) types.Record {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get([]byte(id))
	var obj types.RecordObj
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return obj.ToRecord()
}

// ListRecords - get all records.
func (k Keeper) ListRecords(ctx sdk.Context) []types.Record {
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

// MatchRecords - get all matching records.
func (k Keeper) MatchRecords(ctx sdk.Context, matchFn func(*types.Record) bool) []types.Record {
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

// DeleteRecord - deletes a record from the store.
func (k Keeper) DeleteRecord(ctx sdk.Context, id types.ID) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(id))
}

// ClearRecords - Deletes all records.
// NOTE: FOR LOCAL TESTING PURPOSES ONLY!
func (k Keeper) ClearRecords(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		store.Delete(itr.Key())
	}
}
