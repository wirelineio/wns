//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"path"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/cachekv"
	"github.com/cosmos/cosmos-sdk/store/dbadapter"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/go-amino"
	tmlite "github.com/tendermint/tendermint/lite"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	dbm "github.com/tendermint/tm-db"
	app "github.com/wirelineio/wns"
	"github.com/wirelineio/wns/x/nameservice"
)

// AppState is used to import initial app state (records, names) into the db.
type AppState struct {
	Nameservice nameservice.GenesisState `json:"nameservice" yaml:"nameservice"`
}

// GenesisState is used to import initial state into the db.
type GenesisState struct {
	AppState AppState `json:"app_state" yaml:"app_state"`
}

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
	Log      *logrus.Logger
}

// NewContext creates a context object.
func NewContext(config *Config) *Context {
	log := logrus.New()

	// TODO(ashwin): Configure using CLI flag.
	log.SetLevel(logrus.DebugLevel)

	db := dbm.NewDB("graph", dbm.GoLevelDBBackend, path.Join(config.Home, "data"))
	var dbStore store.KVStore = dbadapter.Store{DB: db}
	store := cachekv.NewStore(dbStore)

	codec := app.MakeCodec()

	nodeAddress := config.NodeAddress

	ctx := Context{
		Config:  config,
		Codec:   codec,
		DBStore: dbStore,
		Store:   store,
		Keeper:  NewKeeper(codec, dbStore),
		Log:     log,
	}

	if nodeAddress != "" {
		ctx.Client = rpcclient.NewHTTP(nodeAddress, "/websocket")
		ctx.Verifier = CreateVerifier(config)
	}

	return &ctx
}
