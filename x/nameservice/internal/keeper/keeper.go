//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"strings"

	"github.com/Masterminds/semver"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/wirelineio/wns/x/bond"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// prefixCIDToRecordIndex is the prefix for CID -> Record index.
// Note: This is the primary index in the system.
// Note: Golang doesn't support const arrays.
var prefixCIDToRecordIndex = []byte{0x00}

// prefixWRNToNameRecordIndex is the prefix for the WRN -> NamingRecord index.
var prefixWRNToNameRecordIndex = []byte{0x01}

// prefixBaseWRNToNameRecordIndex is the prefix for the Base WRN -> NamingRecord index.
// Note: BaseWRL => WRN minus `version`, i.e. latest version.
var prefixBaseWRNToNameRecordIndex = []byte{0x02}

// prefixBondIDToRecordsIndex is the prefix for the Bond ID -> [Record] index.
var prefixBondIDToRecordsIndex = []byte{0x03}

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	RecordKeeper RecordKeeper
	BondKeeper   bond.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	paramstore params.Subspace
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(recordKeeper RecordKeeper, bondKeeper bond.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec, paramstore params.Subspace) Keeper {
	return Keeper{
		RecordKeeper: recordKeeper,
		BondKeeper:   bondKeeper,
		storeKey:     storeKey,
		cdc:          cdc,
		paramstore:   paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// RecordKeeper exposes the bare minimal read-only API for other modules.
type RecordKeeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewRecordKeeper creates new instances of the nameservice RecordKeeper
func NewRecordKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) RecordKeeper {
	return RecordKeeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// PutRecord - saves a record to the store and updates ID -> Record index.
func (k Keeper) PutRecord(ctx sdk.Context, record types.Record) {
	store := ctx.KVStore(k.storeKey)
	store.Set(append(prefixCIDToRecordIndex, []byte(record.ID)...), k.cdc.MustMarshalBinaryBare(record.ToRecordObj()))
}

// Generates Bond ID -> Bond index key.
func getRecordIndexKey(id types.ID) []byte {
	return append(prefixCIDToRecordIndex, []byte(id)...)
}

// Generates Bond ID -> Records index key.
func getBondIDToRecordsIndexKey(bondID bond.ID, id types.ID) []byte {
	return append(append(prefixBondIDToRecordsIndex, []byte(bondID)...), []byte(id)...)
}

// AddBondToRecordIndexEntry adds the Bond ID -> [Record] index entry.
func (k Keeper) AddBondToRecordIndexEntry(ctx sdk.Context, bondID bond.ID, id types.ID) {
	store := ctx.KVStore(k.storeKey)
	store.Set(getBondIDToRecordsIndexKey(bondID, id), []byte{})
}

// RemoveBondToRecordIndexEntry removes the Bond ID -> [Record] index entry.
func (k Keeper) RemoveBondToRecordIndexEntry(ctx sdk.Context, bondID bond.ID, id types.ID) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(getBondIDToRecordsIndexKey(bondID, id))
}

// SetNameRecord - sets a name record.
func (k Keeper) SetNameRecord(ctx sdk.Context, wrn string, nameRecord types.NameRecord) {
	store := ctx.KVStore(k.storeKey)
	store.Set(append(prefixWRNToNameRecordIndex, []byte(wrn)...), k.cdc.MustMarshalBinaryBare(nameRecord))
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(ctx sdk.Context, id types.ID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(getRecordIndexKey(id))
}

// HasNameRecord - checks if a name record exists.
func (k Keeper) HasNameRecord(ctx sdk.Context, wrn string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(append(prefixWRNToNameRecordIndex, []byte(wrn)...))
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(ctx sdk.Context, id types.ID) types.Record {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(getRecordIndexKey(id))
	var obj types.RecordObj
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return obj.ToRecord()
}

// GetNameRecord - gets a name record from the store.
func (k Keeper) GetNameRecord(ctx sdk.Context, wrn string) types.NameRecord {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(append(prefixWRNToNameRecordIndex, []byte(wrn)...))
	var obj types.NameRecord
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return obj
}

// ListRecords - get all records.
func (k Keeper) ListRecords(ctx sdk.Context) []types.Record {
	var records []types.Record

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, prefixCIDToRecordIndex)
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

// ListNameRecords - get all name records.
func (k Keeper) ListNameRecords(ctx sdk.Context) map[string]types.NameRecord {
	nameRecords := make(map[string]types.NameRecord)

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, prefixWRNToNameRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var record types.NameRecord
			k.cdc.MustUnmarshalBinaryBare(bz, &record)
			nameRecords[string(itr.Key()[len(prefixWRNToNameRecordIndex):])] = record
		}
	}

	return nameRecords
}

// ResolveWRN resolves a WRN to a record.
// Note: Version part of the WRN might have a semver range.
func (k Keeper) ResolveWRN(ctx sdk.Context, wrn string) *types.Record {
	segments := strings.Split(wrn, "#")
	if len(segments) == 2 {
		baseWRN, semver := segments[0], segments[1]
		if strings.ContainsAny(semver, "^~<>=!") {
			// Handle semver range.
			return k.ResolveBaseWRN(ctx, baseWRN, semver)
		}
	}

	return k.ResolveFullWRN(ctx, wrn)
}

// ResolveFullWRN resolves a WRN (full path) to a record.
// Note: Version part of the WRN MUST NOT have a semver range.
func (k Keeper) ResolveFullWRN(ctx sdk.Context, wrn string) *types.Record {
	store := ctx.KVStore(k.storeKey)
	nameKey := append(prefixWRNToNameRecordIndex, []byte(wrn)...)

	if store.Has(nameKey) {
		bz := store.Get(nameKey)
		var obj types.NameRecord
		k.cdc.MustUnmarshalBinaryBare(bz, &obj)

		record := k.GetRecord(ctx, obj.ID)
		return &record
	}

	return nil
}

// ResolveBaseWRN resolves a BaseWRN + semver range to a record (picks the highest matching version).
func (k Keeper) ResolveBaseWRN(ctx sdk.Context, baseWRN string, semverRange string) *types.Record {
	semverConstraint, err := semver.NewConstraint(semverRange)
	if err != nil {
		// Handle constraint not being parsable.
		return nil
	}

	store := ctx.KVStore(k.storeKey)

	baseNameKey := append(prefixWRNToNameRecordIndex, []byte(baseWRN)...)
	if !store.Has(baseNameKey) {
		return nil
	}

	var highestSemver, _ = semver.NewVersion("0.0.0")
	var highestNameRecord types.NameRecord = types.NameRecord{}

	baseWRNPrefix := append(prefixWRNToNameRecordIndex, []byte(baseWRN)...)
	itr := sdk.KVStorePrefixIterator(store, baseWRNPrefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var record types.NameRecord
			k.cdc.MustUnmarshalBinaryBare(bz, &record)

			semver, err := semver.NewVersion(record.Version)
			if err == nil && semverConstraint.Check(semver) && semver.GreaterThan(highestSemver) {
				highestSemver = semver
				highestNameRecord = record
			}
		}
	}

	if highestNameRecord.ID != "" {
		record := k.GetRecord(ctx, highestNameRecord.ID)
		return &record
	}

	return nil
}

