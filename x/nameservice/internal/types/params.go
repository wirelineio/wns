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

// Equal returns a boolean determining if two Param types are identical.
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
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

// MustUnmarshalParams - unmarshal the current nameservice params value from store key or panic.
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// UnmarshalParams - unmarshal the current nameservice params value from store key.
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// Validate a set of params.
func (p Params) Validate() error {
	if p.RecordAnnualRent == "" {
		return fmt.Errorf("nameservice parameter RecordAnnualRent can't be an empty string")
	}

	return nil
}
