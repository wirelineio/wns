//
// Copyright 2019 Wireline, Inc.
//

package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/Masterminds/semver"
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/tendermint/tendermint/crypto"
	"golang.org/x/crypto/ripemd160"
)

// MarshalMapToJSONBytes converts map[string]interface{} to bytes.
func MarshalMapToJSONBytes(val map[string]interface{}) (bytes []byte) {
	bytes, err := json.Marshal(val)
	if err != nil {
		panic("Marshal error.")
	}

	return
}

// UnMarshalMapFromJSONBytes converts bytes to map[string]interface{}.
func UnMarshalMapFromJSONBytes(bytes []byte) map[string]interface{} {
	var val map[string]interface{}
	err := json.Unmarshal(bytes, &val)

	if err != nil {
		panic("Marshal error.")
	}

	return val
}

// GetCid gets the content ID.
func GetCid(content []byte) string {
	pref := cid.Prefix{
		Version:  0,
		Codec:    cid.DagCBOR,
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}

	cid, err := pref.Sum(content)
	if err != nil {
		panic("CID generation error.")
	}

	return cid.String()
}

// GetAddressFromPubKey gets an address from the public key.
func GetAddressFromPubKey(pubKey crypto.PubKey) string {
	hasherSHA256 := sha256.New()
	hasherSHA256.Write(pubKey.Bytes())
	sha := hasherSHA256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha)
	ripemd := hasherRIPEMD160.Sum(nil)

	return BytesToHex(ripemd)
}

// BytesToBase64 encodes a byte array as a base64 string.
func BytesToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// BytesFromBase64 decodes a byte array from a base64 string.
func BytesFromBase64(str string) []byte {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic("Error decoding string to bytes.")
	}

	return bytes
}

// BytesToHex encodes a byte array as a hex string.
func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// BytesFromHex decodes a byte array from a hex string.
func BytesFromHex(str string) []byte {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		panic("Error decoding hex to bytes.")
	}

	return bytes
}

// Intersection computes the intersection of two string slices.
func Intersection(a, b []string) (c []string) {
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}

	return
}

// GetSemver returns the semver object for a string version.
func GetSemver(versionStr string) *semver.Version {
	version, err := semver.NewVersion(versionStr)
	if err != nil {
		panic("Invalid version string.")
	}

	return version
}
