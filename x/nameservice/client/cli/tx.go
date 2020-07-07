//
// Copyright 2019 Wireline, Inc.
//

package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/wirelineio/wns/x/nameservice/internal/helpers"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// GetTxCmd returns transaction commands for this module.
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	nameserviceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Nameservice transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nameserviceTxCmd.AddCommand(client.PostCommands(
		GetCmdSetResource(cdc),
		GetCmdReserveName(cdc),
		GetCmdClearResources(cdc),
		GetCmdAssociateBond(cdc),
		GetCmdDissociateBond(cdc),
		GetCmdDissociateRecords(cdc),
		GetCmdReassociateRecords(cdc),
		GetCmdRenewRecord(cdc),
	)...)

	return nameserviceTxCmd
}

// GetCmdSetResource is the CLI command for creating/updating a record.
func GetCmdSetResource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [payload file path] [bond-id]",
		Short: "Set record.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			payload, err := getPayloadFromFile(args[0])
			if err != nil {
				return err
			}

			signOnly := viper.GetBool("sign-only")
			if signOnly {
				return signResource(payload)
			}

			msg := types.NewMsgSetRecord(payload.ToPayloadObj(), args[1], cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Bool("sign-only", false, "Only sign the transaction payload.")

	return cmd
}

// GetCmdReserveName is the CLI command for reserving a name.
func GetCmdReserveName(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reserve-name [name]",
		Short: "Reserve name.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgReserveName(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdAssociateBond is the CLI command for associating a record with a bond.
func GetCmdAssociateBond(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "associate-bond [record-id] [bond-id]",
		Short: "Associate record with bond.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgAssociateBond(args[0], args[1], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdDissociateBond is the CLI command for dissociating a record from a bond.
func GetCmdDissociateBond(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dissociate-bond [record-id]",
		Short: "Dissociate record from (existing) bond.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgDissociateBond(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdDissociateRecords is the CLI command for dissociating all records from a bond.
func GetCmdDissociateRecords(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dissociate-records [bond-id]",
		Short: "Dissociate all records from bond.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgDissociateRecords(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdReassociateRecords is the CLI command for reassociating all records from old to new bond.
func GetCmdReassociateRecords(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reassociate-records [old-bond-id] [new-bond-id]",
		Short: "Reassociates all records from old to new bond.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgReassociateRecords(args[0], args[1], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdClearResources is the CLI command for clearing all records.
// NOTE: FOR LOCAL TESTING PURPOSES ONLY!
func GetCmdClearResources(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear records.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgClearRecords(cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdRenewRecord is the CLI command for renewing an expired record.
func GetCmdRenewRecord(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew-record [record-id]",
		Short: "Renew (expired) record.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgRenewRecord(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// Load payload object from YAML file.
func getPayloadFromFile(filePath string) (types.Payload, error) {
	var payload types.Payload

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return payload, err
	}

	err = yaml.Unmarshal(data, &payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

// Sign payload object.
func signResource(payload types.Payload) error {
	name := viper.GetString("from")

	cid, sigBytes, signedJSON, pubKey, err := requestSignature(payload.Record, name)
	if err != nil {
		return err
	}

	fmt.Println("CID       :", cid)
	fmt.Println("Address   :", helpers.GetAddressFromPubKey(pubKey))
	fmt.Println("PubKey    :", helpers.BytesToBase64(pubKey.Bytes()))
	fmt.Println("Signature :", helpers.BytesToBase64(sigBytes))
	fmt.Println("SigData   :", string(signedJSON))

	return nil
}

// requestSignature returns a cryptographic signature for an object.
func requestSignature(attributes map[string]interface{}, name string) (types.ID, []byte, []byte, crypto.PubKey, error) {
	keybase, err := keys.NewKeyBaseFromHomeFlag()
	if err != nil {
		return "", nil, nil, nil, err
	}

	passphrase, err := keys.GetPassphrase(name)
	if err != nil {
		return "", nil, nil, nil, err
	}

	record := types.Record{Attributes: attributes}
	signBytes, signedJSON := record.GetSignBytes()
	sigBytes, pubKey, err := keybase.Sign(name, passphrase, signBytes)
	if err != nil {
		return "", nil, nil, nil, err
	}

	return record.GetCID(), sigBytes, signedJSON, pubKey, nil
}
