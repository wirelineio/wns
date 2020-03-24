//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
	if ctx.keeper.HasStatusRecord() {
		ctx.log.Fatalln("Node already initialized, aborting.")
	}

	// Create <home>/config directory if it doesn't exist.
	configDirPath := filepath.Join(ctx.config.Home, "config")
	os.Mkdir(configDirPath, 0755)

	// Import genesis.json, if present.
	genesisJSONPath := filepath.Join(configDirPath, "genesis.json")
	if _, err := os.Stat(genesisJSONPath); err == nil {
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
	}

	// Create sync status record.
	ctx.keeper.SaveStatus(Status{LastSyncedHeight: height})
}

// Start initiates the sync process.
func Start(ctx *Context) {
	// Fail if node has no sync status record.
	if !ctx.keeper.HasStatusRecord() {
		ctx.log.Fatalln("Node not initialized, aborting.")
	}

	syncStatus := ctx.keeper.GetStatusRecord()
	lastSyncedHeight := syncStatus.LastSyncedHeight

	for {
		chainCurrentHeight, err := ctx.getCurrentHeight()
		if err != nil {
			logErrorAndWait(ctx, err)
			continue
		}

		if lastSyncedHeight > chainCurrentHeight {
			ctx.log.Panicln("Last synced height cannot be greater than current chain height.")
		}

		newSyncHeight := lastSyncedHeight + 1
		err = ctx.syncAtHeight(newSyncHeight)
		if err != nil {
			logErrorAndWait(ctx, err)
			continue
		}

		// Saved last synced height in db.
		lastSyncedHeight = newSyncHeight
		ctx.keeper.SaveStatus(Status{LastSyncedHeight: lastSyncedHeight})

		waitAfterSync(chainCurrentHeight, lastSyncedHeight)
	}
}

// syncAtHeight runs a sync cycle for the given height.
func (ctx *Context) syncAtHeight(height int64) error {
	ctx.log.Infoln("Syncing at height:", height, time.Now().UTC())

	changeset, err := ctx.getBlockChangeset(height)
	if err != nil {
		return err
	}

	if changeset.Height <= 0 {
		// No changeset for this block, ignore.
		return nil
	}

	ctx.log.Debugln("Syncing changeset:", changeset)

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
	ctx.cache.Write()

	return nil
}

func (ctx *Context) syncRecords(height int64, records []nameservice.ID) error {
	for _, id := range records {
		recordKey := nameservice.GetRecordIndexKey(id)
		value, err := ctx.getStoreValue(recordKey, height)
		if err != nil {
			return err
		}

		ctx.cache.Set(recordKey, value)
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
