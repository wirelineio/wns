//
// Copyright 2019 Wireline, Inc.
//

package gql

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/rpc/core"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"github.com/wirelineio/wns/x/nameservice"
)

func decodeStdTx(tx string) (*auth.StdTx, error) {
	bytes, err := base64.StdEncoding.DecodeString(tx)
	if err != nil {
		return nil, err
	}

	// Note: json.Unmarshal doesn't known which Msg struct to use, so we do it "manually".
	// See https://stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(bytes, &objmap)
	if err != nil {
		return nil, err
	}

	var msgsArr []*json.RawMessage
	err = json.Unmarshal(*objmap["msg"], &msgsArr)
	if err != nil {
		return nil, err
	}

	var firstMsg map[string]*json.RawMessage
	err = json.Unmarshal(*msgsArr[0], &firstMsg)
	if err != nil {
		return nil, err
	}

	var messageType string
	err = json.Unmarshal(*firstMsg["type"], &messageType)
	if err != nil {
		return nil, err
	}

	var msgs []sdk.Msg

	switch messageType {
	case "nameservice/SetRecord":
		{
			var setMsg nameservice.MsgSetRecord
			err = json.Unmarshal(*firstMsg["value"], &setMsg)
			if err != nil {
				return nil, err
			}
			msgs = []sdk.Msg{setMsg}
		}
	}

	var fee auth.StdFee
	err = json.Unmarshal(*objmap["fee"], &fee)
	if err != nil {
		return nil, err
	}

	var sigs []*json.RawMessage
	err = json.Unmarshal(*objmap["signatures"], &sigs)
	if err != nil {
		return nil, err
	}

	var sig map[string]*json.RawMessage
	err = json.Unmarshal(*sigs[0], &sig)
	if err != nil {
		return nil, err
	}

	var pubKeyStr string
	err = json.Unmarshal(*sig["pub_key"], &pubKeyStr)
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, err
	}

	pubKey, err := cryptoAmino.PubKeyFromBytes(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	var signature []byte
	err = json.Unmarshal(*sig["signature"], &signature)
	if err != nil {
		return nil, err
	}

	var memo string
	err = json.Unmarshal(*objmap["memo"], &memo)
	if err != nil {
		return nil, err
	}

	stdTx := auth.StdTx{
		Msgs: msgs,
		Fee:  fee,
		Signatures: []auth.StdSignature{auth.StdSignature{
			PubKey:    pubKey,
			Signature: signature,
		}},
		Memo: memo,
	}

	return &stdTx, nil
}

func broadcastTx(r *mutationResolver, stdTx *auth.StdTx) (*ctypes.ResultBroadcastTxCommit, error) {
	txBytes, err := r.Resolver.codec.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		return nil, err
	}

	ctx := &rpctypes.Context{}
	res, err := core.BroadcastTxCommit(ctx, txBytes)
	if err != nil {
		return nil, err
	}

	if res.CheckTx.IsErr() {
		return nil, errors.New(res.CheckTx.String())
	}

	if res.DeliverTx.IsErr() {
		return nil, errors.New(res.DeliverTx.String())
	}

	return res, nil
}
