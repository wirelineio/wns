//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tendermint/go-amino"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	nameservice "github.com/wirelineio/wns/x/nameservice"
)

// AggressiveSyncIntervalInMillis is the interval for aggressive sync, to catch up quickly to the current height.
const AggressiveSyncIntervalInMillis = 250

// SyncIntervalInMillis is the interval for initiating incremental sync, when already caught up to current height.
const SyncIntervalInMillis = 5 * 1000

// ErrorWaitDurationMillis is the wait duration in case of errors.
const ErrorWaitDurationMillis = 5 * 1000

// Config represents config for sync functionality.
type Config struct {
	NodeAddress      string
	Client           *rpcclient.HTTP
	Codec            *amino.Codec
	LastSyncedHeight int64
}

// GetCurrentHeight gets the current WNS block height.
func GetCurrentHeight(config *Config) (int64, error) {
	status, err := config.Client.Status()
	if err != nil {
		return 0, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

// Start initiates the sync process.
func Start(config *Config) {
	lastSyncedHeight := config.LastSyncedHeight

	for {
		chainCurrentHeight, err := GetCurrentHeight(config)
		if err != nil {
			logErrorAndWait(err)
			continue
		}

		if lastSyncedHeight > chainCurrentHeight {
			panic("Last synced height cannot be greater than current chain height")
		}

		err = syncAtHeight(config, lastSyncedHeight)
		if err != nil {
			logErrorAndWait(err)
			continue
		}

		// TODO(ashwin): Saved last synced height in db.
		lastSyncedHeight = lastSyncedHeight + 1

		waitAfterSync(chainCurrentHeight, lastSyncedHeight)
	}
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
	fmt.Println("Error", err)

	// TODO(ashwin): Exponential backoff logic.
	time.Sleep(ErrorWaitDurationMillis * time.Millisecond)
}

// syncAtHeight runs a sync cycle for the given height.
func syncAtHeight(config *Config, height int64) error {
	fmt.Println("Syncing at height", height, time.Now().UTC())

	cdc := config.Codec

	value, err := getStoreValue(config, nameservice.GetBlockChangesetIndexKey(height), height)
	if err != nil {
		return err
	}

	var changeset nameservice.BlockChangeset
	cdc.MustUnmarshalBinaryBare(value, &changeset)

	if changeset.Height <= 0 {
		// No changeset for this block, ignore.
		return nil
	}

	fmt.Println(string(cdc.MustMarshalJSON(changeset)))

	for _, id := range changeset.Records {
		value, err := getStoreValue(config, nameservice.GetRecordIndexKey(id), height)
		if err != nil {
			return err
		}

		var record nameservice.RecordObj
		cdc.MustUnmarshalBinaryBare(value, &record)

		jsonBytes, _ := json.MarshalIndent(record.ToRecord(), "", "  ")
		fmt.Println(string(jsonBytes))
	}

	for _, name := range changeset.Names {
		value, err := getStoreValue(config, nameservice.GetNameRecordIndexKey(name), height)
		if err != nil {
			return err
		}

		var nameRecord nameservice.NameRecord
		cdc.MustUnmarshalBinaryBare(value, &nameRecord)

		jsonBytes, _ := json.MarshalIndent(nameRecord, "", "  ")
		fmt.Println(name, string(jsonBytes))
	}

	return nil
}

func getStoreValue(config *Config, key []byte, height int64) ([]byte, error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: height,
		Prove:  true,
	}

	res, err := config.Client.ABCIQueryWithOptions("/store/nameservice/key", key, opts)
	if err != nil {
		return nil, err
	}

	// TODO(ashwin): Verify proof.

	return res.Response.Value, nil
}
