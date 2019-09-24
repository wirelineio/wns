package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetRecord defines a SetResource message.
type MsgSetRecord struct {
	Payload PayloadObj
	Signer  sdk.AccAddress
}

// NewMsgSetRecord is the constructor function for MsgSetRecord.
func NewMsgSetRecord(payload PayloadObj, signer sdk.AccAddress) MsgSetRecord {
	return MsgSetRecord{
		Payload: payload,
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

	record := RecordObjToRecord(msg.Payload.Record)
	if record.Type() == "" {
		return sdk.ErrInternal("Record 'type' not set.")
	}

	if record.Name() == "" {
		return sdk.ErrInternal("Record 'name' not set.")
	}

	if record.Version() == "" {
		return sdk.ErrInternal("Record 'version' not set.")
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
