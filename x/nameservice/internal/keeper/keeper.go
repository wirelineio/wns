//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"bytes"
	"encoding/binary"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/go-amino"
	"github.com/wirelineio/wns/x/bond"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// PrefixCIDToRecordIndex is the prefix for CID -> Record index.
// Note: This is the primary index in the system.
// Note: Golang doesn't support const arrays.
var PrefixCIDToRecordIndex = []byte{0x00}

// PrefixWRNToNameRecordIndex is the prefix for the WRN -> NamingRecord index.
var PrefixWRNToNameRecordIndex = []byte{0x01}

// PrefixBondIDToRecordsIndex is the prefix for the Bond ID -> [Record] index.
var PrefixBondIDToRecordsIndex = []byte{0x03}

// PrefixBlockChangesetIndex is the prefix for the block changeset index.
var PrefixBlockChangesetIndex = []byte{0x04}

// PrefixExpiryTimeToRecordsIndex is the prefix for the Expiry Time -> [Record] index.
var PrefixExpiryTimeToRecordsIndex = []byte{0x10}

// KeySyncStatus is the key for the sync status record.
// Only used by WNS lite but defined here to prevent conflicts with existing prefixes.
var KeySyncStatus = []byte{0xff}

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
	store.Set(GetRecordIndexKey(record.ID), k.cdc.MustMarshalBinaryBare(record.ToRecordObj()))
	k.updateBlockChangesetForRecord(ctx, record.ID)
}

// Generates Bond ID -> Bond index key.
func GetRecordIndexKey(id types.ID) []byte {
	return append(PrefixCIDToRecordIndex, []byte(id)...)
}

// Generates WRN -> NameRecord index key.
func GetNameRecordIndexKey(wrn string) []byte {
	return append(PrefixWRNToNameRecordIndex, []byte(wrn)...)
}

// Generates Bond ID -> Records index key.
func getBondIDToRecordsIndexKey(bondID bond.ID, id types.ID) []byte {
	return append(append(PrefixBondIDToRecordsIndex, []byte(bondID)...), []byte(id)...)
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
	store.Set(GetNameRecordIndexKey(wrn), k.cdc.MustMarshalBinaryBare(nameRecord))
	k.updateBlockChangesetForName(ctx, wrn)
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(ctx sdk.Context, id types.ID) bool {
	return HasRecord(ctx.KVStore(k.storeKey), id)
}

// HasRecord - checks if a record by the given ID exists.
func HasRecord(store sdk.KVStore, id types.ID) bool {
	return store.Has(GetRecordIndexKey(id))
}

// HasNameRecord - checks if a name record exists.
func (k Keeper) HasNameRecord(ctx sdk.Context, wrn string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetNameRecordIndexKey(wrn))
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(ctx sdk.Context, id types.ID) types.Record {
	return GetRecord(ctx.KVStore(k.storeKey), k.cdc, id)
}

// GetRecord - gets a record from the store.
func GetRecord(store sdk.KVStore, codec *amino.Codec, id types.ID) types.Record {
	bz := store.Get(GetRecordIndexKey(id))
	var obj types.RecordObj
	codec.MustUnmarshalBinaryBare(bz, &obj)

	return obj.ToRecord()
}

// GetNameRecord - gets a name record from the store.
func (k Keeper) GetNameRecord(ctx sdk.Context, wrn string) types.NameRecord {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetNameRecordIndexKey(wrn))
	var obj types.NameRecord
	k.cdc.MustUnmarshalBinaryBare(bz, &obj)

	return obj
}

// ListRecords - get all records.
func (k Keeper) ListRecords(ctx sdk.Context) []types.Record {
	var records []types.Record

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixCIDToRecordIndex)
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
	itr := sdk.KVStorePrefixIterator(store, PrefixWRNToNameRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var record types.NameRecord
			k.cdc.MustUnmarshalBinaryBare(bz, &record)
			nameRecords[string(itr.Key()[len(PrefixWRNToNameRecordIndex):])] = record
		}
	}

	return nameRecords
}

// ResolveWRN resolves a WRN to a record.
// Note: Version part of the WRN might have a semver range.
func (k Keeper) ResolveWRN(ctx sdk.Context, wrn string) *types.Record {
	return ResolveWRN(ctx.KVStore(k.storeKey), k.cdc, wrn)
}

// ResolveWRN resolves a WRN to a record.
// Note: Version part of the WRN might have a semver range.
func ResolveWRN(store sdk.KVStore, codec *amino.Codec, wrn string) *types.Record {
	segments := strings.Split(wrn, "#")
	if len(segments) == 2 {
		baseWRN, semver := segments[0], segments[1]
		if strings.ContainsAny(semver, "^~<>=!") {
			// Handle semver range.
			return ResolveBaseWRN(store, codec, baseWRN, semver)
		}
	}

	return ResolveFullWRN(store, codec, wrn)
}

// ResolveFullWRN resolves a WRN (full path) to a record.
// Note: Version part of the WRN MUST NOT have a semver range.
func (k Keeper) ResolveFullWRN(ctx sdk.Context, wrn string) *types.Record {
	return ResolveFullWRN(ctx.KVStore(k.storeKey), k.cdc, wrn)
}

// ResolveFullWRN resolves a WRN (full path) to a record.
// Note: Version part of the WRN MUST NOT have a semver range.
func ResolveFullWRN(store sdk.KVStore, codec *amino.Codec, wrn string) *types.Record {
	nameKey := GetNameRecordIndexKey(wrn)

	if store.Has(nameKey) {
		bz := store.Get(nameKey)
		var obj types.NameRecord
		codec.MustUnmarshalBinaryBare(bz, &obj)

		record := GetRecord(store, codec, obj.ID)
		return &record
	}

	return nil
}

