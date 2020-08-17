//
// Copyright 2020 Wireline, Inc.
//

package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the parameters for the auction module.
type Params struct {
	// Duration of commits phase in seconds.
	CommitsDuration time.Duration `json:"commits_duration"`

	// Duration of reveals phase in seconds.
	RevealsDuration time.Duration `json:"reveals_duration"`

	// Commit and reveal fees.
	CommitFee sdk.Coin `json:"commit_fee"`
	RevealFee sdk.Coin `json:"reveal_fee"`

	MinimumBid sdk.Coin `json:"minimum_bid"`
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// ParamSetPairs - implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return ""
}

// Validate a set of params.
func (p Params) Validate() error {
	return nil
}
