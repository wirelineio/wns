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
		switch msg := msg.(type) {
		case types.MsgCreateBond:
			return handleMsgCreateBond(ctx, keeper, msg)
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
		return sdk.ErrInsufficientCoins("Insufficient bond amount.").Result()
	}

	// Move funds into the bond account module.
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, ownerAddress, types.ModuleName, msg.Coins)
	if err != nil {
		return err.Result()
	}

	// Generate bond ID.
	account := keeper.AccountKeeper.GetAccount(ctx, ownerAddress)
	bondID := helpers.BondID{
		Address:  ownerAddress,
		AccNum:   account.GetAccountNumber(),
		Sequence: account.GetSequence(),
	}.Generate()

	// Save bond in store.
	keeper.CreateBond(ctx, types.Bond{ID: types.ID(bondID), Owner: ownerAddress.String(), Balance: msg.Coins})

	return sdk.Result{}
}

// Handle handleMsgClear.
func handleMsgClear(ctx sdk.Context, keeper Keeper, msg types.MsgClear) sdk.Result {
	keeper.Clear(ctx)
	return sdk.Result{}
}
