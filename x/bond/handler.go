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

// NewHandler returns a handler for "bond" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
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

func getMaxBondAmount(ctx sdk.Context, keeper Keeper) (sdk.Coins, error) {
	maxBondAmount, err := sdk.ParseCoin(keeper.MaxBondAmount(ctx))
	if err != nil {
		return nil, err
	}

	maxBondAmountMicroWire, err := sdk.ConvertCoin(maxBondAmount, types.MicroWire)
	if err != nil {
		return nil, err
	}

	return sdk.NewCoins(maxBondAmountMicroWire), nil
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

	maxBondAmountMicroWire, err := getMaxBondAmount(ctx, keeper)
	if err != nil {
		return sdk.ErrInternal("Invalid max bond amount.").Result()
	}

	bond := types.Bond{ID: types.ID(bondID), Owner: ownerAddress.String(), Balance: msg.Coins}
	if bond.Balance.IsAnyGT(maxBondAmountMicroWire) {
		return sdk.ErrInternal("Max bond amount exceeded.").Result()
	}

	// Move funds into the bond account module.
	sdkErr := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, ownerAddress, types.ModuleName, bond.Balance)
	if err != nil {
		return sdkErr.Result()
	}

	// Save bond in store.
	keeper.SaveBond(ctx, bond)

	return sdk.Result{
		Data:   []byte(bond.ID),
		Events: ctx.EventManager().Events(),
	}
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

	maxBondAmountMicroWire, err := getMaxBondAmount(ctx, keeper)
	if err != nil {
		return sdk.ErrInternal("Invalid max bond amount.").Result()
	}

	updatedBalance := bond.Balance.Add(msg.Coins)
	if updatedBalance.IsAnyGT(maxBondAmountMicroWire) {
		return sdk.ErrInternal("Max bond amount exceeded.").Result()
	}

	// Move funds into the bond account module.
	sdkErr := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, ownerAddress, types.ModuleName, msg.Coins)
	if err != nil {
		return sdkErr.Result()
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
