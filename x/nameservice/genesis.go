//
// Copyright 2019 Wireline, Inc.
//

package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

type NameEntry struct {
	Name  string           `json:"name" yaml:"name"`
	Entry types.NameRecord `json:"record" yaml:"record"`
}

type GenesisState struct {
	Names   []NameEntry       `json:"names" yaml:"names"`
	Records []types.RecordObj `json:"records" yaml:"records"`
}

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, nameEntry := range data.Names {
		keeper.SetNameRecord(ctx, nameEntry.Name, nameEntry.Entry)
	}

	for _, record := range data.Records {
		keeper.PutRecord(ctx, record.ToRecord())
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	names := keeper.ListNameRecords(ctx)
	nameEntries := []NameEntry{}
	for name, nameRecord := range names {
		nameEntries = append(nameEntries, NameEntry{name, nameRecord})
	}

	records := keeper.ListRecords(ctx)
	recordEntries := []types.RecordObj{}
	for _, record := range records {
		recordEntries = append(recordEntries, record.ToRecordObj())
	}

	return GenesisState{nameEntries, recordEntries}
}
