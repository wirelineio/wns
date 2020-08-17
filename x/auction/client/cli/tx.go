//
// Copyright 2020 Wireline, Inc.
//

package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/wirelineio/wns/x/auction/internal/types"
)

// GetTxCmd returns transaction commands for this module.
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	auctionTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Auction transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// TODO(ashwin): Add Tx commands.
	auctionTxCmd.AddCommand(client.PostCommands()...)

	return auctionTxCmd
}
