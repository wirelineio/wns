//
// Copyright 2020 Wireline, Inc.
//

package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ID for auctions.
type ID string

// Auction represents funds deposited by an account for record rent payments.
type Auction struct {
	ID      ID        `json:"id,omitempty"`
	Owner   string    `json:"owner,omitempty"`
	Balance sdk.Coins `json:"balance"`
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
