//
// Copyright 2020 Wireline, Inc.
//

package keeper

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/wirelineio/wns/x/bond"
	"github.com/wirelineio/wns/x/nameservice/internal/helpers"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// ProcessSetRecord creates a record.
func (k Keeper) ProcessSetRecord(ctx sdk.Context, msg types.MsgSetRecord) (*types.Record, sdk.Error) {
	payload := msg.Payload.ToPayload()
	record := types.Record{Attributes: payload.Record, BondID: msg.BondID}

	// Check signatures.
	resourceSignBytes, _ := record.GetSignBytes()
	cid, err := record.GetCID()
	if err != nil {
		return nil, sdk.ErrInternal("Invalid record JSON")
	}

	record.ID = cid

	if exists := k.HasRecord(ctx, record.ID); exists {
		// Immutable record already exists. No-op.
		return &record, nil
	}

	record.Owners = []string{}
	for _, sig := range payload.Signatures {
		pubKey, err := cryptoAmino.PubKeyFromBytes(helpers.BytesFromBase64(sig.PubKey))
		if err != nil {
			fmt.Println("Error decoding pubKey from bytes: ", err)
			return nil, sdk.ErrUnauthorized("Invalid public key.")
		}

		sigOK := pubKey.VerifyBytes(resourceSignBytes, helpers.BytesFromBase64(sig.Signature))
		if !sigOK {
			fmt.Println("Signature mismatch: ", sig.PubKey)
			return nil, sdk.ErrUnauthorized("Invalid signature.")
		}

		record.Owners = append(record.Owners, helpers.GetAddressFromPubKey(pubKey))
	}

	// Sort owners list.
	sort.Strings(record.Owners)

	sdkErr := k.processRecord(ctx, &record, false)
	if sdkErr != nil {
		return nil, sdkErr
	}

	return &record, nil
}

