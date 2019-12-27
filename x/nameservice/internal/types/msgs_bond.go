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
	ID     ID             `json:"id"`
	BondID bond.ID        `json:"bondId"`
	Signer sdk.AccAddress `json:"signer"`
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
	ID     ID             `json:"id"`
	Signer sdk.AccAddress `json:"signer"`
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

// MsgDissociateRecords defines a dissociate all records from bond message.
type MsgDissociateRecords struct {
	BondID bond.ID        `json:"bondId"`
	Signer sdk.AccAddress `json:"signer"`
}

// NewMsgDissociateRecords is the constructor function for MsgDissociateRecords.
func NewMsgDissociateRecords(bondID string, signer sdk.AccAddress) MsgDissociateRecords {
	return MsgDissociateRecords{
		BondID: bond.ID(bondID),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgDissociateRecords) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDissociateRecords) Type() string { return "dissociate-records" }

// ValidateBasic Implements Msg.
func (msg MsgDissociateRecords) ValidateBasic() sdk.Error {

	if msg.BondID == "" {
		return sdk.ErrInternal("Bond ID is required.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDissociateRecords) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDissociateRecords) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgReassociateRecords defines a reassociate records message.
type MsgReassociateRecords struct {
	OldBondID bond.ID        `json:"oldBondId"`
	NewBondID bond.ID        `json:"newBondId"`
	Signer    sdk.AccAddress `json:"signer"`
}

// NewMsgReassociateRecords is the constructor function for MsgReassociateRecords.
func NewMsgReassociateRecords(oldBondID string, newBondID string, signer sdk.AccAddress) MsgReassociateRecords {
	return MsgReassociateRecords{
		OldBondID: bond.ID(oldBondID),
		NewBondID: bond.ID(newBondID),
		Signer:    signer,
	}
}

// Route Implements Msg.
func (msg MsgReassociateRecords) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgReassociateRecords) Type() string { return "reassociate-records" }

// ValidateBasic Implements Msg.
func (msg MsgReassociateRecords) ValidateBasic() sdk.Error {

	if msg.OldBondID == "" {
		return sdk.ErrInternal("Old Bond ID is required.")
	}

	if msg.NewBondID == "" {
		return sdk.ErrInternal("New Bond ID is required.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgReassociateRecords) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgReassociateRecords) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
