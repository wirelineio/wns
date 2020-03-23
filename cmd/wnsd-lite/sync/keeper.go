//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/go-amino"
	ns "github.com/wirelineio/wns/x/nameservice"
)

// Keeper is an impl. of an interface similar to the nameservice Keeper.
type Keeper struct {
	Codec *amino.Codec
	Store store.KVStore
}

// NewKeeper creates a new keeper.
func NewKeeper(codec *amino.Codec, store store.KVStore) *Keeper {
	return &Keeper{Codec: codec, Store: store}
}

// Status represents the sync status of the node.
type Status struct {
	LastSyncedHeight int64
}

// HasStatusRecord checks if the store has a status record.
func (k Keeper) HasStatusRecord() bool {
	return k.Store.Has(ns.KeySyncStatus)
}

// GetStatusRecord gets the sync status record.
func (k Keeper) GetStatusRecord() Status {
	bz := k.Store.Get(ns.KeySyncStatus)
	var status Status
	k.Codec.MustUnmarshalBinaryBare(bz, &status)

	return status
}

// SaveStatus saves the sync status record.
func (k Keeper) SaveStatus(status Status) {
	bz := k.Codec.MustMarshalBinaryBare(status)
	k.Store.Set(ns.KeySyncStatus, bz)
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(id ns.ID) bool {
	return ns.HasRecord(k.Store, id)
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(id ns.ID) ns.Record {
	return ns.GetRecord(k.Store, k.Codec, id)
}

// ResolveWRN resolves a WRN to a record.
// Note: Version part of the WRN might have a semver range.
func (k Keeper) ResolveWRN(wrn string) *ns.Record {
	return ns.ResolveWRN(k.Store, k.Codec, wrn)
}

// MatchRecords - get all matching records.
func (k Keeper) MatchRecords(matchFn func(*ns.Record) bool) []*ns.Record {
	return ns.MatchRecords(k.Store, k.Codec, matchFn)
}
