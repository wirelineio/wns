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

	if !msg.Coins.IsValid() {
		return sdk.ErrInvalidCoins("Invalid bond amount.")
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
