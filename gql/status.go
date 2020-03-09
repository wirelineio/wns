//
// Copyright 2020 Wireline, Inc.
//

package gql

import (
	"strconv"

	"github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
)

func getStatusInfo(ctx *rpctypes.Context) (*NodeInfo, *SyncInfo, *ValidatorInfo, error) {
	res, err := core.Status(ctx)

	if err != nil {
		return nil, nil, nil, err
	}

	nodeInfo := res.NodeInfo
	syncInfo := res.SyncInfo
	valInfo := res.ValidatorInfo

	return &NodeInfo{
			ID:      string(nodeInfo.ID()),
			Moniker: nodeInfo.Moniker,
			Network: nodeInfo.Network,
		}, &SyncInfo{
			LatestBlockHash:   syncInfo.LatestBlockHash.String(),
			LatestBlockHeight: strconv.FormatInt(syncInfo.LatestBlockHeight, 10),
			LatestBlockTime:   syncInfo.LatestBlockTime.UTC().String(),
		}, &ValidatorInfo{
			Address:     valInfo.Address.String(),
			VotingPower: strconv.FormatInt(valInfo.VotingPower, 10),
		}, nil
}
