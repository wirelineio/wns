//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"bytes"
	"encoding/binary"
	"sort"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	set "github.com/deckarep/golang-set"
	"github.com/tendermint/go-amino"
	"github.com/wirelineio/wns/x/bond"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// PrefixCIDToRecordIndex is the prefix for CID -> Record index.
// Note: This is the primary index in the system.
// Note: Golang doesn't support const arrays.
var PrefixCIDToRecordIndex = []byte{0x00}

// PrefixNameAuthorityRecordIndex is the prefix for the name -> NameAuthority index.
var PrefixNameAuthorityRecordIndex = []byte{0x01}

// PrefixWRNToNameRecordIndex is the prefix for the WRN -> NamingRecord index.
var PrefixWRNToNameRecordIndex = []byte{0x02}

// PrefixBondIDToRecordsIndex is the prefix for the Bond ID -> [Record] index.
var PrefixBondIDToRecordsIndex = []byte{0x03}

// PrefixBlockChangesetIndex is the prefix for the block changeset index.
var PrefixBlockChangesetIndex = []byte{0x04}

// PrefixExpiryTimeToRecordsIndex is the prefix for the Expiry Time -> [Record] index.
var PrefixExpiryTimeToRecordsIndex = []byte{0x10}

// KeySyncStatus is the key for the sync status record.
// Only used by WNS lite but defined here to prevent conflicts with existing prefixes.
var KeySyncStatus = []byte{0xff}

// PrefixCIDToNamesIndex the the reverse index for naming, i.e. maps CID -> []Names.
// TODO(ashwin): Move out of WNS once we have an indexing service.
var PrefixCIDToNamesIndex = []byte{0xe0}

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	accountKeeper auth.AccountKeeper
	supplyKeeper  supply.Keeper
	recordKeeper  RecordKeeper
	bondKeeper    bond.BondClientKeeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	paramstore params.Subspace
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper, recordKeeper RecordKeeper, bondKeeper bond.BondClientKeeper, storeKey sdk.StoreKey, cdc *codec.Codec, paramstore params.Subspace) Keeper {
	return Keeper{
		accountKeeper: accountKeeper,
		supplyKeeper:  supplyKeeper,
		recordKeeper:  recordKeeper,
		bondKeeper:    bondKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
		paramstore:    paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// RecordKeeper exposes the bare minimal read-only API for other modules.
type RecordKeeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
}

// Record keeper implements the bond usage keeper interface.
var _ bond.BondUsageKeeper = (*RecordKeeper)(nil)

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

func GetCIDToNamesIndexKey(id types.ID) []byte {
	return append(PrefixCIDToNamesIndex, []byte(id)...)
}

// Generates name -> NameAuthority index key.
func GetNameAuthorityIndexKey(name string) []byte {
	return append(PrefixNameAuthorityRecordIndex, []byte(name)...)
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

// AddRecordToNameMapping adds a name to the record ID -> []names index.
func AddRecordToNameMapping(store sdk.KVStore, codec *amino.Codec, id types.ID, wrn string) {
	reverseNameIndexKey := GetCIDToNamesIndexKey(id)

	var names []string
	if store.Has(reverseNameIndexKey) {
		codec.MustUnmarshalBinaryBare(store.Get(reverseNameIndexKey), &names)
	}

	nameSet := sliceToSet(names)
	nameSet.Add(wrn)
	store.Set(reverseNameIndexKey, codec.MustMarshalBinaryBare(setToSlice(nameSet)))
}

// RemoveRecordToNameMapping removes a name from the record ID -> []names index.
func RemoveRecordToNameMapping(store sdk.KVStore, codec *amino.Codec, id types.ID, wrn string) {
	reverseNameIndexKey := GetCIDToNamesIndexKey(id)

	var names []string
	codec.MustUnmarshalBinaryBare(store.Get(reverseNameIndexKey), &names)
	nameSet := sliceToSet(names)
	nameSet.Remove(wrn)

	if nameSet.Cardinality() == 0 {
		// Delete as storing empty slice throws error from baseapp.
		store.Delete(reverseNameIndexKey)
	} else {
		store.Set(reverseNameIndexKey, codec.MustMarshalBinaryBare(setToSlice(nameSet)))
	}
}

// SetNameRecord - sets a name record.
func SetNameRecord(store sdk.KVStore, codec *amino.Codec, wrn string, id types.ID, height int64) {
	nameRecordIndexKey := GetNameRecordIndexKey(wrn)

	var nameRecord types.NameRecord
	if store.Has(nameRecordIndexKey) {
		bz := store.Get(nameRecordIndexKey)
		codec.MustUnmarshalBinaryBare(bz, &nameRecord)
		nameRecord.History = append(nameRecord.History, nameRecord.NameRecordEntry)

		// Update old CID -> []Name index.
		if nameRecord.NameRecordEntry.ID != "" {
			RemoveRecordToNameMapping(store, codec, nameRecord.NameRecordEntry.ID, wrn)
		}
	}

	nameRecord.NameRecordEntry = types.NameRecordEntry{
		ID:     id,
		Height: height,
	}

	store.Set(nameRecordIndexKey, codec.MustMarshalBinaryBare(nameRecord))

	// Update new CID -> []Name index.
	if id != "" {
		AddRecordToNameMapping(store, codec, id, wrn)
	}
}

// SetNameRecord - sets a name record.
func (k Keeper) SetNameRecord(ctx sdk.Context, wrn string, id types.ID) {
	SetNameRecord(ctx.KVStore(k.storeKey), k.cdc, wrn, id, ctx.BlockHeight())

	// Update changeset for name.
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

	return recordObjToRecord(store, codec, obj)
}

// GetNameRecord - gets a name record from the store.
func GetNameRecord(store sdk.KVStore, codec *amino.Codec, wrn string) *types.NameRecord {
	nameRecordKey := GetNameRecordIndexKey(wrn)
	if !store.Has(nameRecordKey) {
		return nil
	}

	bz := store.Get(nameRecordKey)
	var obj types.NameRecord
	codec.MustUnmarshalBinaryBare(bz, &obj)

	return &obj
}

// GetNameRecord - gets a name record from the store.
func (k Keeper) GetNameRecord(ctx sdk.Context, wrn string) *types.NameRecord {
	return GetNameRecord(ctx.KVStore(k.storeKey), k.cdc, wrn)
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
			records = append(records, recordObjToRecord(store, k.cdc, obj))
		}
	}

	return records
}

