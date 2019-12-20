//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	"strings"

	"github.com/Masterminds/semver"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bond "github.com/wirelineio/wns/x/bond"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetRecord defines a SetResource message.
type MsgSetRecord struct {
	Payload PayloadObj
	BondID  bond.ID
	Signer  sdk.AccAddress
}

// NewMsgSetRecord is the constructor function for MsgSetRecord.
func NewMsgSetRecord(payload PayloadObj, bondID string, signer sdk.AccAddress) MsgSetRecord {
	return MsgSetRecord{
		Payload: payload,
		BondID:  bond.ID(bondID),
		Signer:  signer,
	}
}

// Route Implements Msg.
func (msg MsgSetRecord) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetRecord) Type() string { return "set" }

// ValidateBasic Implements Msg.
func (msg MsgSetRecord) ValidateBasic() sdk.Error {

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	owners := msg.Payload.Record.Owners
	for _, owner := range owners {
		if owner == "" {
			return sdk.ErrInternal("Record owner not set.")
		}
	}

	record := msg.Payload.Record.ToRecord()
	wrnType := record.Type()
	if wrnType == "" {
		return sdk.ErrInternal("Record 'type' not set.")
	}

	if !strings.HasPrefix(wrnType, "wrn:") {
		return sdk.ErrInternal("Record 'type' is invalid.")
	}

	// TODO(ashwin): More validation checks.
	if record.Name() == "" {
		return sdk.ErrInternal("Record 'name' not set.")
	}

	version := record.Version()
	if version == "" {
		return sdk.ErrInternal("Record 'version' not set.")
	}

	_, err := semver.NewVersion(version)
	if err != nil {
		return sdk.ErrInternal("Record 'version' is invalid.")
	}

	if msg.BondID == "" {
		return sdk.ErrUnauthorized("Bond ID is required.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetRecord) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgRenewRecord defines a renew record message.
type MsgRenewRecord struct {
	ID     ID
	Signer sdk.AccAddress
}

// NewMsgRenewRecord is the constructor function for MsgRenewRecord.
func NewMsgRenewRecord(id string, signer sdk.AccAddress) MsgRenewRecord {
	return MsgRenewRecord{
		ID:     ID(id),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgRenewRecord) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRenewRecord) Type() string { return "set" }

// ValidateBasic Implements Msg.
func (msg MsgRenewRecord) ValidateBasic() sdk.Error {

	if msg.ID == "" {
		return sdk.ErrInternal("Record ID is required.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRenewRecord) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRenewRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgClearRecords defines a MsgClearRecords message.
type MsgClearRecords struct {
	Signer sdk.AccAddress
}

// NewMsgClearRecords is the constructor function for MsgClearRecords.
func NewMsgClearRecords(signer sdk.AccAddress) MsgClearRecords {
	return MsgClearRecords{
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgClearRecords) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgClearRecords) Type() string { return "clear" }

// ValidateBasic Implements Msg.
func (msg MsgClearRecords) ValidateBasic() sdk.Error {

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgClearRecords) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgClearRecords) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
