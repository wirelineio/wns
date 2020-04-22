//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/wirelineio/wns/x/nameservice"
)

// Special check for errors due to state pruning.
const statePrunedError = "proof is unexpectedly empty; ensure height has not been pruned"

// getCurrentHeight gets the current WNS block height.
func (ctx *Context) getCurrentHeight() (int64, error) {
	ctx.primaryNode.Calls++
	ctx.primaryNode.LastCalledAt = time.Now().UTC()

	// Note: Always get from primary node.
	status, err := ctx.primaryNode.Client.Status()
	if err != nil {
		ctx.primaryNode.Errors++
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

func (ctx *Context) getRandomRPCNodeHandler() *RPCNodeHandler {
	// TODO(ashwin): Make this persistent. Intelligent selection of nodes (e.g. based on QoS).
	nodes := ctx.secondaryNodes
	keys := reflect.ValueOf(nodes).MapKeys()
	address := keys[rand.Intn(len(keys))].Interface().(string)
	return nodes[address]
}

func (ctx *Context) getStoreValue(key []byte, height int64) ([]byte, error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: height,
		Prove:  true,
	}

	path := "/store/nameservice/key"
	rpcNodeHandler := ctx.getRandomRPCNodeHandler()

	rpcNodeHandler.Calls++
	rpcNodeHandler.LastCalledAt = time.Now().UTC()

	res, err := rpcNodeHandler.Client.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		rpcNodeHandler.Errors++
		return nil, err
	}

	if res.Response.IsErr() {
		rpcNodeHandler.Errors++

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

	ctx.primaryNode.Calls++
	ctx.primaryNode.LastCalledAt = time.Now().UTC()

	res, err := ctx.primaryNode.Client.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		ctx.primaryNode.Errors++
		return nil, err
	}

	var KVs []storeTypes.KVPair
	ctx.codec.MustUnmarshalBinaryLengthPrefixed(res.Response.Value, &KVs)

	return KVs, nil
}
