//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Bond params default values.
const (
	// DefaultMaxBondAmount is the default maximum amount a bond can hold.
	DefaultMaxBondAmount int64 = 10000
)

// nolint - Keys for parameter access
var (
	KeyMaxBondAmount = []byte("MaxBondAmount")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for the bond module.
type Params struct {
	MaxBondAmount int64 `json:"max_bond_amount" yaml:"max_bond_amount"`
}

// NewParams creates a new Params instance
func NewParams(maxBondAmount int64) Params {

	return Params{
		MaxBondAmount: maxBondAmount,
	}
}

// ParamSetPairs - implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMaxBondAmount, Value: &p.MaxBondAmount},
	}
}

// Equal returns a boolean determining if two Param types are identical.
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultMaxBondAmount)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Max Bond Amount: %d`, p.MaxBondAmount)
}

// MustUnmarshalParams - unmarshal the current bond params value from store key or panic.
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// UnmarshalParams - unmarshal the current bond params value from store key.
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// Validate a set of params.
func (p Params) Validate() error {
	if p.MaxBondAmount <= 0 {
		return fmt.Errorf("bond parameter MaxBondAmount must be a positive integer")
	}

	return nil
}
