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

	DefaultNameAuctionsEnabled               = true
	DefaultCommitsDuration     time.Duration = time.Hour * 24
	DefaultRevealsDuration     time.Duration = time.Hour * 24
	DefaultCommitFee           string        = "1000000uwire"
	DefaultRevealFee           string        = "1000000uwire"
	DefaultMinimumBid          string        = "5000000uwire"
)

// nolint - Keys for parameter access
var (
	KeyRecordRent       = []byte("RecordRent")
	KeyRecordExpiryTime = []byte("RecordExpiryTime")

	KeyNameAuctions    = []byte("NameAuctionEnabled")
	KeyCommitsDuration = []byte("NameAuctionCommitsDuration")
	KeyRevealsDuration = []byte("NameAuctionRevealsDuration")
	KeyCommitFee       = []byte("NameAuctionCommitFee")
	KeyRevealFee       = []byte("NameAuctionRevealFee")
	KeyMinimumBid      = []byte("NameAuctionMinimumBid")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for nameservice
type Params struct {
	RecordRent       string        `json:"record_rent" yaml:"record_rent"`
	RecordExpiryTime time.Duration `json:"record_expiry_time" yaml:"record_expiry_time"`

	// Are name auctions enabled?
	NameAuctions    bool          `json:"name_auctions" yaml:"name_auctions"`
	CommitsDuration time.Duration `json:"name_auction_commits_duration" yaml:"name_auction_commits_duration"`
	RevealsDuration time.Duration `json:"name_auction_reveals_duration" yaml:"name_auction_reveals_duration"`
	CommitFee       string        `json:"name_auction_commit_fee" yaml:"name_auction_commit_fee"`
	RevealFee       string        `json:"name_auction_reveal_fee" yaml:"name_auction_reveal_fee"`
	MinimumBid      string        `json:"name_auction_minimum_bid" yaml:"name_auction_minimum_bid"`
}

// NewParams creates a new Params instance
func NewParams(recordRent string, recordExpiryTime time.Duration,
	nameAuctions bool, commitsDuration time.Duration, revealsDuration time.Duration,
	commitFee string, revealFee string, minimumBid string) Params {

	return Params{
		RecordRent:       recordRent,
		RecordExpiryTime: recordExpiryTime,

		NameAuctions:    nameAuctions,
		CommitsDuration: commitsDuration,
		RevealsDuration: revealsDuration,
		CommitFee:       commitFee,
		RevealFee:       revealFee,
		MinimumBid:      minimumBid,
	}
}

// ParamSetPairs - implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyRecordRent, Value: &p.RecordRent},
		{Key: KeyRecordExpiryTime, Value: &p.RecordExpiryTime},

		{Key: KeyNameAuctions, Value: &p.NameAuctions},
		{Key: KeyCommitsDuration, Value: &p.CommitsDuration},
		{Key: KeyRevealsDuration, Value: &p.RevealsDuration},
		{Key: KeyCommitFee, Value: &p.CommitFee},
		{Key: KeyRevealFee, Value: &p.RevealFee},
		{Key: KeyMinimumBid, Value: &p.MinimumBid},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultRecordRent, DefaultRecordExpiryTime,
		DefaultNameAuctionsEnabled, DefaultCommitsDuration, DefaultRevealsDuration,
		DefaultCommitFee, DefaultRevealFee, DefaultMinimumBid)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Record Rent                   : %v
  Record Expiry Time            : %v

  Name Auctions Enabled         : %v
  Name Auction Commits Duration : %v
  Name Auction Reveals Duration : %v
  Name Auction Commit Fee       : %v
  Name Auctions Reveal Fee      : %v
  Name Auctions Minimum Bid     : %v`,
		p.RecordRent, p.RecordExpiryTime,
		p.NameAuctions, p.CommitsDuration, p.RevealsDuration, p.CommitFee, p.RevealFee, p.MinimumBid)
}

// Validate a set of params.
func (p Params) Validate() error {
	if p.RecordRent == "" {
		return fmt.Errorf("nameservice parameter RecordRent can't be an empty string")
	}

	if p.RecordExpiryTime <= 0 {
		return fmt.Errorf("nameservice parameter RecordExpiryTime must be a positive integer")
	}

	if p.CommitsDuration <= 0 {
		return fmt.Errorf("nameservice parameter CommitsDuration must be a positive integer")
	}

	if p.RevealsDuration <= 0 {
		return fmt.Errorf("nameservice parameter RevealsDuration must be a positive integer")
	}

	if p.CommitFee == "" {
		return fmt.Errorf("nameservice parameter CommitFee can't be an empty string")
	}

	if p.RevealFee == "" {
		return fmt.Errorf("nameservice parameter RevealFee can't be an empty string")
	}

	if p.MinimumBid == "" {
		return fmt.Errorf("nameservice parameter MinimumBid can't be an empty string")
	}

	return nil
}
