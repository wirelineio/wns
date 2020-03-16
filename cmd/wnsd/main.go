//
// Copyright 2019 Wireline, Inc.
//

package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	genaccscli "github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	baseApp "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	app "github.com/wirelineio/wns"
)

const pruningStrategyFlag = "pruning"
const haltHeightFlag = "halt-height"

const pruningStrategySyncable = "syncable"
const pruningStrategyNothing = "nothing"
const pruningStrategyEverything = "everything"

var invCheckPeriod uint

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "wnsd",
		Short:             "WNS App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome),
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
			genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// Add flags for GQL server.
	rootCmd.PersistentFlags().Bool("gql-server", false, "Start GQL server.")
	rootCmd.PersistentFlags().Bool("gql-playground", false, "Enable GQL playground.")
	rootCmd.PersistentFlags().String("gql-port", "9473", "Port to use for the GQL server.")

	// Invariant checking flag.
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, "inv-check-period", 0, "Assert registered invariants every N blocks.")

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "NS", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	opts := []func(*baseApp.BaseApp){}
	opts = append(opts, getPruningStrategyOption(logger))
	opts = append(opts, getHaltHeightOption(logger))

	return app.NewNameServiceApp(logger, db, invCheckPeriod, opts...)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		nsApp := app.NewNameServiceApp(logger, db, uint(1))
		err := nsApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	nsApp := app.NewNameServiceApp(logger, db, uint(1))

	return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

func getPruningStrategyOption(logger log.Logger) func(*baseApp.BaseApp) {
	pruningStrategy := viper.GetString(pruningStrategyFlag)
	logger.Info(fmt.Sprintf("Pruning strategy: %s", pruningStrategy))

	switch pruningStrategy {
	case pruningStrategySyncable:
		return baseApp.SetPruning(types.PruneSyncable)
	case pruningStrategyNothing:
		return baseApp.SetPruning(types.PruneNothing)
	case pruningStrategyEverything:
		return baseApp.SetPruning(types.PruneEverything)
	default:
		panic(fmt.Sprintf("Invalid pruning strategy: %s", pruningStrategy))
	}
}

func getHaltHeightOption(logger log.Logger) func(*baseApp.BaseApp) {
	haltHeight := viper.GetInt64(haltHeightFlag)
	logger.Info(fmt.Sprintf("Halt height: %d", haltHeight))

	return baseApp.SetHaltHeight(uint64(haltHeight))
}
