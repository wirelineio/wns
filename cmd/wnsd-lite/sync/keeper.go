//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/go-amino"
	"github.com/wirelineio/wns/x/nameservice"
	ns "github.com/wirelineio/wns/x/nameservice"
)

// Keeper is an impl. of an interface similar to the nameservice Keeper.
type Keeper struct {
	config *Config
	codec  *amino.Codec
	store  store.KVStore
}

// NewKeeper creates a new keeper.
func NewKeeper(ctx *Context) *Keeper {
	return &Keeper{config: ctx.config, codec: ctx.codec, store: ctx.store}
}

// Status represents the sync status of the node.
type Status struct {
	LastSyncedHeight int64
	CatchingUp       bool
}

// GetChainID gets the chain ID.
func (k Keeper) GetChainID() string {
	return k.config.ChainID
}

// HasStatusRecord checks if the store has a status record.
func (k Keeper) HasStatusRecord() bool {
	return k.store.Has(ns.KeySyncStatus)
}

// GetStatusRecord gets the sync status record.
func (k Keeper) GetStatusRecord() Status {
	bz := k.store.Get(ns.KeySyncStatus)
	var status Status
	k.codec.MustUnmarshalBinaryBare(bz, &status)

	return status
}

// SaveStatus saves the sync status record.
func (k Keeper) SaveStatus(status Status) {
	bz := k.codec.MustMarshalBinaryBare(status)
	k.store.Set(ns.KeySyncStatus, bz)
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(id ns.ID) bool {
	return ns.HasRecord(k.store, id)
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(id ns.ID) ns.Record {
	return ns.GetRecord(k.store, k.codec, id)
}

// PutRecord - saves a record to the store and updates ID -> Record index.
func (k Keeper) PutRecord(record nameservice.RecordObj) {
	k.store.Set(nameservice.GetRecordIndexKey(record.ID), k.codec.MustMarshalBinaryBare(record))
}

// SetNameAuthorityRecord - sets a name authority record.
func (k Keeper) SetNameAuthorityRecord(name string, nameAuthority nameservice.NameAuthority) {
	k.store.Set(nameservice.GetNameAuthorityIndexKey(name), k.codec.MustMarshalBinaryBare(nameAuthority))
}

// SetNameRecord - sets a name record.
func (k Keeper) SetNameRecord(wrn string, nameRecord nameservice.NameRecord) {
	k.store.Set(nameservice.GetNameRecordIndexKey(wrn), k.codec.MustMarshalBinaryBare(nameRecord))
}

// ResolveWRN resolves a WRN to a record.
// Note: Version part of the WRN might have a semver range.
func (k Keeper) ResolveWRN(wrn string) *ns.Record {
	return ns.ResolveWRN(k.store, k.codec, wrn)
}

// MatchRecords - get all matching records.
func (k Keeper) MatchRecords(matchFn func(*ns.Record) bool) []*ns.Record {
	return ns.MatchRecords(k.store, k.codec, matchFn)
}
