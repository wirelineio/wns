//
// Copyright 2019 Wireline, Inc.
//

package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/wirelineio/wns/x/nameservice/internal/helpers"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSetRecord:
			return handleMsgSetRecord(ctx, keeper, msg)
		case types.MsgClearRecords:
			return handleMsgClearRecords(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetRecord.
func handleMsgSetRecord(ctx sdk.Context, keeper Keeper, msg types.MsgSetRecord) sdk.Result {
	payload := msg.Payload.ToPayload()
	record := types.Record{Attributes: payload.Record}

	// Check signatures.
	resourceSignBytes, _ := record.GetSignBytes()
	record.ID = record.GetCID()

	if exists := keeper.HasRecord(ctx, record.ID); exists {
		// Immutable record already exists. No-op.
		return sdk.Result{}
	}

	record.Owners = []string{}
	for _, sig := range payload.Signatures {
		pubKey, err := cryptoAmino.PubKeyFromBytes(helpers.BytesFromBase64(sig.PubKey))
		if err != nil {
			fmt.Println("Error decoding pubKey from bytes: ", err)
			return sdk.ErrUnauthorized("Invalid public key.").Result()
		}

		sigOK := pubKey.VerifyBytes(resourceSignBytes, helpers.BytesFromBase64(sig.Signature))
		if !sigOK {
			fmt.Println("Signature mismatch: ", sig.PubKey)
			return sdk.ErrUnauthorized("Invalid signature.").Result()
		}

		record.Owners = append(record.Owners, helpers.GetAddressFromPubKey(pubKey))
	}

	keeper.PutRecord(ctx, record)

	return sdk.Result{}
}

// Handle MsgClearRecords.
func handleMsgClearRecords(ctx sdk.Context, keeper Keeper, msg types.MsgClearRecords) sdk.Result {
	keeper.ClearRecords(ctx)

	return sdk.Result{}
}

func checkAccess(owners []string, record types.Record, signatures []types.Signature) bool {
	addresses := []string{}

	// Check signatures.
	resourceSignBytes, _ := record.GetSignBytes()
	for _, sig := range signatures {
		pubKey, err := cryptoAmino.PubKeyFromBytes(helpers.BytesFromBase64(sig.PubKey))
		if err != nil {
			fmt.Println("Error decoding pubKey from bytes: ", err)
			return false
		}

		allow := pubKey.VerifyBytes(resourceSignBytes, helpers.BytesFromBase64(sig.Signature))
		if !allow {
			fmt.Println("Signature mismatch: ", sig.PubKey)

			return false
		}

		addresses = append(addresses, helpers.GetAddressFromPubKey(pubKey))
	}

	// Check one of the addresses matches the owner.
	matches := helpers.Intersection(addresses, owners)
	if len(matches) == 0 {
		return false
	}

	return true
}
