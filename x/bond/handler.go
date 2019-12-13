//
// Copyright 2019 Wireline, Inc.
//

package bond

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/wns/x/bond/internal/helpers"
	"github.com/wirelineio/wns/x/bond/internal/types"
)

// MaxBondBalance is the maximum amount a bond can hold.
// https://github.com/wirelineio/specs/blob/master/wns/testnet-mechanism.md#pricing
// TODO(ashwin): Needs to be made a param under consensus (https://github.com/wirelineio/wns/issues/88).
// TODO(ashwin): Figure out denom unit to use (https://github.com/wirelineio/wns/issues/123).
const MaxBondBalance int64 = 10000

// NewHandler returns a handler for "bond" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgCreateBond:
			return handleMsgCreateBond(ctx, keeper, msg)
		case types.MsgRefillBond:
			return handleMsgRefillBond(ctx, keeper, msg)
		case types.MsgWithdrawBond:
			return handleMsgWithdrawBond(ctx, keeper, msg)
		case types.MsgCancelBond:
			return handleMsgCancelBond(ctx, keeper, msg)
		case types.MsgClear:
			return handleMsgClear(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized bond Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgCreateBond.
func handleMsgCreateBond(ctx sdk.Context, keeper Keeper, msg types.MsgCreateBond) sdk.Result {
	ownerAddress := msg.Signer

	// Check if account has funds.
	if !keeper.CoinKeeper.HasCoins(ctx, ownerAddress, msg.Coins) {
		return sdk.ErrInsufficientCoins("Insufficient funds.").Result()
	}

	// Generate bond ID.
	account := keeper.AccountKeeper.GetAccount(ctx, ownerAddress)
	bondID := helpers.BondID{
		Address:  ownerAddress,
		AccNum:   account.GetAccountNumber(),
		Sequence: account.GetSequence(),
	}.Generate()

	bond := types.Bond{ID: types.ID(bondID), Owner: ownerAddress.String(), Balance: msg.Coins}
	if helpers.AnyCoinAmountExceeds(bond.Balance, MaxBondBalance) {
		return sdk.ErrInternal("Max bond amount exceeded.").Result()
	}

	// Move funds into the bond account module.
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, ownerAddress, types.ModuleName, msg.Coins)
	if err != nil {
		return err.Result()
	}

	// Save bond in store.
	keeper.SaveBond(ctx, bond)

	return sdk.Result{}
}

// Handle handleMsgRefillBond.
func handleMsgRefillBond(ctx sdk.Context, keeper Keeper, msg types.MsgRefillBond) sdk.Result {

	if !keeper.HasBond(ctx, msg.ID) {
		return sdk.ErrInternal("Bond not found.").Result()
	}

	ownerAddress := msg.Signer
	bond := keeper.GetBond(ctx, msg.ID)
	if bond.Owner != ownerAddress.String() {
		return sdk.ErrUnauthorized("Bond owner mismatch.").Result()
	}

	// Check if account has funds.
	if !keeper.CoinKeeper.HasCoins(ctx, ownerAddress, msg.Coins) {
		return sdk.ErrInsufficientCoins("Insufficient funds.").Result()
	}

	updatedBalance := bond.Balance.Add(msg.Coins)
	if helpers.AnyCoinAmountExceeds(updatedBalance, MaxBondBalance) {
		return sdk.ErrInternal("Max bond amount exceeded.").Result()
	}

	// Move funds into the bond account module.
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, ownerAddress, types.ModuleName, msg.Coins)
	if err != nil {
		return err.Result()
	}

	// Update bond balance and save.
	bond.Balance = updatedBalance
	keeper.SaveBond(ctx, bond)

	return sdk.Result{}
}

// Handle handleMsgWithdrawBond.
func handleMsgWithdrawBond(ctx sdk.Context, keeper Keeper, msg types.MsgWithdrawBond) sdk.Result {

	if !keeper.HasBond(ctx, msg.ID) {
		return sdk.ErrInternal("Bond not found.").Result()
	}

	ownerAddress := msg.Signer
	bond := keeper.GetBond(ctx, msg.ID)
	if bond.Owner != ownerAddress.String() {
		return sdk.ErrUnauthorized("Bond owner mismatch.").Result()
	}

	updatedBalance, isNeg := bond.Balance.SafeSub(msg.Coins)
	if isNeg {
		return sdk.ErrInsufficientCoins("Insufficient bond balance.").Result()
	}

	// Move funds from the bond into the account.
	err := keeper.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddress, msg.Coins)
	if err != nil {
		return err.Result()
	}

	// Update bond balance and save.
	bond.Balance = updatedBalance
	keeper.SaveBond(ctx, bond)

	return sdk.Result{}
}

// Handle handleMsgCancelBond.
func handleMsgCancelBond(ctx sdk.Context, keeper Keeper, msg types.MsgCancelBond) sdk.Result {

	if !keeper.HasBond(ctx, msg.ID) {
		return sdk.ErrInternal("Bond not found.").Result()
	}

	ownerAddress := msg.Signer
	bond := keeper.GetBond(ctx, msg.ID)
	if bond.Owner != ownerAddress.String() {
		return sdk.ErrUnauthorized("Bond owner mismatch.").Result()
	}

	// Check if bond is associated with any records.
	if keeper.RecordKeeper.BondHasAssociatedRecords(ctx, msg.ID) {
		return sdk.ErrUnauthorized("Bond has associated records.").Result()
	}

	// Move funds from the bond into the account.
	err := keeper.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddress, bond.Balance)
	if err != nil {
		return err.Result()
	}

	keeper.DeleteBond(ctx, bond)

	return sdk.Result{}
}

// Handle handleMsgClear.
func handleMsgClear(ctx sdk.Context, keeper Keeper, msg types.MsgClear) sdk.Result {
	keeper.Clear(ctx)
	return sdk.Result{}
}
