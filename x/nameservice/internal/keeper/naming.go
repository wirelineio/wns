//
// Copyright 2020 Wireline, Inc.
//

package keeper

import (
	"fmt"
	"net/url"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/wns/x/auction"
	"github.com/wirelineio/wns/x/nameservice/internal/helpers"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// ProcessReserveAuthority reserves a name authority.
func (k Keeper) ProcessReserveAuthority(ctx sdk.Context, msg types.MsgReserveAuthority) (string, sdk.Error) {
	wrn := fmt.Sprintf("wrn://%s", msg.Name)

	parsedWRN, err := url.Parse(wrn)
	if err != nil {
		return "", sdk.ErrInternal("Invalid name.")
	}

	name := parsedWRN.Host
	if fmt.Sprintf("wrn://%s", name) != wrn {
		return "", sdk.ErrInternal("Invalid name.")
	}

	// Check if name already reserved.
	if k.HasNameAuthority(ctx, name) {
		return "", sdk.ErrInternal("Name already exists.")
	}

	if strings.Contains(name, ".") {
		return k.ProcessReserveSubAuthority(ctx, name, msg)
	}

	// Reserve name with signer as owner.
	sdkErr := k.createAuthority(ctx, name, msg.Signer, true)
	if sdkErr != nil {
		return "", sdkErr
	}

	return name, nil
}

func (k Keeper) createAuthority(ctx sdk.Context, name string, owner sdk.AccAddress, isRoot bool) sdk.Error {
	ownerAccount := k.accountKeeper.GetAccount(ctx, owner)
	if ownerAccount == nil {
		return sdk.ErrUnknownAddress("Account not found.")
	}

	pubKey := ownerAccount.GetPubKey()
	if pubKey == nil {
		return sdk.ErrInvalidPubKey("Account public key not set.")
	}

	auctionID := auction.ID("")
	status := types.AuthorityActive

	// Create auction if root authority and name auctions are enabled.
	if isRoot && k.NameAuctionsEnabled(ctx) {
		commitFee, err := sdk.ParseCoin(k.NameAuctionCommitFee(ctx))
		if err != nil {
			return sdk.ErrInvalidCoins("Invalid name auction commit fee.")
		}

		revealFee, err := sdk.ParseCoin(k.NameAuctionRevealFee(ctx))
		if err != nil {
			return sdk.ErrInvalidCoins("Invalid name auction reveal fee.")
		}

		minimumBid, err := sdk.ParseCoin(k.NameAuctionMinimumBid(ctx))
		if err != nil {
			return sdk.ErrInvalidCoins("Invalid name auction minimum bid.")
		}

		// Create an auction.
		msg := auction.NewMsgCreateAuction(auction.Params{
			CommitsDuration: k.NameAuctionCommitsDuration(ctx),
			RevealsDuration: k.NameAuctionRevealsDuration(ctx),
			CommitFee:       commitFee,
			RevealFee:       revealFee,
			MinimumBid:      minimumBid,
		}, owner)

		// TODO(ashwin): Perhaps consume extra gas for auction creation.
		auction, sdkErr := k.auctionKeeper.CreateAuction(ctx, msg)
		if sdkErr != nil {
			return sdkErr
		}

		// TODO(ashwin): Create auction ID -> name authority index.

		status = types.AuthorityUnderAuction
		auctionID = auction.ID
	}

	authority := types.NameAuthority{
		Height:         ctx.BlockHeight(),
		OwnerAddress:   owner.String(),
		OwnerPublicKey: helpers.BytesToBase64(pubKey.Bytes()),
		Status:         status,
		AuctionID:      auctionID,
	}

	k.SetNameAuthority(ctx, name, authority)

	return nil
}

// ProcessReserveSubAuthority reserves a sub-authority.
func (k Keeper) ProcessReserveSubAuthority(ctx sdk.Context, name string, msg types.MsgReserveAuthority) (string, sdk.Error) {
	// Get parent authority name.
	names := strings.Split(name, ".")
	parent := strings.Join(names[1:], ".")

	// Check if parent authority exists.
	parentAuthority := k.GetNameAuthority(ctx, parent)
	if parentAuthority == nil {
		return name, sdk.ErrInternal("Parent authority not found.")
	}

	// Sub-authority creator needs to be the owner of the parent authority.
	if parentAuthority.OwnerAddress != msg.Signer.String() {
		return name, sdk.ErrUnauthorized("Access denied.")
	}

	// Sub-authority owner defaults to parent authority owner.
	subAuthorityOwner := msg.Signer
	if !msg.Owner.Empty() {
		// Override sub-authority owner if provided in message.
		subAuthorityOwner = msg.Owner
	}

	sdkErr := k.createAuthority(ctx, name, subAuthorityOwner, false)
	if sdkErr != nil {
		return "", sdkErr
	}

	return name, nil
}

func (k Keeper) checkWRN(ctx sdk.Context, signer sdk.AccAddress, inputWRN string) sdk.Error {
	parsedWRN, err := url.Parse(inputWRN)
	if err != nil {
		return sdk.ErrInternal("Invalid WRN.")
	}

	name := parsedWRN.Host
	formattedWRN := fmt.Sprintf("wrn://%s%s", name, parsedWRN.RequestURI())
	if formattedWRN != inputWRN {
		return sdk.ErrInternal("Invalid WRN.")
	}

	// Check authority record.
	authority := k.GetNameAuthority(ctx, name)
	if authority == nil {
		return sdk.ErrInternal("Name authority not found.")
	}

	if authority.OwnerAddress != signer.String() {
		return sdk.ErrUnauthorized("Access denied.")
	}

	return nil
}

// ProcessSetName creates a WRN -> Record ID mapping.
func (k Keeper) ProcessSetName(ctx sdk.Context, msg types.MsgSetName) sdk.Error {
	err := k.checkWRN(ctx, msg.Signer, msg.WRN)
	if err != nil {
		return err
	}

	nameRecord := k.GetNameRecord(ctx, msg.WRN)
	if nameRecord != nil && nameRecord.ID == msg.ID {
		// Already pointing to same ID, no-op.
		return nil
	}

	k.SetNameRecord(ctx, msg.WRN, msg.ID)

	return nil
}

// ProcessDeleteName removes a WRN -> Record ID mapping.
func (k Keeper) ProcessDeleteName(ctx sdk.Context, msg types.MsgDeleteName) sdk.Error {
	err := k.checkWRN(ctx, msg.Signer, msg.WRN)
	if err != nil {
		return err
	}

	if !k.HasNameRecord(ctx, msg.WRN) {
		return sdk.ErrInternal("Name not found.")
	}

	// Set CID to empty string.
	k.SetNameRecord(ctx, msg.WRN, "")

	return nil
}
