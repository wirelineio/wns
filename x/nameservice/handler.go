//
// Copyright 2019 Wireline, Inc.
//

package nameservice

import (
	"fmt"
	"reflect"
	"sort"

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
		case types.MsgAssociateBond:
			return handleMsgAssociateBond(ctx, keeper, msg)
		case types.MsgDissociateBond:
			return handleMsgDissociateBond(ctx, keeper, msg)
		case types.MsgDissociateRecords:
			return handleMsgDissociateRecords(ctx, keeper, msg)
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

	if exists := keeper.HasNameRecord(ctx, record.WRN()); exists {
		return sdk.ErrUnauthorized("Name record already exists.").Result()
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

	// Sort owners list.
	sort.Strings(record.Owners)

	// Basic access control - check if record owners === owners of `latest` (according to semver) version.
	if keeper.HasNameRecord(ctx, record.BaseWRN()) {
		latestNameRecord := keeper.GetNameRecord(ctx, record.BaseWRN())
		latestRecord := keeper.GetRecord(ctx, latestNameRecord.ID)
		if !reflect.DeepEqual(latestRecord.Owners, record.Owners) {
			return sdk.ErrUnauthorized("Owners mismatch, operation not allowed.").Result()
		}
	}

	keeper.PutRecord(ctx, record)
	processNameRecords(ctx, keeper, record)

	return sdk.Result{}
}

// Handle MsgClearRecords.
func handleMsgClearRecords(ctx sdk.Context, keeper Keeper, msg types.MsgClearRecords) sdk.Result {
	keeper.ClearRecords(ctx)

	return sdk.Result{}
}

// Handle MsgAssociateBond.
func handleMsgAssociateBond(ctx sdk.Context, keeper Keeper, msg types.MsgAssociateBond) sdk.Result {

	if !keeper.HasRecord(ctx, msg.ID) {
		return sdk.ErrInternal("Record not found.").Result()
	}

	if !keeper.BondKeeper.HasBond(ctx, msg.BondID) {
		return sdk.ErrInternal("Bond not found.").Result()
	}

	// Check if already associated with a bond.
	record := keeper.GetRecord(ctx, msg.ID)
	if record.BondID != "" {
		return sdk.ErrUnauthorized("Bond already exists.").Result()
	}

	// Only the bond owner can associate a record with the bond.
	bond := keeper.BondKeeper.GetBond(ctx, msg.BondID)
	if msg.Signer.String() != bond.Owner {
		return sdk.ErrUnauthorized("Bond owner mismatch.").Result()
	}

	record.BondID = msg.BondID
	keeper.PutRecord(ctx, record)
	keeper.AddBondToRecordIndexEntry(ctx, msg.BondID, msg.ID)

	return sdk.Result{}
}

// Handle MsgDissociateBond.
func handleMsgDissociateBond(ctx sdk.Context, keeper Keeper, msg types.MsgDissociateBond) sdk.Result {

	if !keeper.HasRecord(ctx, msg.ID) {
		return sdk.ErrInternal("Record not found.").Result()
	}

	// Check if associated with a bond.
	record := keeper.GetRecord(ctx, msg.ID)
	bondID := record.BondID
	if bondID == "" {
		return sdk.ErrUnauthorized("Bond not found.").Result()
	}

	// Only the bond owner can dissociate a record from the bond.
	bond := keeper.BondKeeper.GetBond(ctx, bondID)
	if msg.Signer.String() != bond.Owner {
		return sdk.ErrUnauthorized("Bond owner mismatch.").Result()
	}

	// Clear bond ID.
	record.BondID = ""
	keeper.PutRecord(ctx, record)
	keeper.RemoveBondToRecordIndexEntry(ctx, bondID, record.ID)

	return sdk.Result{}
}

// Handle MsgDissociateRecords.
func handleMsgDissociateRecords(ctx sdk.Context, keeper Keeper, msg types.MsgDissociateRecords) sdk.Result {

	if !keeper.BondKeeper.HasBond(ctx, msg.BondID) {
		return sdk.ErrInternal("Bond not found.").Result()
	}

	// Only the bond owner can dissociate all records from the bond.
	bond := keeper.BondKeeper.GetBond(ctx, msg.BondID)
	if msg.Signer.String() != bond.Owner {
		return sdk.ErrUnauthorized("Bond owner mismatch.").Result()
	}

	// Dissociate all records from the bond.
	records := keeper.QueryRecordsByBond(ctx, msg.BondID)
	for _, record := range records {
		// Clear bond ID.
		record.BondID = ""
		keeper.PutRecord(ctx, record)
		keeper.RemoveBondToRecordIndexEntry(ctx, msg.BondID, record.ID)
	}

	return sdk.Result{}
}

func processNameRecords(ctx sdk.Context, keeper Keeper, record types.Record) {
	keeper.SetNameRecord(ctx, record.WRN(), record.ToNameRecord())
	maybeUpdateBaseNameRecord(ctx, keeper, record)
}

func maybeUpdateBaseNameRecord(ctx sdk.Context, keeper Keeper, record types.Record) {
	if !keeper.HasNameRecord(ctx, record.BaseWRN()) {
		// Create base name record.
		keeper.SetNameRecord(ctx, record.BaseWRN(), record.ToNameRecord())
		return
	}

	// Get current base record (which will have current latest version).
	baseNameRecord := keeper.GetNameRecord(ctx, record.BaseWRN())
	latestRecord := keeper.GetRecord(ctx, baseNameRecord.ID)

	latestVersion := helpers.GetSemver(latestRecord.Version())
	createdVersion := helpers.GetSemver(record.Version())
	if createdVersion.GreaterThan(latestVersion) {
		// Need to update the base name record.
		keeper.SetNameRecord(ctx, record.BaseWRN(), record.ToNameRecord())
	}
}
