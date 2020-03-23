//
// Copyright 2020 Wireline, Inc.
//

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wirelineio/wns/cmd/wnsd-lite/gql"
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
		chainID, _ := cmd.Flags().GetString("chain-id")
		home, _ := cmd.Flags().GetString("home")
		height, _ := cmd.Flags().GetInt64("height")

		config := sync.Config{ChainID: chainID, Home: home}
		ctx := sync.NewContext(&config, height)

		sync.Init(ctx)
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

		ctx := sync.NewContext(&config, height)

		go gql.Server(ctx)

		sync.Start(ctx)
	},
}

func init() {
	initCmd.Flags().Int64("height", 1, "Initial height (corresponding to genesis.json, if present)")

	startCmd.Flags().StringP("node", "n", "tcp://localhost:26657", "Upstream WNS node RPC address")

	// TODO(ashwin): Remove this flag after we start saving height in db.
	startCmd.Flags().Int64("height", 1, "Height to start synchronizing at")

	// Add flags for GQL server.
	startCmd.Flags().Bool("gql-server", true, "Start GQL server.")
	startCmd.Flags().Bool("gql-playground", true, "Enable GQL playground.")
	startCmd.Flags().String("gql-port", "9475", "Port to use for the GQL server.")
}