// ListNameAuthorityRecords - get all name authority records.
func (k Keeper) ListNameAuthorityRecords(ctx sdk.Context) map[string]types.NameAuthority {
	nameAuthorityRecords := make(map[string]types.NameAuthority)

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixNameAuthorityRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var record types.NameAuthority
			k.cdc.MustUnmarshalBinaryBare(bz, &record)
			nameAuthorityRecords[string(itr.Key()[len(PrefixNameAuthorityRecordIndex):])] = record
		}
	}

	return nameAuthorityRecords
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
func (k Keeper) ResolveWRN(ctx sdk.Context, wrn string) *types.Record {
	return ResolveWRN(ctx.KVStore(k.storeKey), k.cdc, wrn)
}

// ResolveWRN resolves a WRN to a record.
func ResolveWRN(store sdk.KVStore, codec *amino.Codec, wrn string) *types.Record {
	nameKey := GetNameRecordIndexKey(wrn)

	if store.Has(nameKey) {
		bz := store.Get(nameKey)
		var obj types.NameRecord
		codec.MustUnmarshalBinaryBare(bz, &obj)

		recordExists := HasRecord(store, obj.ID)
		if !recordExists || obj.ID == "" {
			return nil
		}

		record := GetRecord(store, codec, obj.ID)
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
			record := recordObjToRecord(store, codec, obj)
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
			records = append(records, recordObjToRecord(store, k.cdc, obj))
		}
	}

	return records
}

// ModuleName returns the module name.
func (k RecordKeeper) ModuleName() string {
	return types.ModuleName
}

