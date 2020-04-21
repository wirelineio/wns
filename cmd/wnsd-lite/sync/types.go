//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"path/filepath"

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
	ChainID  string   `json:"chain_id" yaml:"chain_id"`
	AppState AppState `json:"app_state" yaml:"app_state"`
}

// Config represents config for sync functionality.
type Config struct {
	LogLevel            string
	NodeAddress         string
	ChainID             string
	Home                string
	InitFromNode        bool
	InitFromGenesisFile bool
	Endpoint            string
}

// Context contains sync context info.
type Context struct {
	config *Config
	codec  *amino.Codec

	// Primary RPC primaryNode, used for verification.
	primaryNode *RPCNode

	// Other RPC secondaryNodes, used for load distribution.
	secondaryNodes map[string]*RPCNode

	log      *logrus.Logger
	verifier tmlite.Verifier
	store    store.KVStore
	cache    *cachekv.Store
	keeper   *Keeper
}

// RPCNode is used to call an RPC endpoint and maintains basic stats.
type RPCNode struct {
	Address string          `json:"address"`
	Client  *rpcclient.HTTP `json:"-"`
	Calls   int64           `json:"calls"`
	Errors  int64           `json:"errors"`
}

// NewContext creates a context object.
func NewContext(config *Config) *Context {
	log := logrus.New()

	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatalln(err)
	}

	log.SetLevel(logLevel)

	db := dbm.NewDB("graph", dbm.GoLevelDBBackend, filepath.Join(config.Home, "data"))
	var dbStore store.KVStore = dbadapter.Store{DB: db}
	cacheStore := cachekv.NewStore(dbStore)

	codec := app.MakeCodec()

	nodeAddress := config.NodeAddress

	ctx := Context{
		config:         config,
		codec:          codec,
		store:          dbStore,
		cache:          cacheStore,
		log:            log,
		secondaryNodes: make(map[string]*RPCNode),
	}

	ctx.keeper = NewKeeper(&ctx)

	if nodeAddress != "" {
		ctx.primaryNode = &RPCNode{
			Client:  rpcclient.NewHTTP(nodeAddress, "/websocket"),
			Address: nodeAddress,
			Calls:   0,
			Errors:  0,
		}

		// Init secondary nodes, as they should have at least one entry.
		// Don't assume endpoint will be passed for discovery of secondary nodes.
		ctx.secondaryNodes[nodeAddress] = ctx.primaryNode

		ctx.verifier = CreateVerifier(config)
	}

	return &ctx
}
