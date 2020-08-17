//
// Copyright 2020 Wireline, Inc.
//

package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Auction status values.
const (
	// Auction is in commit phase.
	AuctionStatusCommitPhase = 0x1

	// Auction is in reveal phase.
	AuctionStatusRevealPhase = 0x2

	// Auction finished (i.e. winner selected).
	AuctionStatusFinished = 0xf
)

// Bid status values.
const (
	BidStatusCommitted = 0x1
	BidStatusRevealed  = 0x2
	BidStatusExpired   = 0x3
)

// ID for auctions.
type ID string

// Auction is a 2nd price sealed-bid on-chain auction.
type Auction struct {
	ID     ID   `json:"id,omitempty"`
	Status int8 `json:"status,omitempty"`

	// Creator of the auction.
	OwnerAddress sdk.AccAddress `json:"ownerAddress,omitempty"`

	// Auction create time.
	CreateTime time.Time `json:"createTime,omitempty"`

	// Time when the commit phase ends.
	CommitsEndTime time.Time `json:"commitsEndTime,omitempty"`

	// Time when the reveal phase ends.
	RevealsEndTime time.Time `json:"revealsEndTime,omitempty"`

	// Commit Fee + Reveal Fee both need to be paid when committing a bid.
	// Reveal Fee is returned ONLY if the bid is revealed.
	CommitFee sdk.Coin `json:"commitFee,omitempty"`
	RevealFee sdk.Coin `json:"revealFee,omitempty"`

	// Minimum bid for a valid commit.
	MinimumBid sdk.Coin `json:"minimumBid,omitempty"`

	// Winner address.
	WinnerAddress sdk.Address `json:"winnerAddress,omitempty"`

	// Winning bid, i.e. highest bid.
	WinnerBid sdk.Coin `json:"winnerBid,omitempty"`

	// Amount winner actually pays, i.e. 2nd highest bid.
	// As it's a 2nd price auction.
	WinnerPrice sdk.Address `json:"winnerPrice,omitempty"`
}

// Bid represents a sealed bid (commit) made during the auction.
type Bid struct {
	AuctionID     ID          `json:"auctionId,omitempty"`
	BidderAddress sdk.Address `json:"bidderAddress,omitempty"`
	Status        int8        `json:"status,omitempty"`
	AuctionFee    sdk.Coin    `json:"auctionFee,omitempty"`
	CommitHash    string      `json:"commitHash,omitempty"`
	CommitTime    time.Time   `json:"commitTime,omitempty"`
	Reveal        string      `json:"reveal,omitempty"`
	RevealTime    time.Time   `json:"revealTime,omitempty"`
	BidAmount     sdk.Coin    `json:"bidAmount,omitempty"`
}

// AuctionID simplifies generation of auction IDs.
type AuctionID struct {
	Address  sdk.Address
	AccNum   uint64
	Sequence uint64
}

// Generate creates the auction ID.
func (auctionID AuctionID) Generate() string {
	hasher := sha256.New()
	str := fmt.Sprintf("%s:%d:%d", auctionID.Address.String(), auctionID.AccNum, auctionID.Sequence)
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}