// UsesBond returns true if the bond has associated records.
func (k RecordKeeper) UsesBond(ctx sdk.Context, bondID bond.ID) bool {
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
		if record.BondID == "" || !k.bondKeeper.HasBond(ctx, record.BondID) {
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
	rent, err := sdk.ParseCoins(k.RecordRent(ctx))
	if err != nil {
		panic("Invalid record rent.")
	}

	sdkErr := k.bondKeeper.TransferCoinsToModuleAccount(ctx, record.BondID, types.RecordRentModuleAccountName, rent)
	if sdkErr != nil {
		// Insufficient funds, mark record as deleted.
		record.Deleted = true
		k.PutRecord(ctx, record)
		k.DeleteRecordExpiryQueue(ctx, record)

		return
	}

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

func (k Keeper) updateBlockChangesetForNameAuthority(ctx sdk.Context, name string) {
	changeset := k.getOrCreateBlockChangeset(ctx, ctx.BlockHeight())
	changeset.NameAuthorities = append(changeset.NameAuthorities, name)
	k.saveBlockChangeset(ctx, changeset)
}

// HasNameAuthority - checks if a name/authority exists.
func (k Keeper) HasNameAuthority(ctx sdk.Context, name string) bool {
	return HasNameAuthority(ctx.KVStore(k.storeKey), name)
}

// HasNameAuthority - checks if a name authority entry exists.
func HasNameAuthority(store sdk.KVStore, name string) bool {
	return store.Has(GetNameAuthorityIndexKey(name))
}

// SetNameAuthority creates the NameAutority record.
func (k Keeper) SetNameAuthority(ctx sdk.Context, name string, ownerAddress string, ownerPublicKey string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetNameAuthorityIndexKey(name), k.cdc.MustMarshalBinaryBare(
		types.NameAuthority{
			OwnerAddress:   ownerAddress,
			OwnerPublicKey: ownerPublicKey,
			Height:         ctx.BlockHeight(),
		}))
	k.updateBlockChangesetForNameAuthority(ctx, name)
}

// GetNameAuthority - gets a name authority from the store.
func GetNameAuthority(store sdk.KVStore, codec *amino.Codec, name string) *types.NameAuthority {
	authorityKey := GetNameAuthorityIndexKey(name)
	if !store.Has(authorityKey) {
		return nil
	}

	bz := store.Get(authorityKey)
	var obj types.NameAuthority
	codec.MustUnmarshalBinaryBare(bz, &obj)

	return &obj
}

// GetNameAuthority - gets a name authority from the store.
func (k Keeper) GetNameAuthority(ctx sdk.Context, name string) *types.NameAuthority {
	return GetNameAuthority(ctx.KVStore(k.storeKey), k.cdc, name)
}

func recordObjToRecord(store sdk.KVStore, codec *amino.Codec, obj types.RecordObj) types.Record {
	record := obj.ToRecord()

	reverseNameIndexKey := GetCIDToNamesIndexKey(obj.ID)
	if store.Has(reverseNameIndexKey) {
		var names []string
		codec.MustUnmarshalBinaryBare(store.Get(reverseNameIndexKey), &names)
		record.Names = names
	}

	return record
}

// GetModuleBalances gets the nameservice module account(s) balances.
func (k Keeper) GetModuleBalances(ctx sdk.Context) map[string]sdk.Coins {
	balances := map[string]sdk.Coins{}
	accountNames := []string{types.RecordRentModuleAccountName}

	for _, accountName := range accountNames {
		moduleAddress := k.supplyKeeper.GetModuleAddress(accountName)
		moduleAccount := k.accountKeeper.GetAccount(ctx, moduleAddress)
		if moduleAccount != nil {
			balances[accountName] = moduleAccount.GetCoins()
		}
	}

	return balances
}

func setToSlice(set set.Set) []string {
	names := []string{}

	for name := range set.Iter() {
		if name, ok := name.(string); ok && name != "" {
			names = append(names, name)
		}
	}

	sort.SliceStable(names, func(i, j int) bool { return names[i] < names[j] })

	return names
}

func sliceToSet(names []string) set.Set {
	set := set.NewThreadUnsafeSet()

	for _, name := range names {
		if name != "" {
			set.Add(name)
		}
	}

	return set
}
