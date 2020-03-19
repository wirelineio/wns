//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	ns "github.com/wirelineio/wns/x/nameservice"
)

// Keeper is an impl. of an interface similar to the nameservice Keeper.
type Keeper struct {
	Context *Context
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(id ns.ID) bool {
	return ns.HasRecord(k.Context.DBStore, id)
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(id ns.ID) ns.Record {
	return ns.GetRecord(k.Context.DBStore, k.Context.Codec, id)
}

// ResolveWRN resolves a WRN to a record.
// Note: Version part of the WRN might have a semver range.
func (k Keeper) ResolveWRN(wrn string) *ns.Record {
	return ns.ResolveWRN(k.Context.DBStore, k.Context.Codec, wrn)
}
