//
// Copyright 2020 Wireline, Inc.
//

package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	app "github.com/wirelineio/wns"
	sync "github.com/wirelineio/wns/cmd/wnsd-lite/sync"
)

// Version => WNS Lite node version.
const Version = "0.1.0"

// SyncIntervalInSeconds is the ticker interval for initiating a sync cycle.
const SyncIntervalInSeconds = 1

// nodeAddress is the Tendermint RPC address of the upstream WNS node.
var nodeAddress string

// height to start sync at. To start at the last saved height, use -1 (the default).
var height int64

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
		config := sync.Config{
			NodeAddress: nodeAddress,
			Client:      rpcclient.NewHTTP(nodeAddress, "/websocket"),
			Codec:       app.MakeCodec(),
		}

		ticker := time.NewTicker(SyncIntervalInSeconds * time.Second)
		for now := range ticker.C {
			err := sync.Synchronize(config, height, now)
			if err != nil {
				fmt.Println("Error", err)

				// Continue waiting at this block height in case of errors.
				continue
			}

			// TODO(ashwin): Saved last synced height in db.
			height = height + 1
		}
	},
}

func init() {
	startCmd.Flags().StringVarP(&nodeAddress, "node", "n", "tcp://localhost:26657", "Upstream WNS node RPC address")
	startCmd.Flags().Int64Var(&height, "height", 1, "Height to start synchronizing at")
}
