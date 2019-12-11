//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ID for bonds.
type ID string

// Bond represents funds deposited by an account for record rent payments.
type Bond struct {
	ID      ID
	Owner   string
	Balance sdk.Coins
}
