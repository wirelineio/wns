//
// Copyright 2020 Wireline, Inc.
//

package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

// DefaultLightNodeHome is the root directory for the wnsd-lite node.
var DefaultLightNodeHome = os.ExpandEnv("$HOME/.wireline/wnsd-lite")

func main() {
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "wnsd-lite",
		Short: "WNS Lite",
	}

	rootCmd.PersistentFlags().String("chain-id", "wireline", "Chain identifier")

	rootCmd.AddCommand(versionCmd, initCmd, startCmd)

	executor := cli.PrepareBaseCmd(rootCmd, "NSL", DefaultLightNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
