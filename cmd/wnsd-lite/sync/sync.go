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

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/cachekv"
	"github.com/cosmos/cosmos-sdk/store/dbadapter"
	"github.com/tendermint/go-amino"
	tmlite "github.com/tendermint/tendermint/lite"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	dbm "github.com/tendermint/tm-db"
	app "github.com/wirelineio/wns"
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
	NodeAddress string
	ChainID     string
	Home        string
}

// Context contains sync context info.
type Context struct {
	Config   *Config
	Codec    *amino.Codec
	Client   *rpcclient.HTTP
	Verifier tmlite.Verifier
	DBStore  store.KVStore
	Store    *cachekv.Store
	Keeper   *Keeper
}

// AppState is used to import initial app state (records, names) into the db.
type AppState struct {
	Nameservice nameservice.GenesisState `json:"nameservice" yaml:"nameservice"`
}

// GenesisState is used to import initial state into the db.
type GenesisState struct {
	AppState AppState `json:"app_state" yaml:"app_state"`
}

// NewContext creates a context object.
func NewContext(config *Config) *Context {

	// TODO(ashwin): Switch from in-mem store to persistent leveldb store.
	var mem store.KVStore = dbadapter.Store{DB: dbm.NewMemDB()}
	store := cachekv.NewStore(mem)

	codec := app.MakeCodec()

	nodeAddress := config.NodeAddress

	ctx := Context{
		Config:  config,
		Codec:   codec,
		DBStore: mem,
		Store:   store,
		Keeper:  NewKeeper(codec, mem),
	}

	if nodeAddress != "" {
		ctx.Client = rpcclient.NewHTTP(nodeAddress, "/websocket")
		ctx.Verifier = CreateVerifier(config)
	}

	return &ctx
}

// GetCurrentHeight gets the current WNS block height.
func (ctx *Context) GetCurrentHeight() (int64, error) {
	status, err := ctx.Client.Status()
	if err != nil {
		return 0, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

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
		chainCurrentHeight, err := ctx.GetCurrentHeight()
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

func (ctx *Context) getBlockChangeset(height int64) (*nameservice.BlockChangeset, error) {
	value, err := ctx.getStoreValue(nameservice.GetBlockChangesetIndexKey(height), height)
	if err != nil {
		return nil, err
	}

	var changeset nameservice.BlockChangeset
	ctx.Codec.MustUnmarshalBinaryBare(value, &changeset)

	return &changeset, nil
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
