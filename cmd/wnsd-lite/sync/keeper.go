//
// Copyright 2020 Wireline, Inc.
//

package sync

import "github.com/wirelineio/wns/x/nameservice"

// Keeper is an impl. of an interface similar to the nameservice Keeper.
type Keeper struct {
	Context *Context
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(id nameservice.ID) bool {
	store := k.Context.DBStore
	return store.Has(nameservice.GetRecordIndexKey(id))
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(id nameservice.ID) nameservice.Record {
	store := k.Context.DBStore

	bz := store.Get(nameservice.GetRecordIndexKey(id))
	var obj nameservice.RecordObj
	k.Context.Codec.MustUnmarshalBinaryBare(bz, &obj)

	return obj.ToRecord()
}
