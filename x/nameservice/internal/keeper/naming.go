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

	if strings.Contains(name, ".") {
		return "", sdk.ErrInternal(("Invalid name (dot is currently not allowed in root authority names)."))
	}

	// Check if name already reserved.
	if k.HasNameAuthority(ctx, name) {
		return "", sdk.ErrInternal("Name already exists.")
	}

	// Reserve name with signer as owner.
	account := k.AccountKeeper.GetAccount(ctx, msg.Signer)
	if account == nil {
		return "", sdk.ErrUnknownAddress("Account not found.")
	}

	k.SetNameAuthority(
		ctx,
		name,
		msg.Signer.String(),
		helpers.BytesToBase64(account.GetPubKey().Bytes()),
	)

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
