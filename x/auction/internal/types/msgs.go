//
// Copyright 2020 Wireline, Inc.
//

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgCreateAuction defines a create auction message.
type MsgCreateAuction struct {
	Coins  sdk.Coins      `json:"coins"`
	Signer sdk.AccAddress `json:"signer"`
}

// NewMsgCreateAuction is the constructor function for MsgCreateAuction.
func NewMsgCreateAuction(denom string, amount int64, signer sdk.AccAddress) MsgCreateAuction {
	return MsgCreateAuction{
		Coins:  sdk.NewCoins(sdk.NewInt64Coin(denom, amount)),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgCreateAuction) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateAuction) Type() string { return "create" }

// ValidateBasic Implements Msg.
func (msg MsgCreateAuction) ValidateBasic() sdk.Error {

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	if len(msg.Coins) == 0 || !msg.Coins.IsValid() {
		return sdk.ErrInvalidCoins("Invalid amount.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCreateAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgRefillAuction defines a refill auction message.
type MsgRefillAuction struct {
	ID     ID             `json:"id"`
	Coins  sdk.Coins      `json:"coins"`
	Signer sdk.AccAddress `json:"signer"`
}

// NewMsgRefillAuction is the constructor function for MsgRefillAuction.
func NewMsgRefillAuction(id string, denom string, amount int64, signer sdk.AccAddress) MsgRefillAuction {
	return MsgRefillAuction{
		ID:     ID(id),
		Coins:  sdk.NewCoins(sdk.NewInt64Coin(denom, amount)),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgRefillAuction) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRefillAuction) Type() string { return "refill" }

// ValidateBasic Implements Msg.
func (msg MsgRefillAuction) ValidateBasic() sdk.Error {

	if string(msg.ID) == "" {
		return sdk.ErrInternal("Invalid auction ID.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	if len(msg.Coins) == 0 || !msg.Coins.IsValid() {
		return sdk.ErrInvalidCoins("Invalid amount.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRefillAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRefillAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgWithdrawAuction defines a withdraw (funds from) auction message.
type MsgWithdrawAuction struct {
	ID     ID             `json:"id"`
	Coins  sdk.Coins      `json:"coins"`
	Signer sdk.AccAddress `json:"signer"`
}

// NewMsgWithdrawAuction is the constructor function for MsgWithdrawAuction.
func NewMsgWithdrawAuction(id string, denom string, amount int64, signer sdk.AccAddress) MsgWithdrawAuction {
	return MsgWithdrawAuction{
		ID:     ID(id),
		Coins:  sdk.NewCoins(sdk.NewInt64Coin(denom, amount)),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgWithdrawAuction) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgWithdrawAuction) Type() string { return "withdraw" }

// ValidateBasic Implements Msg.
func (msg MsgWithdrawAuction) ValidateBasic() sdk.Error {

	if string(msg.ID) == "" {
		return sdk.ErrInternal("Invalid auction ID.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	if len(msg.Coins) == 0 || !msg.Coins.IsValid() {
		return sdk.ErrInvalidCoins("Invalid amount.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgWithdrawAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgWithdrawAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgCancelAuction defines a cancel auction message.
type MsgCancelAuction struct {
	ID     ID             `json:"id"`
	Signer sdk.AccAddress `json:"signer"`
}

// NewMsgCancelAuction is the constructor function for MsgCancelAuction.
func NewMsgCancelAuction(id string, signer sdk.AccAddress) MsgCancelAuction {
	return MsgCancelAuction{
		ID:     ID(id),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgCancelAuction) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCancelAuction) Type() string { return "cancel" }

// ValidateBasic Implements Msg.
func (msg MsgCancelAuction) ValidateBasic() sdk.Error {

	if string(msg.ID) == "" {
		return sdk.ErrInternal("Invalid auction ID.")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCancelAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCancelAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
