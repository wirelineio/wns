//
// Copyright 2020 Wireline, Inc.
//

package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/cachekv"
	"github.com/cosmos/cosmos-sdk/store/dbadapter"
	"github.com/spf13/cobra"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	dbm "github.com/tendermint/tm-db"
	app "github.com/wirelineio/wns"
	sync "github.com/wirelineio/wns/cmd/wnsd-lite/sync"
)

// Version => WNS Lite node version.
const Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the node version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the WNS lite node",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the WNS lite node",
	Run: func(cmd *cobra.Command, args []string) {
		nodeAddress, _ := cmd.Flags().GetString("node")
		height, _ := cmd.Flags().GetInt64("height")
		chainID, _ := cmd.Flags().GetString("chain-id")
		home, _ := cmd.Flags().GetString("home")

		config := sync.Config{
			NodeAddress: nodeAddress,
			ChainID:     chainID,
			Home:        home,
		}

		// TODO(ashwin): Switch from in-mem store to persistent leveldb store.
		mem := dbadapter.Store{DB: dbm.NewMemDB()}
		store := cachekv.NewStore(mem)

		ctx := sync.Context{
			Config:           &config,
			LastSyncedHeight: height,
			Client:           rpcclient.NewHTTP(nodeAddress, "/websocket"),
			Verifier:         sync.CreateVerifier(&config),
			Codec:            app.MakeCodec(),
			Store:            store,
		}

		sync.Start(&ctx)
	},
}

func init() {
	startCmd.Flags().StringP("node", "n", "tcp://localhost:26657", "Upstream WNS node RPC address")

	// TODO(ashwin): Remove this flag after we start saving height in db.
	startCmd.Flags().Int64("height", 1, "Height to start synchronizing at")
}
