//
// Copyright 2019 Wireline, Inc.
//

package bond

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		default:
			errMsg := fmt.Sprintf("Unrecognized bond Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgCreateBond.
func handleMsgCreateBond(ctx sdk.Context, keeper Keeper, msg types.MsgCreateBond) sdk.Result {
	bond, err := keeper.CreateBond(ctx, msg.Signer, msg.Coins)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(bond.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle handleMsgRefillBond.
func handleMsgRefillBond(ctx sdk.Context, keeper Keeper, msg types.MsgRefillBond) sdk.Result {
	bond, err := keeper.RefillBond(ctx, msg.ID, msg.Signer, msg.Coins)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(bond.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle handleMsgWithdrawBond.
func handleMsgWithdrawBond(ctx sdk.Context, keeper Keeper, msg types.MsgWithdrawBond) sdk.Result {
	bond, err := keeper.WithdrawBond(ctx, msg.ID, msg.Signer, msg.Coins)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(bond.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle handleMsgCancelBond.
func handleMsgCancelBond(ctx sdk.Context, keeper Keeper, msg types.MsgCancelBond) sdk.Result {
	bond, err := keeper.CancelBond(ctx, msg.ID, msg.Signer)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(bond.ID),
		Events: ctx.EventManager().Events(),
	}
}
