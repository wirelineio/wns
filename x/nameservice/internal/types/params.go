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
	// DefaultRecordRent is the default record rent for 1 time period (see expiry time).
	DefaultRecordRent string = "1000000uwire"

	// DefaultRecordExpiryTime is the default record expiry time (1 year).
	DefaultRecordExpiryTime time.Duration = time.Hour * 24 * 365
)

// nolint - Keys for parameter access
var (
	KeyRecordRent       = []byte("RecordRent")
	KeyRecordExpiryTime = []byte("RecordExpiryTime")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for nameservice
type Params struct {
	RecordRent       string        `json:"record_rent" yaml:"record_rent"`
	RecordExpiryTime time.Duration `json:"record_expiry_time" yaml:"record_expiry_time"`
}

// NewParams creates a new Params instance
func NewParams(recordRent string, recordExpiryTime time.Duration) Params {

	return Params{
		RecordRent:       recordRent,
		RecordExpiryTime: recordExpiryTime,
	}
}

// ParamSetPairs - implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyRecordRent, Value: &p.RecordRent},
		{Key: KeyRecordExpiryTime, Value: &p.RecordExpiryTime},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultRecordRent, DefaultRecordExpiryTime)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Record Rent        : %s
  Record Expiry Time : %s`, p.RecordRent, p.RecordExpiryTime)
}

// Validate a set of params.
func (p Params) Validate() error {
	if p.RecordRent == "" {
		return fmt.Errorf("nameservice parameter RecordRent can't be an empty string")
	}

	if p.RecordExpiryTime <= 0 {
		return fmt.Errorf("nameservice parameter RecordExpiryTime must be a positive integer")
	}

	return nil
}
