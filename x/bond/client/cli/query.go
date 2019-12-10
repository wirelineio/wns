//
// Copyright 2019 Wireline, Inc.
//

package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/wirelineio/wns/x/bond/internal/types"
)

// GetQueryCmd returns query commands.
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	bondQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the bond module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	bondQueryCmd.AddCommand(client.GetCommands()...)
	return bondQueryCmd
}
