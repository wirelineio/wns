//
// Copyright 2019 Wireline, Inc.
//

package nameservice

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/wirelineio/wns/x/bond"
	"github.com/wirelineio/wns/x/nameservice/internal/helpers"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSetRecord:
			return handleMsgSetRecord(ctx, keeper, msg)
		case types.MsgReserveName:
			return handleMsgReserveName(ctx, keeper, msg)
		case types.MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case types.MsgAssociateBond:
			return handleMsgAssociateBond(ctx, keeper, msg)
		case types.MsgDissociateBond:
			return handleMsgDissociateBond(ctx, keeper, msg)
		case types.MsgDissociateRecords:
			return handleMsgDissociateRecords(ctx, keeper, msg)
		case types.MsgReassociateRecords:
			return handleMsgReassociateRecords(ctx, keeper, msg)
		case types.MsgRenewRecord:
			return handleMsgRenewRecord(ctx, keeper, msg)
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
	record := types.Record{Attributes: payload.Record, BondID: msg.BondID}

	// Check signatures.
	resourceSignBytes, _ := record.GetSignBytes()
	record.ID = record.GetCID()

	if exists := keeper.HasRecord(ctx, record.ID); exists {
		// Immutable record already exists. No-op.
		return sdk.Result{
			Data:   []byte(record.ID),
			Events: ctx.EventManager().Events(),
		}
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

	sdkErr := processRecord(ctx, keeper, &record, false)
	if sdkErr != nil {
		return sdkErr.Result()
	}

	return sdk.Result{
		Data:   []byte(record.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgRenewRecord.
func handleMsgRenewRecord(ctx sdk.Context, keeper Keeper, msg types.MsgRenewRecord) sdk.Result {
	if !keeper.HasRecord(ctx, msg.ID) {
		return sdk.ErrInternal("Record not found.").Result()
	}

	// Check if renewal is required (i.e. expired record marked as deleted).
	record := keeper.GetRecord(ctx, msg.ID)
	if !record.Deleted || record.ExpiryTime.After(ctx.BlockTime()) {
		return sdk.ErrInternal("Renewal not required.").Result()
	}

	err := processRecord(ctx, keeper, &record, true)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(record.ID),
		Events: ctx.EventManager().Events(),
	}
}

func processRecord(ctx sdk.Context, keeper Keeper, record *types.Record, isRenewal bool) sdk.Error {
	// Check that the record has an associated bond.
	if !keeper.BondKeeper.HasBond(ctx, record.BondID) {
		return sdk.ErrUnauthorized("Bond not found.")
	}

	bondObj := keeper.BondKeeper.GetBond(ctx, record.BondID)
	coins, err := sdk.ParseCoins(keeper.RecordRent(ctx))
	if err != nil {
		return sdk.ErrInvalidCoins("Invalid record rent.")
	}

	rent, err := sdk.ConvertCoin(coins[0], bond.MicroWire)
	if err != nil {
		return sdk.ErrInvalidCoins("Invalid record rent.")
	}

	// Deduct rent from bond.
	updatedBalance, isNeg := bondObj.Balance.SafeSub(sdk.NewCoins(rent))
	if isNeg {
		// Check if bond has sufficient funds.
		return sdk.ErrInsufficientCoins("Insufficient funds.")
	}

	// Move funds from bond module to record rent module.
	err = keeper.BondKeeper.SupplyKeeper.SendCoinsFromModuleToModule(ctx, bond.ModuleName, bond.RecordRentModuleAccountName, sdk.NewCoins(rent))
	if err != nil {
		return sdk.ErrInternal("Error withdrawing rent.")
	}

	// Update bond balance.
	bondObj.Balance = updatedBalance
	keeper.BondKeeper.SaveBond(ctx, bondObj)

	record.CreateTime = ctx.BlockHeader().Time
	record.ExpiryTime = ctx.BlockHeader().Time.Add(keeper.RecordExpiryTime(ctx))
	record.Deleted = false

	keeper.PutRecord(ctx, *record)
	keeper.InsertRecordExpiryQueue(ctx, *record)

	// Renewal doesn't change the name and bond indexes.
	if !isRenewal {
		keeper.AddBondToRecordIndexEntry(ctx, record.BondID, record.ID)
	}

	return nil
}

// Handle MsgClearRecords.
func handleMsgClearRecords(ctx sdk.Context, keeper Keeper, msg types.MsgClearRecords) sdk.Result {
	// keeper.ClearRecords(ctx)

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

	// Required so that renewal is triggered (with new bond ID) for expired records.
	if record.Deleted {
		keeper.InsertRecordExpiryQueue(ctx, record)
	}

	return sdk.Result{
		Data:   []byte(record.ID),
		Events: ctx.EventManager().Events(),
	}
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

	return sdk.Result{
		Data:   []byte(record.ID),
		Events: ctx.EventManager().Events(),
	}
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
	records := keeper.RecordKeeper.QueryRecordsByBond(ctx, msg.BondID)
	for _, record := range records {
		// Clear bond ID.
		record.BondID = ""
		keeper.PutRecord(ctx, record)
		keeper.RemoveBondToRecordIndexEntry(ctx, msg.BondID, record.ID)
	}

	return sdk.Result{
		Data:   []byte(bond.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgReassociateRecords.
func handleMsgReassociateRecords(ctx sdk.Context, keeper Keeper, msg types.MsgReassociateRecords) sdk.Result {

	if !keeper.BondKeeper.HasBond(ctx, msg.OldBondID) {
		return sdk.ErrInternal("Old bond not found.").Result()
	}

	if !keeper.BondKeeper.HasBond(ctx, msg.NewBondID) {
		return sdk.ErrInternal("New bond not found.").Result()
	}

	// Only the bond owner can reassociate all records.
	oldBond := keeper.BondKeeper.GetBond(ctx, msg.OldBondID)
	if msg.Signer.String() != oldBond.Owner {
		return sdk.ErrUnauthorized("Old bond owner mismatch.").Result()
	}

	newBond := keeper.BondKeeper.GetBond(ctx, msg.NewBondID)
	if msg.Signer.String() != newBond.Owner {
		return sdk.ErrUnauthorized("New bond owner mismatch.").Result()
	}

	// Reassociate all records.
	records := keeper.RecordKeeper.QueryRecordsByBond(ctx, msg.OldBondID)
	for _, record := range records {
		// Switch bond ID.
		record.BondID = msg.NewBondID
		keeper.PutRecord(ctx, record)

		keeper.RemoveBondToRecordIndexEntry(ctx, msg.OldBondID, record.ID)
		keeper.AddBondToRecordIndexEntry(ctx, msg.NewBondID, record.ID)

		// Required so that renewal is triggered (with new bond ID) for expired records.
		if record.Deleted {
			keeper.InsertRecordExpiryQueue(ctx, record)
		}
	}

	return sdk.Result{
		Data:   []byte(newBond.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgReserveName.
func handleMsgReserveName(ctx sdk.Context, keeper Keeper, msg types.MsgReserveName) sdk.Result {
	wrn := fmt.Sprintf("wrn://%s", msg.Name)

	parsedWRN, err := url.Parse(wrn)
	if err != nil {
		return sdk.ErrInternal("Invalid name.").Result()
	}

	name := parsedWRN.Host
	if fmt.Sprintf("wrn://%s", name) != wrn {
		return sdk.ErrInternal("Invalid name.").Result()
	}

	if strings.Contains(name, ".") {
		return sdk.ErrInternal(("Invalid name (dot is currently not allowed in root authority names).")).Result()
	}

	// Check if name already reserved.
	if keeper.HasNameAuthority(ctx, name) {
		return sdk.ErrInternal("Name already exists.").Result()
	}

	// Reserve name with signer as owner.
	account := keeper.AccountKeeper.GetAccount(ctx, msg.Signer)
	if account == nil {
		return sdk.ErrUnknownAddress("Account not found.").Result()
	}

	keeper.SetNameAuthority(ctx, name, NameAuthority{
		OwnerPublicKey: helpers.BytesToBase64(account.GetPubKey().Bytes()),
		OwnerAddress:   msg.Signer.String(),
		Height:         ctx.BlockHeight(),
	})

	return sdk.Result{
		Data:   []byte(name),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgSetName.
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg types.MsgSetName) sdk.Result {
	parsedWRN, err := url.Parse(msg.WRN)
	if err != nil {
		return sdk.ErrInternal("Invalid WRN.").Result()
	}

	name := parsedWRN.Host
	wrn := fmt.Sprintf("wrn://%s%s", name, parsedWRN.RequestURI())
	if wrn != msg.WRN {
		return sdk.ErrInternal("Invalid WRN.").Result()
	}

	// Check authority record.
	if !keeper.HasNameAuthority(ctx, name) {
		return sdk.ErrInternal("Name authority not found.").Result()
	}

	authority := keeper.GetNameAuthority(ctx, name)
	if authority.OwnerAddress != msg.Signer.String() {
		return sdk.ErrUnauthorized("Access denied.").Result()
	}

	keeper.SetNameRecord(ctx, wrn, msg.ID)

	return sdk.Result{
		Data:   []byte(wrn),
		Events: ctx.EventManager().Events(),
	}
}
