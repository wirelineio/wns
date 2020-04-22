//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"time"

	nameservice "github.com/wirelineio/wns/x/nameservice"
)

// AggressiveSyncIntervalInMillis is the interval for aggressive sync, to catch up quickly to the current height.
const AggressiveSyncIntervalInMillis = 250

// SyncIntervalInMillis is the interval for initiating incremental sync, when already caught up to current height.
const SyncIntervalInMillis = 5 * 1000

// ErrorWaitDurationMillis is the wait duration in case of errors.
const ErrorWaitDurationMillis = 1 * 1000

// SyncLaggingMinHeightDiff is the min. difference in height to consider a lite node as lagging the full node.
const SyncLaggingMinHeightDiff = 5

// DumpRPCNodeStatsFrequencyMillis controls frequency to dump RPC node stats.
const DumpRPCNodeStatsFrequencyMillis = 60 * 1000

// DiscoverRPCNodesFrequencyMillis controls frequency to discover new RPC endpoints.
const DiscoverRPCNodesFrequencyMillis = 60 * 1000

// Init sets up the lite node.
func Init(ctx *Context, height int64) {
	// If sync record exists, abort with error.
	if ctx.keeper.HasStatusRecord() {
		ctx.log.Fatalln("Node already initialized, aborting.")
	}

	if !ctx.config.InitFromNode && !ctx.config.InitFromGenesisFile {
		ctx.log.Fatalln("Must pass one of `--from-node` and `--from-genesis-file`.")
	}

	if ctx.config.InitFromNode {
		initFromNode(ctx)
	} else if ctx.config.InitFromGenesisFile {
		initFromGenesisFile(ctx, height)
	}
}

// Start initiates the sync process.
func Start(ctx *Context) {
	// Fail if node has no sync status record.
	if !ctx.keeper.HasStatusRecord() {
		ctx.log.Fatalln("Node not initialized, aborting.")
	}

	go dumpConnectionStatsOnTimer(ctx)

	if ctx.config.Endpoint != "" {
		go discoverRPCNodesOnTimer(ctx)
		ctx.log.Infoln("RPC endpoint discovery ON:", ctx.config.Endpoint)
	} else {
		ctx.log.Infoln("RPC endpoint discovery OFF.")
	}

	syncStatus := ctx.keeper.GetStatusRecord()
	lastSyncedHeight := syncStatus.LastSyncedHeight

	for {
		chainCurrentHeight, err := ctx.primaryNode.getCurrentHeight()
		if err != nil {
			logErrorAndWait(ctx, err)
			continue
		}

		if lastSyncedHeight > chainCurrentHeight {
			// Maybe we've connected to a new primary node (after restart) and that isn't fully caught up, yet. Just wait.
			logErrorAndWait(ctx, errors.New("last synced height greater than current chain height"))
			continue
		}

		newSyncHeight := lastSyncedHeight + 1
		if newSyncHeight > chainCurrentHeight {
			// Can't sync beyond chain height, just wait.
			waitAfterSync(chainCurrentHeight, chainCurrentHeight)
			continue
		}

		err = syncAtHeight(ctx, newSyncHeight)
		if err != nil {
			logErrorAndWait(ctx, err)
			continue
		}

		// Saved last synced height in db.
		lastSyncedHeight = newSyncHeight
		catchingUp := (chainCurrentHeight - lastSyncedHeight) > SyncLaggingMinHeightDiff

		ctx.keeper.SaveStatus(Status{
			LastSyncedHeight: lastSyncedHeight,
			CatchingUp:       catchingUp,
		})

		waitAfterSync(chainCurrentHeight, lastSyncedHeight)
	}
}

// syncAtHeight runs a sync cycle for the given height.
func syncAtHeight(ctx *Context, height int64) error {
	rpcNodeHandler := getRandomRPCNodeHandler(ctx)

	ctx.log.Infoln("Syncing from", rpcNodeHandler.Address, "at height:", height)

	changeset, err := rpcNodeHandler.getBlockChangeset(ctx, height)
	if err != nil {
		return err
	}

	if changeset.Height <= 0 {
		// No changeset for this block, ignore.
		return nil
	}

	ctx.log.Debugln("Syncing changeset:", changeset)

	// Sync records.
	err = rpcNodeHandler.syncRecords(ctx, height, changeset.Records)
	if err != nil {
		return err
	}

	// Sync name records.
	err = rpcNodeHandler.syncNameRecords(ctx, height, changeset.Names)
	if err != nil {
		return err
	}

	// Flush cache changes to underlying store.
	ctx.cache.Write()

	return nil
}

func (rpcNodeHandler *RPCNodeHandler) syncRecords(ctx *Context, height int64, records []nameservice.ID) error {
	for _, id := range records {
		recordKey := nameservice.GetRecordIndexKey(id)
		value, err := rpcNodeHandler.getStoreValue(ctx, recordKey, height)
		if err != nil {
			return err
		}

		ctx.cache.Set(recordKey, value)
	}

	return nil
}

func (rpcNodeHandler *RPCNodeHandler) syncNameRecords(ctx *Context, height int64, names []string) error {
	for _, name := range names {
		nameRecordKey := nameservice.GetNameRecordIndexKey(name)
		value, err := rpcNodeHandler.getStoreValue(ctx, nameRecordKey, height)
		if err != nil {
			return err
		}

		ctx.cache.Set(nameRecordKey, value)
	}

	return nil
}