// MatchRecords - get all matching records.
func (k Keeper) MatchRecords(ctx sdk.Context, matchFn func(*types.Record) bool) []*types.Record {
	var records []*types.Record

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, prefixCIDToRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.RecordObj
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			record := obj.ToRecord()
			if matchFn(&record) {
				records = append(records, &record)
			}
		}
	}

	return records
}

// ClearRecords - Deletes all records and indexes.
// NOTE: FOR LOCAL TESTING PURPOSES ONLY!
func (k Keeper) ClearRecords(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	// Note: Clear everything, records and indexes.
	itr := store.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		store.Delete(itr.Key())
	}
}

// QueryRecordsByBond - get all records for the given bond.
func (k RecordKeeper) QueryRecordsByBond(ctx sdk.Context, bondID bond.ID) []types.Record {
	var records []types.Record

	bondIDPrefix := append(prefixBondIDToRecordsIndex, []byte(bondID)...)
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, bondIDPrefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		cid := itr.Key()[len(bondIDPrefix):]
		bz := store.Get(append(prefixCIDToRecordIndex, cid...))
		if bz != nil {
			var obj types.RecordObj
			k.cdc.MustUnmarshalBinaryBare(bz, &obj)
			records = append(records, obj.ToRecord())
		}
	}

	return records
}

// BondHasAssociatedRecords returns true if the bond has associated records.
func (k RecordKeeper) BondHasAssociatedRecords(ctx sdk.Context, bondID bond.ID) bool {
	bondIDPrefix := append(prefixBondIDToRecordsIndex, []byte(bondID)...)
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, bondIDPrefix)
	defer itr.Close()
	return itr.Valid()
}
