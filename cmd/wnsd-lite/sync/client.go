//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"fmt"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/wirelineio/wns/x/nameservice"
)

// getCurrentHeight gets the current WNS block height.
func (ctx *Context) getCurrentHeight() (int64, error) {
	status, err := ctx.Client.Status()
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
	ctx.Codec.MustUnmarshalBinaryBare(value, &changeset)

	return &changeset, nil
}

func (ctx *Context) getStoreValue(key []byte, height int64) ([]byte, error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: height,
		Prove:  true,
	}

	path := "/store/nameservice/key"
	res, err := ctx.Client.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return nil, err
	}

	if res.Response.Height == 0 && res.Response.Value != nil {
		panic("Invalid response height/value.")
	}

	if res.Response.Height > 0 && res.Response.Height != height {
		panic(fmt.Sprintf("Invalid response height: %d", res.Response.Height))
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
