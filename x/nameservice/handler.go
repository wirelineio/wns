//
// Copyright 2019 Wireline, Inc.
//

package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSetRecord:
			return handleMsgSetRecord(ctx, keeper, msg)
		case types.MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case types.MsgDeleteName:
			return handleMsgDeleteName(ctx, keeper, msg)
		case types.MsgReserveAuthority:
			return handleMsgReserveAuthority(ctx, keeper, msg)
		case types.MsgSetAuthorityBond:
			return handleMsgSetAuthorityBond(ctx, keeper, msg)
		case types.MsgAssociateBond:
			return handleMsgAssociateBond(ctx, keeper, msg)
		case types.MsgDissociateBond:
			return handleMsgDissociateBond(ctx, keeper, msg)
		case types.MsgDissociateRecords:
			return handleMsgDissociateRecords(ctx, keeper, msg)
		case types.MsgReassociateRecords:
			return handleMsgReassociateRecords(ctx, keeper, msg)
		case types.MsgRenewRecord:
			return handleMsgRenewRecord(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetRecord.
func handleMsgSetRecord(ctx sdk.Context, keeper Keeper, msg types.MsgSetRecord) sdk.Result {
	record, err := keeper.ProcessSetRecord(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(record.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgRenewRecord.
func handleMsgRenewRecord(ctx sdk.Context, keeper Keeper, msg types.MsgRenewRecord) sdk.Result {
	record, err := keeper.ProcessRenewRecord(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(record.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgAssociateBond.
func handleMsgAssociateBond(ctx sdk.Context, keeper Keeper, msg types.MsgAssociateBond) sdk.Result {
	record, err := keeper.ProcessAssociateBond(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(record.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgDissociateBond.
func handleMsgDissociateBond(ctx sdk.Context, keeper Keeper, msg types.MsgDissociateBond) sdk.Result {
	record, err := keeper.ProcessDissociateBond(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(record.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgDissociateRecords.
func handleMsgDissociateRecords(ctx sdk.Context, keeper Keeper, msg types.MsgDissociateRecords) sdk.Result {
	bond, err := keeper.ProcessDissociateRecords(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(bond.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgReassociateRecords.
func handleMsgReassociateRecords(ctx sdk.Context, keeper Keeper, msg types.MsgReassociateRecords) sdk.Result {
	bond, err := keeper.ProcessReassociateRecords(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(bond.ID),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgReserveName.
func handleMsgReserveAuthority(ctx sdk.Context, keeper Keeper, msg types.MsgReserveAuthority) sdk.Result {
	name, err := keeper.ProcessReserveAuthority(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(name),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgSetAuthorityBond.
func handleMsgSetAuthorityBond(ctx sdk.Context, keeper Keeper, msg types.MsgSetAuthorityBond) sdk.Result {
	name, err := keeper.ProcessSetAuthorityBond(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(name),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgSetName.
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg types.MsgSetName) sdk.Result {
	err := keeper.ProcessSetName(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(msg.WRN),
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgDeleteName.
func handleMsgDeleteName(ctx sdk.Context, keeper Keeper, msg types.MsgDeleteName) sdk.Result {
	err := keeper.ProcessDeleteName(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data:   []byte(msg.WRN),
		Events: ctx.EventManager().Events(),
	}
}
