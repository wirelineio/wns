//
// Copyright 2019 Wireline, Inc.
//

package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wirelineio/wns/x/bond/internal/types"
)

// RegisterInvariants registers all bond module invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-accounts", BalanceInvariants(k))
}

// BalanceInvariants checks that the 'bond' and 'rent' module account balances are non-negative.
func BalanceInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := checkModuleNegativeBalance(ctx, k, types.ModuleName)
		if stop {
			return res, stop
		}

		return checkModuleNegativeBalance(ctx, k, types.RecordRentModuleAccountName)
	}
}

func checkModuleNegativeBalance(ctx sdk.Context, k Keeper, name string) (string, bool) {
	moduleAccount := k.SupplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	if moduleAccount.GetCoins().IsAnyNegative() {
		return sdk.FormatInvariant(types.ModuleName, "module-accounts", fmt.Sprintf("Module account '%s' has negative balance.", name)), true
	}

	return "", false
}

func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return BalanceInvariants(k)(ctx)
	}
}
