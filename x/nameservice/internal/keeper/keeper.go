//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"bytes"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/wirelineio/wns/x/bond"
	"github.com/wirelineio/wns/x/nameservice/internal/helpers"
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

// prefixExpiryTimeToRecordsIndex is the prefix for the Expiry Time -> [Record] index.
var prefixExpiryTimeToRecordsIndex = []byte{0x10}

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

// getRecordExpiryQueueTimeKey gets the prefix for the record expiry queue.
func getRecordExpiryQueueTimeKey(timestamp time.Time) []byte {
	timeBytes := sdk.FormatTimeBytes(timestamp)
	return append(prefixExpiryTimeToRecordsIndex, timeBytes...)
}

// GetRecordExpiryQueueTimeSlice gets a specific record queue timeslice.
// A timeslice is a slice of CIDs corresponding to records that expire at a certain time.
func (k Keeper) GetRecordExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (cids []types.ID) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(getRecordExpiryQueueTimeKey(timestamp))
	if bz == nil {
		return []types.ID{}
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &cids)
	return cids
}

// SetRecordExpiryQueueTimeSlice sets a specific record expiry queue timeslice.
func (k Keeper) SetRecordExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time, cids []types.ID) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(cids)
	store.Set(getRecordExpiryQueueTimeKey(timestamp), bz)
}

// DeleteRecordExpiryQueueTimeSlice deletes a specific record expiry queue timeslice.
func (k Keeper) DeleteRecordExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(getRecordExpiryQueueTimeKey(timestamp))
}

// InsertRecordExpiryQueue inserts a record CID to the appropriate timeslice in the record expiry queue.
func (k Keeper) InsertRecordExpiryQueue(ctx sdk.Context, val types.Record) {
	timeSlice := k.GetRecordExpiryQueueTimeSlice(ctx, val.ExpiryTime)
	timeSlice = append(timeSlice, val.ID)
	k.SetRecordExpiryQueueTimeSlice(ctx, val.ExpiryTime, timeSlice)
}

// DeleteRecordExpiryQueue deletes a record CID from the record expiry queue.
func (k Keeper) DeleteRecordExpiryQueue(ctx sdk.Context, record types.Record) {
	timeSlice := k.GetRecordExpiryQueueTimeSlice(ctx, record.ExpiryTime)
	newTimeSlice := []types.ID{}

	for _, cid := range timeSlice {
		if !bytes.Equal([]byte(cid), []byte(record.ID)) {
			newTimeSlice = append(newTimeSlice, cid)
		}
	}

	if len(newTimeSlice) == 0 {
		k.DeleteRecordExpiryQueueTimeSlice(ctx, record.ExpiryTime)
	} else {
		k.SetRecordExpiryQueueTimeSlice(ctx, record.ExpiryTime, newTimeSlice)
	}
}

// RecordExpiryQueueIterator returns all the record expiry queue timeslices from time 0 until endTime.
func (k Keeper) RecordExpiryQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	rangeEndBytes := sdk.InclusiveEndBytes(getRecordExpiryQueueTimeKey(endTime))
	return store.Iterator(prefixExpiryTimeToRecordsIndex, rangeEndBytes)
}

// GetAllExpiredRecords returns a concatenated list of all the timeslices before currTime.
func (k Keeper) GetAllExpiredRecords(ctx sdk.Context, currTime time.Time) (expiredRecordCIDs []types.ID) {
	// Gets an iterator for all timeslices from time 0 until the current block header time.
	itr := k.RecordExpiryQueueIterator(ctx, ctx.BlockHeader().Time)
	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		timeslice := []types.ID{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(itr.Value(), &timeslice)
		expiredRecordCIDs = append(expiredRecordCIDs, timeslice...)
	}

	return expiredRecordCIDs
}

// ProcessRecordExpiryQueue tries to renew expiring records (by collecting rent) else marks them as deleted.
func (k Keeper) ProcessRecordExpiryQueue(ctx sdk.Context) {
	cids := k.GetAllExpiredRecords(ctx, ctx.BlockHeader().Time)
	for _, cid := range cids {
		record := k.GetRecord(ctx, cid)

		// If record doesn't have an associated bond or if bond no longer exists, mark it deleted.
		if record.BondID == "" || !k.BondKeeper.HasBond(ctx, record.BondID) {
			record.Deleted = true
			k.PutRecord(ctx, record)
			k.DeleteRecordExpiryQueue(ctx, record)

			return
		}

		// Try to renew the record by taking rent.
		k.TryTakeRecordRent(ctx, record)
	}
}

// ProcessNameRecords creates name records.
func (k Keeper) ProcessNameRecords(ctx sdk.Context, record types.Record) {
	k.SetNameRecord(ctx, record.WRN(), record.ToNameRecord())
	k.MaybeUpdateBaseNameRecord(ctx, record)
}

// MaybeUpdateBaseNameRecord updates the base name record if required.
func (k Keeper) MaybeUpdateBaseNameRecord(ctx sdk.Context, record types.Record) {
	if !k.HasNameRecord(ctx, record.BaseWRN()) {
		// Create base name record.
		k.SetNameRecord(ctx, record.BaseWRN(), record.ToNameRecord())
		return
	}

	// Get current base record (which will have current latest version).
	baseNameRecord := k.GetNameRecord(ctx, record.BaseWRN())
	latestRecord := k.GetRecord(ctx, baseNameRecord.ID)

	latestVersion := helpers.GetSemver(latestRecord.Version())
	createdVersion := helpers.GetSemver(record.Version())
	if createdVersion.GreaterThan(latestVersion) {
		// Need to update the base name record.
		k.SetNameRecord(ctx, record.BaseWRN(), record.ToNameRecord())
	}
}

// TryTakeRecordRent tries to take rent from the record bond.
func (k Keeper) TryTakeRecordRent(ctx sdk.Context, record types.Record) {
	bondObj := k.BondKeeper.GetBond(ctx, record.BondID)
	coins, err := sdk.ParseCoins(k.RecordRent(ctx))
	if err != nil {
		panic("Invalid record rent.")
	}

	rent, err := sdk.ConvertCoin(coins[0], bond.MicroWire)
	if err != nil {
		panic("Invalid record rent.")
	}

	// Try deducting rent from bond.
	updatedBalance, isNeg := bondObj.Balance.SafeSub(sdk.NewCoins(rent))
	if isNeg {
		// Insufficient funds, mark record as deleted.
		record.Deleted = true
		k.PutRecord(ctx, record)
		k.DeleteRecordExpiryQueue(ctx, record)

		return
	}

	// Move funds from bond module to record rent module.
	err = k.BondKeeper.SupplyKeeper.SendCoinsFromModuleToModule(ctx, bond.ModuleName, bond.RecordRentModuleAccountName, sdk.NewCoins(rent))
	if err != nil {
		panic("Error withdrawing rent.")
	}

	// Update bond balance.
	bondObj.Balance = updatedBalance
	k.BondKeeper.SaveBond(ctx, bondObj)

	// Delete old expiry queue entry, create new one.
	k.DeleteRecordExpiryQueue(ctx, record)
	record.ExpiryTime = ctx.BlockHeader().Time.Add(k.RecordExpiryTime(ctx))
	k.InsertRecordExpiryQueue(ctx, record)

	// Save record.
	record.Deleted = false
	k.PutRecord(ctx, record)
	k.AddBondToRecordIndexEntry(ctx, record.BondID, record.ID)
}
