//
// Copyright 2019 Wireline, Inc.
//

package bond

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wirelineio/wns/x/bond/internal/types"
)

type GenesisState struct {
	Params types.Params `json:"params" yaml:"params"`
}

func NewGenesisState(params types.Params) GenesisState {
	return GenesisState{Params: params}
}

func ValidateGenesis(data GenesisState) error {
	err := data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{Params: types.DefaultParams()}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)

	return GenesisState{Params: params}
}
