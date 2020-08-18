//
// Copyright 2020 Wireline, Inc.
//

package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	wnsUtils "github.com/wirelineio/wns/utils"
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
	auctionTxCmd.AddCommand(client.PostCommands(
		GetCmdCommitBid(cdc),
	)...)

	return auctionTxCmd
}

// GetCmdCommitBid is the CLI command for committing a bid.
func GetCmdCommitBid(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-bid [auction-id] [bid-amount] [auction-fee]",
		Short: "Commit sealed bid.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Validate bid amount.
			bidAmount, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			mnemonic, err := wnsUtils.GenerateMnemonic()
			if err != nil {
				return err
			}

			chainID := viper.GetString("chain-id")
			auctionID := args[0]

			reveal := map[string]interface{}{
				"chainId":   chainID,
				"auctionId": auctionID,
				"bidAmount": bidAmount.String(),
				"salt":      mnemonic,
			}

			auctionFee, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			commitHash, content, err := wnsUtils.GenerateHash(reveal)
			if err != nil {
				return err
			}

			// Save reveal file.
			ioutil.WriteFile(fmt.Sprintf("%s.json", commitHash), []byte(content), 0600)

			msg := types.NewMsgCommitBid(auctionID, commitHash, auctionFee, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
