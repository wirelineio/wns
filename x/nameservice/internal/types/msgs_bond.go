//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bond "github.com/wirelineio/wns/x/bond"
)

// MsgAssociateBond defines a associate bond message.
type MsgAssociateBond struct {
	ID     ID
	BondID bond.ID
	Signer sdk.AccAddress
}

// NewMsgAssociateBond is the constructor function for MsgAssociateBond.
func NewMsgAssociateBond(id string, bondID string, signer sdk.AccAddress) MsgAssociateBond {
	return MsgAssociateBond{
		ID:     ID(id),
		BondID: bond.ID(bondID),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgAssociateBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgAssociateBond) Type() string { return "associate-bond" }

// ValidateBasic Implements Msg.
func (msg MsgAssociateBond) ValidateBasic() sdk.Error {

	if msg.ID == "" {
		return sdk.ErrInternal("Record ID is required.")
	}

	if msg.BondID == "" {
		return sdk.ErrInternal("Bond ID is required.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAssociateBond) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgAssociateBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgDissociateBond defines a dissociate bond message.
type MsgDissociateBond struct {
	ID     ID
	Signer sdk.AccAddress
}

// NewMsgDissociateBond is the constructor function for MsgDissociateBond.
func NewMsgDissociateBond(id string, signer sdk.AccAddress) MsgDissociateBond {
	return MsgDissociateBond{
		ID:     ID(id),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgDissociateBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDissociateBond) Type() string { return "dissociate-bond" }

// ValidateBasic Implements Msg.
func (msg MsgDissociateBond) ValidateBasic() sdk.Error {

	if msg.ID == "" {
		return sdk.ErrInternal("Record ID is required.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDissociateBond) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDissociateBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