// ProcessRenewRecord renews a record.
func (k Keeper) ProcessRenewRecord(ctx sdk.Context, msg types.MsgRenewRecord) (*types.Record, sdk.Error) {
	if !k.HasRecord(ctx, msg.ID) {
		return nil, sdk.ErrInternal("Record not found.")
	}

	// Check if renewal is required (i.e. expired record marked as deleted).
	record := k.GetRecord(ctx, msg.ID)
	if !record.Deleted || record.ExpiryTime.After(ctx.BlockTime()) {
		return nil, sdk.ErrInternal("Renewal not required.")
	}

	err := k.processRecord(ctx, &record, true)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (k Keeper) processRecord(ctx sdk.Context, record *types.Record, isRenewal bool) sdk.Error {
	// Check that the record has an associated bond.
	if !k.BondKeeper.HasBond(ctx, record.BondID) {
		return sdk.ErrUnauthorized("Bond not found.")
	}

	bondObj := k.BondKeeper.GetBond(ctx, record.BondID)
	rent, err := sdk.ParseCoins(k.RecordRent(ctx))
	if err != nil {
		return sdk.ErrInvalidCoins("Invalid record rent.")
	}

	// Deduct rent from bond.
	updatedBalance, isNeg := bondObj.Balance.SafeSub(rent)
	if isNeg {
		// Check if bond has sufficient funds.
		return sdk.ErrInsufficientCoins("Insufficient funds.")
	}

	// Move funds from bond module to record rent module.
	err = k.SupplyKeeper.SendCoinsFromModuleToModule(ctx, bond.ModuleName, types.RecordRentModuleAccountName, rent)
	if err != nil {
		return sdk.ErrInternal("Error withdrawing rent.")
	}

	// Update bond balance.
	bondObj.Balance = updatedBalance
	k.BondKeeper.SaveBond(ctx, bondObj)

	record.CreateTime = ctx.BlockHeader().Time
	record.ExpiryTime = ctx.BlockHeader().Time.Add(k.RecordExpiryTime(ctx))
	record.Deleted = false

	k.PutRecord(ctx, *record)
	k.InsertRecordExpiryQueue(ctx, *record)

	// Renewal doesn't change the name and bond indexes.
	if !isRenewal {
		k.AddBondToRecordIndexEntry(ctx, record.BondID, record.ID)
	}

	return nil
}

// ProcessAssociateBond associates a record with a bond.
func (k Keeper) ProcessAssociateBond(ctx sdk.Context, msg types.MsgAssociateBond) (*types.Record, sdk.Error) {

	if !k.HasRecord(ctx, msg.ID) {
		return nil, sdk.ErrInternal("Record not found.")
	}

	if !k.BondKeeper.HasBond(ctx, msg.BondID) {
		return nil, sdk.ErrInternal("Bond not found.")
	}

	// Check if already associated with a bond.
	record := k.GetRecord(ctx, msg.ID)
	if record.BondID != "" {
		return nil, sdk.ErrUnauthorized("Bond already exists.")
	}

	// Only the bond owner can associate a record with the bond.
	bond := k.BondKeeper.GetBond(ctx, msg.BondID)
	if msg.Signer.String() != bond.Owner {
		return nil, sdk.ErrUnauthorized("Bond owner mismatch.")
	}

	record.BondID = msg.BondID
	k.PutRecord(ctx, record)
	k.AddBondToRecordIndexEntry(ctx, msg.BondID, msg.ID)

	// Required so that renewal is triggered (with new bond ID) for expired records.
	if record.Deleted {
		k.InsertRecordExpiryQueue(ctx, record)
	}

	return &record, nil
}

// ProcessDissociateBond dissociates a record from its bond.
func (k Keeper) ProcessDissociateBond(ctx sdk.Context, msg types.MsgDissociateBond) (*types.Record, sdk.Error) {

	if !k.HasRecord(ctx, msg.ID) {
		return nil, sdk.ErrInternal("Record not found.")
	}

	// Check if associated with a bond.
	record := k.GetRecord(ctx, msg.ID)
	bondID := record.BondID
	if bondID == "" {
		return nil, sdk.ErrUnauthorized("Bond not found.")
	}

	// Only the bond owner can dissociate a record from the bond.
	bond := k.BondKeeper.GetBond(ctx, bondID)
	if msg.Signer.String() != bond.Owner {
		return nil, sdk.ErrUnauthorized("Bond owner mismatch.")
	}

	// Clear bond ID.
	record.BondID = ""
	k.PutRecord(ctx, record)
	k.RemoveBondToRecordIndexEntry(ctx, bondID, record.ID)

	return &record, nil
}

// ProcessDissociateRecords dissociates all records associated with a given bond.
func (k Keeper) ProcessDissociateRecords(ctx sdk.Context, msg types.MsgDissociateRecords) (*bond.Bond, sdk.Error) {

	if !k.BondKeeper.HasBond(ctx, msg.BondID) {
		return nil, sdk.ErrInternal("Bond not found.")
	}

	// Only the bond owner can dissociate all records from the bond.
	bond := k.BondKeeper.GetBond(ctx, msg.BondID)
	if msg.Signer.String() != bond.Owner {
		return nil, sdk.ErrUnauthorized("Bond owner mismatch.")
	}

	// Dissociate all records from the bond.
	records := k.RecordKeeper.QueryRecordsByBond(ctx, msg.BondID)
	for _, record := range records {
		// Clear bond ID.
		record.BondID = ""
		k.PutRecord(ctx, record)
		k.RemoveBondToRecordIndexEntry(ctx, msg.BondID, record.ID)
	}

	return &bond, nil
}

// ProcessReassociateRecords switches records from and old to new bond.
func (k Keeper) ProcessReassociateRecords(ctx sdk.Context, msg types.MsgReassociateRecords) (*bond.Bond, sdk.Error) {

	if !k.BondKeeper.HasBond(ctx, msg.OldBondID) {
		return nil, sdk.ErrInternal("Old bond not found.")
	}

	if !k.BondKeeper.HasBond(ctx, msg.NewBondID) {
		return nil, sdk.ErrInternal("New bond not found.")
	}

	// Only the bond owner can reassociate all records.
	oldBond := k.BondKeeper.GetBond(ctx, msg.OldBondID)
	if msg.Signer.String() != oldBond.Owner {
		return nil, sdk.ErrUnauthorized("Old bond owner mismatch.")
	}

	newBond := k.BondKeeper.GetBond(ctx, msg.NewBondID)
	if msg.Signer.String() != newBond.Owner {
		return nil, sdk.ErrUnauthorized("New bond owner mismatch.")
	}

	// Reassociate all records.
	records := k.RecordKeeper.QueryRecordsByBond(ctx, msg.OldBondID)
	for _, record := range records {
		// Switch bond ID.
		record.BondID = msg.NewBondID
		k.PutRecord(ctx, record)

		k.RemoveBondToRecordIndexEntry(ctx, msg.OldBondID, record.ID)
		k.AddBondToRecordIndexEntry(ctx, msg.NewBondID, record.ID)

		// Required so that renewal is triggered (with new bond ID) for expired records.
		if record.Deleted {
			k.InsertRecordExpiryQueue(ctx, record)
		}
	}

	return &newBond, nil
}
