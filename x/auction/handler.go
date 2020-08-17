//
// Copyright 2020 Wireline, Inc.
//

package auction

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/wns/x/auction/internal/types"
)

// NewHandler returns a handler for "auction" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateAuction:
			return handleMsgCreateAuction(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized auction Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgCreateAuction.
func handleMsgCreateAuction(ctx sdk.Context, keeper Keeper, msg types.MsgCreateAuction) sdk.Result {
	auction, err := keeper.CreateAuction(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(auction.ID),
		Events: ctx.EventManager().Events(),
	}
}
