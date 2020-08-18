//
// Copyright 2020 Wireline, Inc.
//

package utils

import (
	"bytes"

	canonicalJson "github.com/gibson042/canonicaljson-go"
	cbor "github.com/ipfs/go-ipld-cbor"
	mh "github.com/multiformats/go-multihash"
)

// GenerateHash returns the hash of the canonicalized JSON input.
func GenerateHash(json map[string]interface{}) (string, []byte, error) {
	content, err := canonicalJson.Marshal(json)
	if err != nil {
		return "", nil, err
	}

	cid, err := cbor.FromJSON(bytes.NewReader(content), mh.SHA2_256, -1)
	if err != nil {
		return "", nil, err
	}

	return cid.String(), content, nil
}