func waitAfterSync(chainCurrentHeight int64, lastSyncedHeight int64) {
	if chainCurrentHeight == lastSyncedHeight {
		// Caught up to current chain height, don't have to poll aggressively now.
		time.Sleep(SyncIntervalInMillis * time.Millisecond)
	} else {
		// Still catching up to current height, poll more aggressively.
		time.Sleep(AggressiveSyncIntervalInMillis * time.Millisecond)
	}
}

func logErrorAndWait(ctx *Context, err error) {
	ctx.log.Errorln(err)

	// TODO(ashwin): Exponential backoff logic.
	time.Sleep(ErrorWaitDurationMillis * time.Millisecond)
}

func initFromNode(ctx *Context) {
	height, err := ctx.primaryNode.getCurrentHeight()
	if err != nil {
		ctx.log.Fatalln("Error fetching current height:", err)
	}

	ctx.log.Debugln("Current block height:", height)

	recordKVs, err := ctx.getStoreSubspace("nameservice", nameservice.PrefixCIDToRecordIndex, height)
	if err != nil {
		ctx.log.Fatalln("Error fetching records", err)
	}

	for _, kv := range recordKVs {
		var record nameservice.RecordObj
		ctx.codec.MustUnmarshalBinaryBare(kv.Value, &record)
		ctx.log.Debugln("Importing record", record.ID)
		ctx.keeper.PutRecord(record)
	}

	namesKVs, err := ctx.getStoreSubspace("nameservice", nameservice.PrefixWRNToNameRecordIndex, height)
	if err != nil {
		ctx.log.Fatalln("Error fetching name records", err)
	}

	for _, kv := range namesKVs {
		var nameRecord nameservice.NameRecord
		ctx.codec.MustUnmarshalBinaryBare(kv.Value, &nameRecord)
		wrn := string(kv.Key[len(nameservice.PrefixWRNToNameRecordIndex):])
		ctx.log.Debugln("Importing name", wrn)
		ctx.keeper.SetNameRecord(wrn, nameRecord)
	}

	// Create sync status record.
	ctx.keeper.SaveStatus(Status{LastSyncedHeight: height})
}

func initFromGenesisFile(ctx *Context, height int64) {
	// Create <home>/config directory if it doesn't exist.
	configDirPath := filepath.Join(ctx.config.Home, "config")
	os.Mkdir(configDirPath, 0755)

	// Import genesis.json.
	genesisJSONPath := filepath.Join(configDirPath, "genesis.json")
	_, err := os.Stat(genesisJSONPath)
	if err != nil {
		ctx.log.Fatalln("Genesis file error:", err)
	}

	geneisState := GenesisState{}
	bytes, err := ioutil.ReadFile(genesisJSONPath)
	if err != nil {
		ctx.log.Fatalln(err)
	}

	err = ctx.codec.UnmarshalJSON(bytes, &geneisState)
	if err != nil {
		ctx.log.Fatalln(err)
	}

	// Check that chain-id matches.
	if geneisState.ChainID != ctx.config.ChainID {
		ctx.log.Fatalln("Chain ID mismatch:", genesisJSONPath)
	}

	names := geneisState.AppState.Nameservice.Names
	for _, nameEntry := range names {
		ctx.keeper.SetNameRecord(nameEntry.Name, nameEntry.Entry)
	}

	records := geneisState.AppState.Nameservice.Records
	for _, record := range records {
		ctx.keeper.PutRecord(record)
	}

	// Create sync status record.
	ctx.keeper.SaveStatus(Status{LastSyncedHeight: height})
}

func getRandomRPCNodeHandler(ctx *Context) *RPCNodeHandler {
	ctx.nodeLock.RLock()
	defer ctx.nodeLock.RUnlock()

	// TODO(ashwin): Make address book persistent. Intelligent selection of nodes (e.g. based on QoS).
	nodes := ctx.secondaryNodes
	keys := reflect.ValueOf(nodes).MapKeys()
	address := keys[rand.Intn(len(keys))].Interface().(string)
	rpcNodeHandler := nodes[address]

	return rpcNodeHandler
}

func dumpConnectionStatsOnTimer(ctx *Context) {
	for {
		time.Sleep(DumpRPCNodeStatsFrequencyMillis * time.Millisecond)
		dumpConnectionStats(ctx)
	}
}

func dumpConnectionStats(ctx *Context) {
	ctx.nodeLock.RLock()
	defer ctx.nodeLock.RUnlock()

	// Log RPC node stats.
	bytes, _ := json.Marshal(ctx.secondaryNodes)
	ctx.log.Debugln(string(bytes))
}

func discoverRPCNodesOnTimer(ctx *Context) {
	for {
		time.Sleep(DiscoverRPCNodesFrequencyMillis * time.Millisecond)
		discoverRPCNodes(ctx)
	}
}

// Discover new RPC nodes.
func discoverRPCNodes(ctx *Context) {
	rpcEndpoints, err := DiscoverRPCEndpoints(ctx, ctx.config.Endpoint)
	if err != nil {
		ctx.log.Errorln("Error discovering RPC endpoints", err)
		return
	}

	ctx.log.Debugln("RPC endpoints:", rpcEndpoints)

	ctx.nodeLock.Lock()
	defer ctx.nodeLock.Unlock()

	for _, rpcEndpoint := range rpcEndpoints {
		if _, exists := ctx.secondaryNodes[rpcEndpoint]; !exists {
			ctx.log.Infoln("Added new RPC endpoint:", rpcEndpoint)
			rpcNodeHandler := NewRPCNodeHandler(rpcEndpoint)
			ctx.secondaryNodes[rpcEndpoint] = rpcNodeHandler
		}
	}
}
