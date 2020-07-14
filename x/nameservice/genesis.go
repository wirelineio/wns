//
// Copyright 2019 Wireline, Inc.
//

package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

type AuthorityEntry struct {
	Name  string              `json:"name" yaml:"name"`
	Entry types.NameAuthority `json:"record" yaml:"record"`
}

type NameEntry struct {
	Name  string           `json:"name" yaml:"name"`
	Entry types.NameRecord `json:"record" yaml:"record"`
}

type GenesisState struct {
	Params      types.Params      `json:"params" yaml:"params"`
	Authorities []AuthorityEntry  `json:"authorities" yaml:"authorities"`
	Names       []NameEntry       `json:"names" yaml:"names"`
	Records     []types.RecordObj `json:"records" yaml:"records"`
}

func NewGenesisState(params types.Params, authorities []AuthorityEntry, names []NameEntry, records []types.RecordObj) GenesisState {
	return GenesisState{
		Params:      params,
		Authorities: authorities,
		Names:       names,
		Records:     records,
	}
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

	for _, record := range data.Records {
		obj := record.ToRecord()
		keeper.PutRecord(ctx, obj)

		// Add to record expiry queue if expiry time is in the future.
		if obj.ExpiryTime.After(ctx.BlockTime()) {
			keeper.InsertRecordExpiryQueue(ctx, obj)
		}

		// Note: Bond genesis runs first, so bonds will already be present.
		if record.BondID != "" {
			keeper.AddBondToRecordIndexEntry(ctx, record.BondID, record.ID)
		}
	}

	for _, authorityEntry := range data.Authorities {
		keeper.SetNameAuthority(ctx, authorityEntry.Name, authorityEntry.Entry)
	}

	for _, nameEntry := range data.Names {
		keeper.SetNameRecord(ctx, nameEntry.Name, nameEntry.Entry.ID)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)

	authorities := keeper.ListNameAuthorityRecords(ctx)
	authorityEntries := []AuthorityEntry{}
	for name, nameAuthorityRecord := range authorities {
		authorityEntries = append(authorityEntries, AuthorityEntry{name, nameAuthorityRecord})
	}

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

	return GenesisState{
		Params:      params,
		Authorities: authorityEntries,
		Names:       nameEntries,
		Records:     recordEntries,
	}
}
