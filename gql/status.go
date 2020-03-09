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

func getNetInfo(ctx *rpctypes.Context) (string, []*PeerInfo, error) {
	res, err := core.NetInfo(ctx)
	if err != nil {
		return "", nil, err
	}

	peers := res.Peers
	peersInfo := make([]*PeerInfo, len(peers))
	for index, peer := range peers {
		peersInfo[index] = &PeerInfo{
			Node: NodeInfo{
				ID:      string(peer.NodeInfo.ID()),
				Moniker: peer.NodeInfo.Moniker,
				Network: peer.NodeInfo.Network,
			},
			IsOutbound: peer.IsOutbound,
			RemoteIP:   peer.RemoteIP,
		}
	}

	return strconv.FormatInt(int64(res.NPeers), 10), peersInfo, nil
}
