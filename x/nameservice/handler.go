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
			return handleMsgSetResource(ctx, keeper, msg)
		case types.MsgClearRecords:
			return handleMsgClearResources(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetRecord.
func handleMsgSetResource(ctx sdk.Context, keeper Keeper, msg types.MsgSetRecord) sdk.Result {
	payload := types.PayloadObjToPayload(msg.Payload)
	record := &payload.Record

	record.ID = types.ID(helpers.GenRecordHash(record))
	if exists := keeper.HasResource(ctx, record.ID); exists {
		return sdk.ErrUnauthorized("Record already exists.").Result()
	}

	record.Owners = []string{}
	for _, sig := range payload.Signatures {
		pubKey, err := cryptoAmino.PubKeyFromBytes(helpers.BytesFromBase64(sig.PubKey))
		if err != nil {
			fmt.Println("Error decoding pubKey from bytes: ", err)
			return sdk.ErrUnauthorized("Invalid public key.").Result()
		}

		record.Owners = append(record.Owners, helpers.GetAddressFromPubKey(pubKey))
	}

	keeper.PutResource(ctx, payload.Record)

	return sdk.Result{}
}

// Handle MsgClearRecords.
func handleMsgClearResources(ctx sdk.Context, keeper Keeper, msg types.MsgClearRecords) sdk.Result {
	keeper.ClearResources(ctx)

	return sdk.Result{}
}

func checkAccess(owners []string, record types.Record, signatures []types.Signature) bool {
	addresses := []string{}

	// Check signatures.
	resourceSignBytes := helpers.GenRecordHash(&record)
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
