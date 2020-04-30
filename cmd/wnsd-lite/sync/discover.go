//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"context"

	"github.com/machinebox/graphql"
	"github.com/wirelineio/wns/gql"
)

const query = `
{
	records: queryRecords(attributes: [{ key: "type", value: { string: "wrn:xbox" } }]) {
		id
		type
		version
		attributes {
			key
			value {
				string
			}
		}
	}
}
`

// Response represents the GQL response.
type Response struct {
	Records []gql.Record `json:"records"`
}

// DiscoverRPCEndpoints queries for WNS RPC endpoints.
func DiscoverRPCEndpoints(ctx *Context, gqlEndpoint string) ([]string, error) {
	client := graphql.NewClient(gqlEndpoint)
	req := graphql.NewRequest(query)
	req.Header.Set("Cache-Control", "no-cache")
	gqlContext := context.Background()

	var response Response
	if err := client.Run(gqlContext, req, &response); err != nil {
		ctx.log.Errorln(err)
		return nil, err
	}

	rpcEndpoints := []string{}

	for _, record := range response.Records {
		for _, kv := range record.Attributes {
			if kv.Key == "wns" {
				var wnsData map[string]interface{}
				if server, ok := wnsData["rpc"].(string); ok {
					rpcEndpoints = append(rpcEndpoints, server)
				}
			}
		}
	}

	return rpcEndpoints, nil
}
