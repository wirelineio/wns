//
// Copyright 2020 Wireline, Inc.
//

package sync

import ns "github.com/wirelineio/wns/x/nameservice"

// Keeper is an impl. of an interface similar to the nameservice Keeper.
type Keeper struct {
	Context *Context
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(id ns.ID) bool {
	store := k.Context.DBStore
	return ns.HasRecord(store, id)
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(id ns.ID) ns.Record {
	store := k.Context.DBStore
	return ns.GetRecord(store, k.Context.Codec, id)
}
