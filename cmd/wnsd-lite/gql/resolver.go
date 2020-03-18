//
// Copyright 2020 Wireline, Inc.
//

package gql

import (
	"context"
)

// Resolver is the GQL query resolver.
type Resolver struct{}

type queryResolver struct{ *Resolver }

// Query is the entry point to query execution.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *queryResolver) GetRecordsByIds(ctx context.Context, ids []string) ([]*Record, error) {
	records := make([]*Record, len(ids))
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
func (r *queryResolver) QueryRecords(ctx context.Context, attributes []*KeyValueInput) ([]*Record, error) {

	return nil, nil
}

// ResolveRecords resolves records by ref/WRN, with semver range support.
func (r *queryResolver) ResolveRecords(ctx context.Context, refs []string) ([]*Record, error) {

	return nil, nil
}

func (r *queryResolver) GetStatus(ctx context.Context) (*Status, error) {
	return &Status{}, nil
}

func (r *queryResolver) GetRecord(ctx context.Context, id string) (*Record, error) {

	return nil, nil
}
