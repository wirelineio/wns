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
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateBond{}, "bond/CreateBond", nil)
	cdc.RegisterConcrete(MsgRefillBond{}, "bond/RefillBond", nil)
	cdc.RegisterConcrete(MsgClear{}, "bond/Clear", nil)
}
