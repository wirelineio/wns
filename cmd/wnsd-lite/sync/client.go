//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"fmt"
	"strings"

	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/wirelineio/wns/x/nameservice"
)

// Special check for errors due to state pruning.
const statePrunedError = "proof is unexpectedly empty; ensure height has not been pruned"

// getCurrentHeight gets the current WNS block height.
func (ctx *Context) getCurrentHeight() (int64, error) {
	status, err := ctx.client.Status()
	if err != nil {
		return 0, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

func (ctx *Context) getBlockChangeset(height int64) (*nameservice.BlockChangeset, error) {
	value, err := ctx.getStoreValue(nameservice.GetBlockChangesetIndexKey(height), height)
	if err != nil {
		return nil, err
	}

	var changeset nameservice.BlockChangeset
	ctx.codec.MustUnmarshalBinaryBare(value, &changeset)

	return &changeset, nil
}

func (ctx *Context) getStoreValue(key []byte, height int64) ([]byte, error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: height,
		Prove:  true,
	}

	path := "/store/nameservice/key"
	res, err := ctx.client.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return nil, err
	}

	if res.Response.IsErr() {
		// Check if state has been pruned.
		if strings.Contains(res.Response.GetLog(), statePrunedError) {
			ctx.log.Errorln("Error fetching pruned state. Re-init sync with a recent genesis.json OR connect to a node that doesn't prune state.")
		}

		ctx.log.Panicln(res.Response)
	}

	if res.Response.Height == 0 && res.Response.Value != nil {
		ctx.log.Panicln("Invalid response height/value.")
	}

	if res.Response.Height > 0 && res.Response.Height != height {
		ctx.log.Panicln(fmt.Sprintf("Invalid response height: %d", res.Response.Height))
	}

	if res.Response.Height > 0 {
		// Note: Fails with `panic: runtime error: invalid memory address or nil pointer dereference` if called with empty response.
		err = VerifyProof(ctx, path, res.Response)
		if err != nil {
			return nil, err
		}
	}

	return res.Response.Value, nil
}

func (ctx *Context) getStoreSubspace(subspace string, key []byte, height int64) ([]storeTypes.KVPair, error) {
	opts := rpcclient.ABCIQueryOptions{Height: height}
	path := fmt.Sprintf("/store/%s/subspace", subspace)

	res, err := ctx.client.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return nil, err
	}

	var KVs []storeTypes.KVPair
	ctx.codec.MustUnmarshalBinaryLengthPrefixed(res.Response.Value, &KVs)

	return KVs, nil
}
