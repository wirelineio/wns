//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	nameservice "github.com/wirelineio/wns/x/nameservice"
)

// AggressiveSyncIntervalInMillis is the interval for aggressive sync, to catch up quickly to the current height.
const AggressiveSyncIntervalInMillis = 250

// SyncIntervalInMillis is the interval for initiating incremental sync, when already caught up to current height.
const SyncIntervalInMillis = 5 * 1000

// ErrorWaitDurationMillis is the wait duration in case of errors.
const ErrorWaitDurationMillis = 5 * 1000

// Init sets up the lite node.
func Init(ctx *Context, height int64) {
	// If sync record exists, abort with error.
	if ctx.Keeper.HasStatusRecord() {
		logErrorAndExit(errors.New("node already initialized, aborting"), 1)
	}

	// TODO(ashwin): Create <home>/config and <home>data directories.
	// TODO(ashwin): Create db in data directory.

	// Import genesis.json, if present.
	genesisJSONPath := path.Join(ctx.Config.Home, "config", "genesis.json")
	if _, err := os.Stat(genesisJSONPath); err == nil {
		geneisState := GenesisState{}
		bytes, err := ioutil.ReadFile(genesisJSONPath)
		if err != nil {
			logErrorAndExit(err, 1)
		}

		err = ctx.Codec.UnmarshalJSON(bytes, &geneisState)
		if err != nil {
			logErrorAndExit(err, 1)
		}

		names := geneisState.AppState.Nameservice.Names
		for _, nameEntry := range names {
			ctx.Keeper.SetNameRecord(nameEntry.Name, nameEntry.Entry)
		}

		records := geneisState.AppState.Nameservice.Records
		for _, record := range records {
			ctx.Keeper.PutRecord(record)
		}
	}

	// Create sync status record.
	ctx.Keeper.SaveStatus(Status{LastSyncedHeight: height})
}

// Start initiates the sync process.
func Start(ctx *Context) {
	// Fail if node has no sync status record.
	if !ctx.Keeper.HasStatusRecord() {
		logErrorAndExit(errors.New("node not initialized, aborting"), 1)
	}

	syncStatus := ctx.Keeper.GetStatusRecord()
	lastSyncedHeight := syncStatus.LastSyncedHeight

	for {
		chainCurrentHeight, err := ctx.getCurrentHeight()
		if err != nil {
			logErrorAndWait(err)
			continue
		}

		if lastSyncedHeight > chainCurrentHeight {
			panic("Last synced height cannot be greater than current chain height")
		}

		err = ctx.syncAtHeight(lastSyncedHeight)
		if err != nil {
			logErrorAndWait(err)
			continue
		}

		// Saved last synced height in db.
		lastSyncedHeight = lastSyncedHeight + 1
		ctx.Keeper.SaveStatus(Status{LastSyncedHeight: lastSyncedHeight})

		waitAfterSync(chainCurrentHeight, lastSyncedHeight)
	}
}

// syncAtHeight runs a sync cycle for the given height.
func (ctx *Context) syncAtHeight(height int64) error {
	fmt.Println("Syncing at height", height, time.Now().UTC())

	changeset, err := ctx.getBlockChangeset(height)
	if err != nil {
		return err
	}

	if changeset.Height <= 0 {
		// No changeset for this block, ignore.
		return nil
	}

	// Sync records.
	err = ctx.syncRecords(height, changeset.Records)
	if err != nil {
		return err
	}

	// Sync name records.
	err = ctx.syncNameRecords(height, changeset.Names)
	if err != nil {
		return err
	}

	// Flush cache changes to underlying store.
	ctx.Store.Write()

	return nil
}

func (ctx *Context) syncRecords(height int64, records []nameservice.ID) error {
	for _, id := range records {
		recordKey := nameservice.GetRecordIndexKey(id)
		value, err := ctx.getStoreValue(recordKey, height)
		if err != nil {
			return err
		}

		ctx.Store.Set(recordKey, value)
	}

	return nil
}

func (ctx *Context) syncNameRecords(height int64, names []string) error {
	for _, name := range names {
		nameRecordKey := nameservice.GetNameRecordIndexKey(name)
		value, err := ctx.getStoreValue(nameRecordKey, height)
		if err != nil {
			return err
		}

		ctx.Store.Set(nameRecordKey, value)
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

func logErrorAndWait(err error) {
	fmt.Println("Error:", err)

	// TODO(ashwin): Exponential backoff logic.
	time.Sleep(ErrorWaitDurationMillis * time.Millisecond)
}

func logErrorAndExit(err error, exitStatus int) {
	fmt.Println("Error:", err)
	os.Exit(exitStatus)
}
