//
// Copyright 2020 Wireline, Inc.
//

package gql

import (
	"context"
	"os"
	"strconv"

	"github.com/wirelineio/wns/cmd/wnsd-lite/sync"
	baseGql "github.com/wirelineio/wns/gql"
	"github.com/wirelineio/wns/x/nameservice"
)

// LiteNodeDataPath is the path to the lite node data folder.
var LiteNodeDataPath = os.ExpandEnv("$HOME/.wire/wnsd-lite/data")

// Resolver is the GQL query resolver.
type Resolver struct {
	Keeper  *sync.Keeper
	LogFile string
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
	var records = r.Keeper.MatchRecords(func(record *nameservice.Record) bool {
		return baseGql.MatchOnAttributes(record, attributes)
	})

	return baseGql.QueryRecords(ctx, r, records, attributes)
}

// ResolveRecords resolves records by ref/WRN, with semver range support.
func (r *queryResolver) ResolveNames(ctx context.Context, names []string) (*baseGql.RecordResult, error) {
	gqlResponse := []*baseGql.Record{}

	for _, ref := range names {
		record := r.Keeper.ResolveWRN(ref)
		gqlRecord, err := baseGql.GetGQLRecord(ctx, r, record)
		if err != nil {
			return nil, err
		}

		gqlResponse = append(gqlResponse, gqlRecord)
	}

	result := baseGql.RecordResult{
		Meta: baseGql.ResultMeta{
			Height: string(r.Keeper.GetStatusRecord().LastSyncedHeight),
		},
		Records: gqlResponse,
	}

	return &result, nil
}

func (r *queryResolver) LookupAuthorities(ctx context.Context, names []string) (*baseGql.AuthorityResult, error) {
	// TODO(ashwin): Implement.
	return nil, nil
}

func (r *queryResolver) LookupNames(ctx context.Context, names []string) (*baseGql.NameResult, error) {
	// TODO(ashwin): Implement.
	return nil, nil
}

func (r *queryResolver) GetLogs(ctx context.Context, count *int) ([]string, error) {
	return baseGql.GetLogs(ctx, r.LogFile, count)
}

func (r *queryResolver) GetStatus(ctx context.Context) (*baseGql.Status, error) {
	statusRecord := r.Keeper.GetStatusRecord()

	diskUsage, err := baseGql.GetDiskUsage(LiteNodeDataPath)
	if err != nil {
		return nil, err
	}

	return &baseGql.Status{
		Version: baseGql.NamserviceVersion,
		Node:    baseGql.NodeInfo{Network: r.Keeper.GetChainID()},
		Sync: baseGql.SyncInfo{
			LatestBlockHeight: strconv.FormatInt(statusRecord.LastSyncedHeight, 10),
			CatchingUp:        statusRecord.CatchingUp,
		},
		DiskUsage: diskUsage,
	}, nil
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
