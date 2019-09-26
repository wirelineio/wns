//
// Copyright 2019 Wireline, Inc.
//

package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	nameserviceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the nameservice module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	nameserviceQueryCmd.AddCommand(client.GetCommands(
		GetCmdList(storeKey, cdc),
		GetCmdGetResource(storeKey, cdc),
	)...)
	return nameserviceQueryCmd
}

// GetCmdList queries all records.
func GetCmdList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List records.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.Set("trust-node", true)

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/list", queryRoute), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

// GetCmdGetResource queries a record record.
func GetCmdGetResource(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [ID]",
		Short: "Get record.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.Set("trust-node", true)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			id := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get/%s", queryRoute, id), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
