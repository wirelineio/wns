//
// Copyright 2020 Wireline, Inc.
//

package gql

import (
	"context"

	"github.com/wirelineio/wns/cmd/wnsd-lite/sync"
	baseGql "github.com/wirelineio/wns/gql"
	"github.com/wirelineio/wns/x/nameservice"
)

// Resolver is the GQL query resolver.
type Resolver struct {
	Keeper *sync.Keeper
}

type queryResolver struct{ *Resolver }

// Query is the entry point to query execution.
func (r *Resolver) Query() baseGql.QueryResolver {
	return &queryResolver{r}
}

func (r *queryResolver) GetRecordsByIds(ctx context.Context, ids []string) ([]*baseGql.Record, error) {
	records := make([]*baseGql.Record, len(ids))
	for index, id := range ids {
		record, err := r.GetRecord(ctx, id)
		if err != nil {
			return nil, err
		}

		records[index] = record
	}

	return records, nil
}

// QueryRecords filters records by K=V conditions.
func (r *queryResolver) QueryRecords(ctx context.Context, attributes []*baseGql.KeyValueInput) ([]*baseGql.Record, error) {
	gqlResponse := []*baseGql.Record{}

	var records = r.Keeper.MatchRecords(func(record *nameservice.Record) bool {
		return baseGql.MatchOnAttributes(record, attributes)
	})

	if baseGql.RequestedLatestVersionsOnly(attributes) {
		records = baseGql.GetLatestVersions(records)
	}

	for _, record := range records {
		gqlRecord, err := baseGql.GetGQLRecord(ctx, r, record)
		if err != nil {
			return nil, err
		}

		gqlResponse = append(gqlResponse, gqlRecord)
	}

	return gqlResponse, nil
}

// ResolveRecords resolves records by ref/WRN, with semver range support.
func (r *queryResolver) ResolveRecords(ctx context.Context, refs []string) ([]*baseGql.Record, error) {
	gqlResponse := []*baseGql.Record{}

	for _, ref := range refs {
		record := r.Keeper.ResolveWRN(ref)
		gqlRecord, err := baseGql.GetGQLRecord(ctx, r, record)
		if err != nil {
			return nil, err
		}

		gqlResponse = append(gqlResponse, gqlRecord)
	}

	return gqlResponse, nil
}

func (r *queryResolver) GetStatus(ctx context.Context) (*baseGql.Status, error) {
	return &baseGql.Status{}, nil
}

func (r *queryResolver) GetRecord(ctx context.Context, id string) (*baseGql.Record, error) {
	dbID := nameservice.ID(id)
	if r.Keeper.HasRecord(dbID) {
		record := r.Keeper.GetRecord(dbID)
		if !record.Deleted {
			return baseGql.GetGQLRecord(ctx, r, &record)
		}
	}

	return nil, nil
}
