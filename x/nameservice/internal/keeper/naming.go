//
// Copyright 2020 Wireline, Inc.
//

package keeper

import (
	"fmt"
	"net/url"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	sdkErr := k.createAuthority(ctx, name, msg.Signer)
	if sdkErr != nil {
		return "", sdkErr
	}

	return name, nil
}

func (k Keeper) createAuthority(ctx sdk.Context, name string, owner sdk.AccAddress) sdk.Error {
	ownerAccount := k.accountKeeper.GetAccount(ctx, owner)
	if ownerAccount == nil {
		return sdk.ErrUnknownAddress("Account not found.")
	}

	pubKey := ownerAccount.GetPubKey()
	if pubKey == nil {
		return sdk.ErrInvalidPubKey("Account public key not set.")
	}

	k.SetNameAuthority(
		ctx,
		name,
		owner.String(),
		helpers.BytesToBase64(pubKey.Bytes()),
	)

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

	sdkErr := k.createAuthority(ctx, name, subAuthorityOwner)
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
