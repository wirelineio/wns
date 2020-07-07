//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetRecord{}, "nameservice/SetRecord", nil)
	cdc.RegisterConcrete(MsgRenewRecord{}, "nameservice/RenewRecord", nil)

	cdc.RegisterConcrete(MsgReserveName{}, "nameservice/ReserveName", nil)
	cdc.RegisterConcrete(MsgSetName{}, "nameservice/SetName", nil)

	cdc.RegisterConcrete(MsgAssociateBond{}, "nameservice/AssociateBond", nil)
	cdc.RegisterConcrete(MsgDissociateBond{}, "nameservice/DissociateBond", nil)
	cdc.RegisterConcrete(MsgDissociateRecords{}, "nameservice/DissociateRecords", nil)
	cdc.RegisterConcrete(MsgReassociateRecords{}, "nameservice/ReassociateRecords", nil)
	cdc.RegisterConcrete(MsgClearRecords{}, "nameservice/ClearRecords", nil)
}
