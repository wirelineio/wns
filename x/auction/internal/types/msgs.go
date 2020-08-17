//
// Copyright 2020 Wireline, Inc.
//

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/types/time"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgCreateAuction defines a create auction message.
type MsgCreateAuction struct {
	CommitsEndTime int64          `json:"commitsEndTime,omitempty"`
	RevealsEndTime int64          `json:"revealsEndTime,omitempty"`
	CommitFee      sdk.Coin       `json:"commitFee,omitempty"`
	RevealFee      sdk.Coin       `json:"revealFee,omitempty"`
	MinimumBid     sdk.Coin       `json:"minimumBid,omitempty"`
	Signer         sdk.AccAddress `json:"signer"`
}

// NewMsgCreateAuction is the constructor function for MsgCreateAuction.
func NewMsgCreateAuction(commitsEndTime int64, revealsEndTime int64,
	commitFee sdk.Coin, revealFee sdk.Coin, minimumBid sdk.Coin,
	signer sdk.AccAddress) MsgCreateAuction {

	return MsgCreateAuction{
		CommitsEndTime: commitsEndTime,
		RevealsEndTime: revealsEndTime,
		CommitFee:      commitFee,
		RevealFee:      revealFee,
		MinimumBid:     minimumBid,
		Signer:         signer,
	}
}

// Route Implements Msg.
func (msg MsgCreateAuction) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateAuction) Type() string { return "create" }

// ValidateBasic Implements Msg.
func (msg MsgCreateAuction) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	now := time.Now().Unix()
	if msg.CommitsEndTime < now {
		return sdk.ErrInternal("Commit phase end time can't be in the past.")
	}

	if msg.RevealsEndTime < now {
		return sdk.ErrInternal("Reveal phase end time can't be in the past.")
	}

	if msg.RevealsEndTime <= msg.CommitsEndTime {
		return sdk.ErrInternal("Reveal phase end time must be greater than commit phase end time.")
	}

	if !msg.MinimumBid.IsPositive() {
		return sdk.ErrInternal("Minimum bid should be greater than zero.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCreateAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgCommitBid defines a commit bid message.
type MsgCommitBid struct {
	AuctionID  ID             `json:"auctionId,omitempty"`
	CommitHash string         `json:"commit,omitempty"`
	AuctionFee sdk.Coin       `json:"auctionFee,omitempty"`
	Signer     sdk.AccAddress `json:"signer"`
}

// NewMsgCommitBid is the constructor function for MsgCommitBid.
func NewMsgCommitBid(auctionID string, commitHash string, auctionFee sdk.Coin,
	signer sdk.AccAddress) MsgCommitBid {

	return MsgCommitBid{
		AuctionID:  ID(auctionID),
		CommitHash: commitHash,
		AuctionFee: auctionFee,
		Signer:     signer,
	}
}

// Route Implements Msg.
func (msg MsgCommitBid) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCommitBid) Type() string { return "commit" }

// ValidateBasic Implements Msg.
func (msg MsgCommitBid) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	if msg.AuctionID == "" {
		return sdk.ErrInternal("Invalid auction ID.")
	}

	if msg.CommitHash == "" {
		return sdk.ErrInternal("Invalid commit hash.")
	}

	if !msg.AuctionFee.IsPositive() {
		return sdk.ErrInternal("Auction fee should be greater than zero.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCommitBid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCommitBid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgRevealBid defines a reveal bid message.
type MsgRevealBid struct {
	AuctionID ID             `json:"auctionId,omitempty"`
	Reveal    string         `json:"commit,omitempty"`
	Signer    sdk.AccAddress `json:"signer"`
}

// NewMsgRevealBid is the constructor function for MsgRevealBid.
func NewMsgRevealBid(auctionID string, reveal string, signer sdk.AccAddress) MsgRevealBid {

	return MsgRevealBid{
		AuctionID: ID(auctionID),
		Reveal:    reveal,
		Signer:    signer,
	}
}

// Route Implements Msg.
func (msg MsgRevealBid) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRevealBid) Type() string { return "reveal" }

// ValidateBasic Implements Msg.
func (msg MsgRevealBid) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	if msg.AuctionID == "" {
		return sdk.ErrInternal("Invalid auction ID.")
	}

	if msg.Reveal == "" {
		return sdk.ErrInternal("Invalid commit hash.")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRevealBid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRevealBid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
