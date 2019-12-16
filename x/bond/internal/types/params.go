//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	"fmt"

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

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultMaxBondAmount)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Max Bond Amount: %d`, p.MaxBondAmount)
}

// Validate a set of params.
func (p Params) Validate() error {
	if p.MaxBondAmount <= 0 {
		return fmt.Errorf("bond parameter MaxBondAmount must be a positive integer")
	}

	return nil
}
