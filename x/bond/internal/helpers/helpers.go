//
// Copyright 2019 Wireline, Inc.
//

package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BondID simplifies generation of bond IDs.
type BondID struct {
	Address  sdk.Address
	AccNum   uint64
	Sequence uint64
}

// Generate creates the bond ID.
func (bondID BondID) Generate() string {
	hasher := sha256.New()
	str := fmt.Sprintf("%s:%d:%d", bondID.Address.String(), bondID.AccNum, bondID.Sequence)
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}
