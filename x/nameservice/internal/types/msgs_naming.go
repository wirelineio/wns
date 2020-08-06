//
// Copyright 2020 Wireline, Inc.
//

package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// MsgReserveAuthority defines a ReserveName message.
type MsgReserveAuthority struct {
	Name   string         `json:"name"`
	Signer sdk.AccAddress `json:"signer"`
}

// NewMsgReserveAuthority is the constructor function for MsgReserveAuthority.
func NewMsgReserveAuthority(name string, signer sdk.AccAddress) MsgReserveAuthority {
	return MsgReserveAuthority{
		Name:   name,
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgReserveAuthority) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgReserveAuthority) Type() string { return "reserve-authority" }

// ValidateBasic Implements Msg.
func (msg MsgReserveAuthority) ValidateBasic() sdk.Error {

	if msg.Name == "" {
		return sdk.ErrInternal("Name is required.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgReserveAuthority) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgReserveAuthority) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgSetName defines a SetName message.
type MsgSetName struct {
	WRN    string         `json:"wrn"`
	ID     ID             `json:"id"`
	Signer sdk.AccAddress `json:"signer"`
}

// NewMsgSetName is the constructor function for MsgSetName.
func NewMsgSetName(wrn string, id string, signer sdk.AccAddress) MsgSetName {
	return MsgSetName{
		WRN:    wrn,
		ID:     ID(id),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgSetName) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetName) Type() string { return "set-name" }

// ValidateBasic Implements Msg.
func (msg MsgSetName) ValidateBasic() sdk.Error {

	if msg.WRN == "" {
		return sdk.ErrInternal("WRN is required.")
	}

	if msg.ID == "" {
		return sdk.ErrInternal("ID is required.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgDeleteName defines a DeleteName message.
type MsgDeleteName struct {
	WRN    string         `json:"wrn"`
	Signer sdk.AccAddress `json:"signer"`
}

// NewMsgDeleteName is the constructor function for MsgDeleteName.
func NewMsgDeleteName(wrn string, signer sdk.AccAddress) MsgDeleteName {
	return MsgDeleteName{
		WRN:    wrn,
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgDeleteName) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDeleteName) Type() string { return "delete-name" }

// ValidateBasic Implements Msg.
func (msg MsgDeleteName) ValidateBasic() sdk.Error {

	if msg.WRN == "" {
		return sdk.ErrInternal("WRN is required.")
	}
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDeleteName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
