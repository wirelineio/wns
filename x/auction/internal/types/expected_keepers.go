//
// Copyright 2020 Wireline, Inc.
//

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuctionUsageKeeper keep track of auction usage in other modules.
// Used to, for example, prevent deletion of a auction that's in use.
type AuctionUsageKeeper interface {
	ModuleName() string
	UsesAuction(ctx sdk.Context, auctionID ID) bool
}
