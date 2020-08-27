//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bond "github.com/wirelineio/wns/x/bond"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetRecord defines a SetResource message.
type MsgSetRecord struct {
	Payload PayloadObj     `json:"payload"`
	BondID  bond.ID        `json:"bondId"`
	Signer  sdk.AccAddress `json:"signer"`
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
	ID     ID             `json:"id"`
	Signer sdk.AccAddress `json:"signer"`
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
