//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

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
