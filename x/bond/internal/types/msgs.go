//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgCreateBond defines a create bond message.
type MsgCreateBond struct {
	Coins  sdk.Coins
	Signer sdk.AccAddress
}

// NewMsgCreateBond is the constructor function for MsgCreateBond.
func NewMsgCreateBond(denom string, amount int64, signer sdk.AccAddress) MsgCreateBond {
	return MsgCreateBond{
		Coins:  sdk.NewCoins(sdk.NewInt64Coin(denom, amount)),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgCreateBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateBond) Type() string { return "create" }

// ValidateBasic Implements Msg.
func (msg MsgCreateBond) ValidateBasic() sdk.Error {

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	if len(msg.Coins) == 0 || !msg.Coins.IsValid() {
		return sdk.ErrInvalidCoins("Invalid amount.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateBond) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCreateBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgRefillBond defines a refill bond message.
type MsgRefillBond struct {
	ID     ID
	Coins  sdk.Coins
	Signer sdk.AccAddress
}

// NewMsgRefillBond is the constructor function for MsgRefillBond.
func NewMsgRefillBond(id string, denom string, amount int64, signer sdk.AccAddress) MsgRefillBond {
	return MsgRefillBond{
		ID:     ID(id),
		Coins:  sdk.NewCoins(sdk.NewInt64Coin(denom, amount)),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgRefillBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRefillBond) Type() string { return "refill" }

// ValidateBasic Implements Msg.
func (msg MsgRefillBond) ValidateBasic() sdk.Error {

	if string(msg.ID) == "" {
		return sdk.ErrInternal("Invalid bond ID.")
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
func (msg MsgRefillBond) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRefillBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgWithdrawBond defines a withdraw (funds from) bond message.
type MsgWithdrawBond struct {
	ID     ID
	Coins  sdk.Coins
	Signer sdk.AccAddress
}

// NewMsgWithdrawBond is the constructor function for MsgWithdrawBond.
func NewMsgWithdrawBond(id string, denom string, amount int64, signer sdk.AccAddress) MsgWithdrawBond {
	return MsgWithdrawBond{
		ID:     ID(id),
		Coins:  sdk.NewCoins(sdk.NewInt64Coin(denom, amount)),
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgWithdrawBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgWithdrawBond) Type() string { return "withdraw" }

// ValidateBasic Implements Msg.
func (msg MsgWithdrawBond) ValidateBasic() sdk.Error {

	if string(msg.ID) == "" {
		return sdk.ErrInternal("Invalid bond ID.")
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
func (msg MsgWithdrawBond) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgWithdrawBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgClear defines a MsgClear message.
type MsgClear struct {
	Signer sdk.AccAddress
}

// NewMsgClear is the constructor function for MsgClear.
func NewMsgClear(signer sdk.AccAddress) MsgClear {
	return MsgClear{
		Signer: signer,
	}
}

// Route Implements Msg.
func (msg MsgClear) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgClear) Type() string { return "clear" }

// ValidateBasic Implements Msg.
func (msg MsgClear) ValidateBasic() sdk.Error {

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgClear) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgClear) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
