//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RecordKeeper interface {
	BondHasAssociatedRecords(ctx sdk.Context, bondID ID) bool
}
