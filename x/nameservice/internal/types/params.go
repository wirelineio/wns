//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// Nameservice params default values.
const (
	// DefaultRecordAnnualRent is the default record rent for 1 year.
	DefaultRecordAnnualRent string = "1wire"

	// DefaultRecordExpiryTime is the default record expiry time (1 year).
	DefaultRecordExpiryTime time.Duration = time.Hour * 24 * 7 * 365
)

// nolint - Keys for parameter access
var (
	KeyRecordAnnualRent = []byte("RecordAnnualRent")
	KeyRecordExpiryTime = []byte("RecordExpiryTime")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for nameservice
type Params struct {
	RecordAnnualRent string        `json:"record_annual_rent" yaml:"record_annual_rent"`
	RecordExpiryTime time.Duration `json:"record_expiry_time" yaml:"record_expiry_time"`
}

// NewParams creates a new Params instance
func NewParams(recordAnnualRent string, recordExpiryTime time.Duration) Params {

	return Params{
		RecordAnnualRent: recordAnnualRent,
		RecordExpiryTime: recordExpiryTime,
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
	return NewParams(DefaultRecordAnnualRent, DefaultRecordExpiryTime)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Record Annual Rent: %s
  Record Expiry Time: %s`, p.RecordAnnualRent, p.RecordExpiryTime)
}

// Validate a set of params.
func (p Params) Validate() error {
	if p.RecordAnnualRent == "" {
		return fmt.Errorf("nameservice parameter RecordAnnualRent can't be an empty string")
	}

	if p.RecordExpiryTime <= 0 {
		return fmt.Errorf("nameservice parameter RecordExpiryTime must be a positive integer")
	}

	return nil
}
