//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// Nameservice params default values.
const (
	// DefaultRecordAnnualRent is the default record rent for 1 year.
	DefaultRecordAnnualRent string = "1wire"
)

// nolint - Keys for parameter access
var (
	KeyRecordAnnualRent = []byte("RecordAnnualRent")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for nameservice
type Params struct {
	RecordAnnualRent string `json:"record_annual_rent" yaml:"record_annual_rent"`
}

// NewParams creates a new Params instance
func NewParams(recordAnnualRent string) Params {

	return Params{
		RecordAnnualRent: recordAnnualRent,
	}
}

// ParamSetPairs - implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyRecordAnnualRent, Value: &p.RecordAnnualRent},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultRecordAnnualRent)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Record Annual Rent: %s`, p.RecordAnnualRent)
}

// Validate a set of params.
func (p Params) Validate() error {
	if p.RecordAnnualRent == "" {
		return fmt.Errorf("nameservice parameter RecordAnnualRent can't be an empty string")
	}

	return nil
}
