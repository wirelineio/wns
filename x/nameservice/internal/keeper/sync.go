//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
	wnsUtils "github.com/wirelineio/wns/utils"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

func GetBlockChangesetIndexKey(height int64) []byte {
	return append(PrefixBlockChangesetIndex, wnsUtils.Int64ToBytes(height)...)
}

func getOrCreateBlockChangeset(ctx sdk.Context, store sdk.KVStore, codec *amino.Codec, height int64) *types.BlockChangeset {

	bz := store.Get(GetBlockChangesetIndexKey(height))

	if bz != nil {
		var changeset types.BlockChangeset
		codec.MustUnmarshalBinaryBare(bz, &changeset)

		return &changeset
	}

	return &types.BlockChangeset{
		Height:  height,
		Records: []types.ID{},
		Names:   []string{},
	}
}

func (k Keeper) getOrCreateBlockChangeset(ctx sdk.Context, height int64) *types.BlockChangeset {
	return getOrCreateBlockChangeset(ctx, ctx.KVStore(k.storeKey), k.cdc, height)
}

func saveBlockChangeset(ctx sdk.Context, store sdk.KVStore, codec *amino.Codec, changeset *types.BlockChangeset) {
	bz := codec.MustMarshalBinaryBare(*changeset)
	store.Set(GetBlockChangesetIndexKey(changeset.Height), bz)
}

func (k Keeper) saveBlockChangeset(ctx sdk.Context, changeset *types.BlockChangeset) {
	saveBlockChangeset(ctx, ctx.KVStore(k.storeKey), k.cdc, changeset)
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

func updateBlockChangesetForNameAuthority(ctx sdk.Context, store sdk.KVStore, codec *amino.Codec, name string) {
	changeset := getOrCreateBlockChangeset(ctx, store, codec, ctx.BlockHeight())
	changeset.NameAuthorities = append(changeset.NameAuthorities, name)
	saveBlockChangeset(ctx, store, codec, changeset)
}
