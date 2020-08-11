//
// Copyright 2020 Wireline, Inc.
//

package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint - Keys for parameter access
var (
	KeyMaxAuctionAmount = []byte("MaxAuctionAmount")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for the auction module.
type Params struct {
	MaxAuctionAmount string `json:"max_auction_amount" yaml:"max_auction_amount"`
}

// NewParams creates a new Params instance
func NewParams(maxAuctionAmount string) Params {
	return Params{
		MaxAuctionAmount: maxAuctionAmount,
	}
}

// ParamSetPairs - implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMaxAuctionAmount, Value: &p.MaxAuctionAmount},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams("")
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Max Auction Amount: %s`, p.MaxAuctionAmount)
}

// Validate a set of params.
func (p Params) Validate() error {
	_, err := sdk.ParseCoins(p.MaxAuctionAmount)
	if err != nil {
		return err
	}

	return nil
}
