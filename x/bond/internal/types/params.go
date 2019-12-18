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
	DefaultMaxBondAmount string = "10wire"
)

// nolint - Keys for parameter access
var (
	KeyMaxBondAmount = []byte("MaxBondAmount")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for the bond module.
type Params struct {
	MaxBondAmount string `json:"max_bond_amount" yaml:"max_bond_amount"`
}

// NewParams creates a new Params instance
func NewParams(maxBondAmount string) Params {

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
  Max Bond Amount: %s`, p.MaxBondAmount)
}

// Validate a set of params.
func (p Params) Validate() error {
	if p.MaxBondAmount == "" {
		return fmt.Errorf("bond parameter MaxBondAmount can't be an empty string")
	}

	return nil
}
