//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	"crypto/sha256"
	"fmt"
	"strings"

	canonicalJson "github.com/gibson042/canonicaljson-go"
	"github.com/wirelineio/wns/x/nameservice/internal/helpers"
)

// ID for records.
type ID string

// Record represents a WNS record.
type Record struct {
	ID         ID                     `json:"id,omitempty"`
	Owners     []string               `json:"owners,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
}

// Type of Record.
func (r Record) Type() string {
	return r.Attributes["type"].(string)
}

// Name of Record.
func (r Record) Name() string {
	return r.Attributes["name"].(string)
}

// Version of Record.
func (r Record) Version() string {
	return r.Attributes["version"].(string)
}

// WRN returns the record `wrn`, e.g. `wrn:bot:wireline.io/chess#0.1.0`.
func (r Record) WRN() string {
	return strings.ToLower(fmt.Sprintf("%s#%s", r.BaseWRN(), r.Version()))
}

// BaseWRN returns the record `wrn` minus the version, e.g. `wrn:bot:wireline.io/chess`.
func (r Record) BaseWRN() string {
	return strings.ToLower(fmt.Sprintf("%s:%s", r.Type(), r.Name()))
}

// GetOwners returns the list of owners (for GQL).
func (r Record) GetOwners() []*string {
	owners := []*string{}
	for _, owner := range r.Owners {
		owners = append(owners, &owner)
	}

	return owners
}

// ToRecordObj convers Record to RecordObj.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (r *Record) ToRecordObj() RecordObj {
	var resourceObj RecordObj

	resourceObj.ID = r.ID
	resourceObj.Owners = r.Owners
	resourceObj.Attributes = helpers.MarshalMapToJSONBytes(r.Attributes)

	return resourceObj
}

// ToNameRecord gets a naming record entry for the record.
func (r *Record) ToNameRecord() NameRecord {
	var nameRecord NameRecord
	nameRecord.ID = r.ID
	nameRecord.Version = r.Version()

	return nameRecord
}

// CanonicalJSON returns the canonical JSON respresentation of the record.
func (r *Record) CanonicalJSON() []byte {
	bytes, err := canonicalJson.Marshal(r.Attributes)
	if err != nil {
		panic("Record marshal error.")
	}

	return bytes
}

// GetSignBytes generates a record hash to be signed.
func (r *Record) GetSignBytes() ([]byte, []byte) {
	// Double SHA256 hash.

	// Input to the first round of hashing.
	bytes := r.CanonicalJSON()

	// First round.
	first := sha256.New()
	first.Write(bytes)
	firstHash := first.Sum(nil)

	// Second round of hashing takes as input the output of the first round.
	second := sha256.New()
	second.Write(firstHash)
	secondHash := second.Sum(nil)

	return secondHash, bytes
}

// GetCID gets the record CID.
func (r *Record) GetCID() ID {
	return ID(helpers.GetCid(r.CanonicalJSON()))
}

// Signature represents a record signature.
type Signature struct {
	PubKey    string `json:"pubKey"`
	Signature string `json:"sig"`
}

// PayloadObj represents a signed record payload.
type PayloadObj struct {
	Record     RecordObj   `json:"record"`
	Signatures []Signature `json:"signatures"`
}

// ToPayload converts Payload to PayloadObj object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (payloadObj PayloadObj) ToPayload() Payload {
	var payload Payload

	payload.Record = helpers.UnMarshalMapFromJSONBytes(payloadObj.Record.Attributes)
	payload.Signatures = payloadObj.Signatures

	return payload
}

// RecordObj represents a WNS record.
type RecordObj struct {
	ID         ID       `json:"id,omitempty"`
	Owners     []string `json:"owners,omitempty"`
	Attributes []byte   `json:"attributes"`
}

// ToRecord convers RecordObj to Record.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (resourceObj *RecordObj) ToRecord() Record {
	var record Record

	record.ID = resourceObj.ID
	record.Owners = resourceObj.Owners
	record.Attributes = helpers.UnMarshalMapFromJSONBytes(resourceObj.Attributes)

	return record
}

// Payload represents a signed record payload that can be serialized from/to YAML.
type Payload struct {
	Record     map[string]interface{} `json:"record"`
	Signatures []Signature            `json:"signatures"`
}

// ToPayloadObj converts Payload to PayloadObj object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (payload *Payload) ToPayloadObj() PayloadObj {
	var payloadObj PayloadObj

	payloadObj.Record.Attributes = helpers.MarshalMapToJSONBytes(payload.Record)
	payloadObj.Signatures = payload.Signatures

	return payloadObj
}

// NameRecord is a naming record entry for a WRN.
type NameRecord struct {
	ID      ID     `json:"id"`
	Version string `json:"version"`
}
