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
		logLevel, _ := cmd.Flags().GetString("log-level")
		chainID, _ := cmd.Flags().GetString("chain-id")
		home, _ := cmd.Flags().GetString("home")
		nodeAddress, _ := cmd.Flags().GetString("node")
		height, _ := cmd.Flags().GetInt64("height")
		initFromNode, _ := cmd.Flags().GetBool("from-node")
		initFromGenesisFile, _ := cmd.Flags().GetBool("from-genesis-file")

		config := sync.Config{
			LogLevel:            logLevel,
			ChainID:             chainID,
			Home:                home,
			NodeAddress:         nodeAddress,
			InitFromNode:        initFromNode,
			InitFromGenesisFile: initFromGenesisFile,
		}
		ctx := sync.NewContext(&config)

		sync.Init(ctx, height)
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the WNS lite node",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel, _ := cmd.Flags().GetString("log-level")
		chainID, _ := cmd.Flags().GetString("chain-id")
		home, _ := cmd.Flags().GetString("home")
		nodeAddress, _ := cmd.Flags().GetString("node")

		config := sync.Config{
			LogLevel:    logLevel,
			ChainID:     chainID,
			Home:        home,
			NodeAddress: nodeAddress,
		}

		ctx := sync.NewContext(&config)

		go gql.Server(ctx)

		sync.Start(ctx)
	},
}

func init() {
	// Init command flags.
	initCmd.Flags().Bool("from-node", false, "Initialize from trusted node.")
	initCmd.Flags().Bool("from-genesis-file", false, "Initialize from genesis file.")
	initCmd.Flags().Int64("height", 1, "Initial height (if using --from-genesis-file option)")

	// Start command flags.
	startCmd.Flags().Bool("gql-server", true, "Start GQL server.")
	startCmd.Flags().Bool("gql-playground", true, "Enable GQL playground.")
	startCmd.Flags().String("gql-port", "9473", "Port to use for the GQL server.")
	startCmd.Flags().String("gql-playground-api-base", "", "GQL API base path to use in GQL playground.")
}