// ResolveBaseWRN resolves a BaseWRN + semver range to a record (picks the highest matching version).
func (k Keeper) ResolveBaseWRN(ctx sdk.Context, baseWRN string, semverRange string) *types.Record {
	return ResolveBaseWRN(ctx.KVStore(k.storeKey), k.cdc, baseWRN, semverRange)
}

// ResolveBaseWRN resolves a BaseWRN + semver range to a record (picks the highest matching version).
func ResolveBaseWRN(store sdk.KVStore, codec *amino.Codec, baseWRN string, semverRange string) *types.Record {
	semverConstraint, err := semver.NewConstraint(semverRange)
	if err != nil {
		// Handle constraint not being parsable.
		return nil
	}

	baseNameKey := GetNameRecordIndexKey(baseWRN)
	if !store.Has(baseNameKey) {
		return nil
	}

	var highestSemver, _ = semver.NewVersion("0.0.0")
	var highestNameRecord types.NameRecord = types.NameRecord{}

	itr := sdk.KVStorePrefixIterator(store, baseNameKey)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var record types.NameRecord
			codec.MustUnmarshalBinaryBare(bz, &record)

			semver, err := semver.NewVersion(record.Version)
			if err == nil && semverConstraint.Check(semver) && semver.GreaterThan(highestSemver) {
				highestSemver = semver
				highestNameRecord = record
			}
		}
	}

	if highestNameRecord.ID != "" {
		record := GetRecord(store, codec, highestNameRecord.ID)
		return &record
	}

	return nil
}

// MatchRecords - get all matching records.
func (k Keeper) MatchRecords(ctx sdk.Context, matchFn func(*types.Record) bool) []*types.Record {
	return MatchRecords(ctx.KVStore(k.storeKey), k.cdc, matchFn)
}

// MatchRecords - get all matching records.
func MatchRecords(store sdk.KVStore, codec *amino.Codec, matchFn func(*types.Record) bool) []*types.Record {
	var records []*types.Record

	itr := sdk.KVStorePrefixIterator(store, PrefixCIDToRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.RecordObj
			codec.MustUnmarshalBinaryBare(bz, &obj)
			record := obj.ToRecord()
			if matchFn(&record) {
				records = append(records, &record)
			}
		}
	}

	return records
}

// QueryRecordsByBond - get all records for the given bond.
func (k RecordKeeper) QueryRecordsByBond(ctx sdk.Context, bondID bond.ID) []types.Record {
	var records []types.Record

	bondIDPrefix := append(PrefixBondIDToRecordsIndex, []byte(bondID)...)
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, bondIDPrefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		cid := itr.Key()[len(bondIDPrefix):]
		bz := store.Get(append(PrefixCIDToRecordIndex, cid...))
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
	bondIDPrefix := append(PrefixBondIDToRecordsIndex, []byte(bondID)...)
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, bondIDPrefix)
	defer itr.Close()
	return itr.Valid()
}

// getRecordExpiryQueueTimeKey gets the prefix for the record expiry queue.
func getRecordExpiryQueueTimeKey(timestamp time.Time) []byte {
	timeBytes := sdk.FormatTimeBytes(timestamp)
	return append(PrefixExpiryTimeToRecordsIndex, timeBytes...)
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
	return store.Iterator(PrefixExpiryTimeToRecordsIndex, rangeEndBytes)
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

func int64ToBytes(num int64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, num)
	return buf.Bytes()
}

func GetBlockChangesetIndexKey(height int64) []byte {
	return append(PrefixBlockChangesetIndex, int64ToBytes(height)...)
}

func (k Keeper) getOrCreateBlockChangeset(ctx sdk.Context, height int64) *types.BlockChangeset {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetBlockChangesetIndexKey(height))

	if bz != nil {
		var changeset types.BlockChangeset
		k.cdc.MustUnmarshalBinaryBare(bz, &changeset)

		return &changeset
	}

	return &types.BlockChangeset{
		Height:  height,
		Records: []types.ID{},
		Names:   []string{},
	}
}

// SetRecordExpiryQueueTimeSlice sets a specific record expiry queue timeslice.
func (k Keeper) saveBlockChangeset(ctx sdk.Context, changeset *types.BlockChangeset) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(*changeset)
	store.Set(GetBlockChangesetIndexKey(changeset.Height), bz)
}

func (k Keeper) updateBlockChangesetForRecord(ctx sdk.Context, id types.ID) {
	changeset := k.getOrCreateBlockChangeset(ctx, ctx.BlockHeight())
	changeset.Records = append(changeset.Records, id)
	k.saveBlockChangeset(ctx, changeset)
}

func (k Keeper) updateBlockChangesetForName(ctx sdk.Context, wrn string) {
	changeset := k.getOrCreateBlockChangeset(ctx, ctx.BlockHeight())
	changeset.Names = append(changeset.Names, wrn)
	k.saveBlockChangeset(ctx, changeset)
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

	// Clear bonds.
	k.BondKeeper.Clear(ctx)
}

// HasNameAuthority - checks if a name/authority exists.
func (k Keeper) HasNameAuthority(ctx sdk.Context, wrn string) bool {
	return HasNameAuthority(ctx.KVStore(k.storeKey), wrn)
}

// HasNameAuthority - checks if a name authority entry exists.
func HasNameAuthority(store sdk.KVStore, wrn string) bool {
	return store.Has(GetNameRecordIndexKey(wrn))
}
