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
func GenerateHash(json map[string]interface{}) (string, string, error) {
	content, err := canonicalJson.Marshal(json)
	if err != nil {
		return "", "", err
	}

	cid, err := cbor.FromJSON(bytes.NewReader(content), mh.SHA2_256, -1)
	if err != nil {
		return "", "", err
	}

	return cid.String(), string(content), nil
}
