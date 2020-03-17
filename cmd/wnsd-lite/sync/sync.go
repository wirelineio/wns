//
// Copyright 2020 Wireline, Inc.
//

package sync

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tendermint/go-amino"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	nameservice "github.com/wirelineio/wns/x/nameservice"
)

// Config represents config for sync functionality.
type Config struct {
	NodeAddress string
	Client      *rpcclient.HTTP
	Codec       *amino.Codec
}

// Synchronize runs a sync cycle.
func Synchronize(config Config, height int64, syncTime time.Time) error {
	fmt.Println("Syncing at height", height, "time", syncTime.UTC())

	cdc := config.Codec

	opts := rpcclient.ABCIQueryOptions{
		Height: height,
		Prove:  true,
	}

	blockHeightKey := nameservice.GetBlockChangesetIndexKey(height)
	res, err := config.Client.ABCIQueryWithOptions("/store/nameservice/key", blockHeightKey, opts)
	if err != nil {
		return err
	}

	// TODO(ashwin): Verify proof.

	var changeset nameservice.BlockChangeset
	cdc.MustUnmarshalBinaryBare(res.Response.Value, &changeset)

	if changeset.Height > 0 {
		fmt.Println(string(cdc.MustMarshalJSON(changeset)))

		for _, id := range changeset.Records {
			recordKey := nameservice.GetRecordIndexKey(id)
			res, err := config.Client.ABCIQueryWithOptions("/store/nameservice/key", recordKey, opts)
			if err != nil {
				return err
			}

			// TODO(ashwin): Verify proof.

			var record nameservice.RecordObj
			cdc.MustUnmarshalBinaryBare(res.Response.Value, &record)

			jsonBytes, _ := json.MarshalIndent(record.ToRecord(), "", "  ")
			fmt.Println(string(jsonBytes))
		}

		for _, name := range changeset.Names {
			nameRecordKey := nameservice.GetNameRecordIndexKey(name)
			res, err := config.Client.ABCIQueryWithOptions("/store/nameservice/key", nameRecordKey, opts)
			if err != nil {
				return err
			}

			// TODO(ashwin): Verify proof.

			var nameRecord nameservice.NameRecord
			cdc.MustUnmarshalBinaryBare(res.Response.Value, &nameRecord)

			jsonBytes, _ := json.MarshalIndent(nameRecord, "", "  ")
			fmt.Println(name, string(jsonBytes))
		}
	}

	return nil
}
