//
// Copyright 2019 Wireline, Inc.
//

package types

import (
	"crypto/sha256"

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
	Extension  map[string]interface{} `json:"extension"`
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
	resourceObj.Extension = helpers.MarshalMapToJSONBytes(r.Extension)

	return resourceObj
}

// CanonicalJSON returns the canonical JSON respresentation of the record.
func (r *Record) CanonicalJSON() []byte {
	record := Record{
		Attributes: r.Attributes,
		Extension:  r.Extension,
	}

	bytes, err := canonicalJson.Marshal(record)
	if err != nil {
		panic("Record marshal error.")
	}

	return bytes
}

// GetSignBytes generates a transaction hash.
func (r *Record) GetSignBytes() []byte {
	// Double SHA256 hash.

	// First round.
	first := sha256.New()
	bytes := r.CanonicalJSON()

	first.Write(bytes)
	firstHash := first.Sum(nil)

	// Second round.
	second := sha256.New()
	second.Write(firstHash)
	secondHash := second.Sum(nil)

	return secondHash
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

	payload.Record = payloadObj.Record.ToRecord()
	payload.Signatures = payloadObj.Signatures

	return payload
}

// RecordObj represents a WNS record.
type RecordObj struct {
	ID         ID       `json:"id,omitempty"`
	Owners     []string `json:"owners,omitempty"`
	Attributes []byte   `json:"attributes"`
	Extension  []byte   `json:"extension"`
}

// ToRecord convers RecordObj to Record.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (resourceObj *RecordObj) ToRecord() Record {
	var record Record

	record.ID = resourceObj.ID
	record.Owners = resourceObj.Owners
	record.Attributes = helpers.UnMarshalMapFromJSONBytes(resourceObj.Attributes)
	record.Extension = helpers.UnMarshalMapFromJSONBytes(resourceObj.Extension)

	return record
}

// Payload represents a signed record payload that can be serialized from/to YAML.
type Payload struct {
	Record     Record      `json:"record"`
	Signatures []Signature `json:"signatures"`
}

// ToPayloadObj converts Payload to PayloadObj object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func (payload *Payload) ToPayloadObj() PayloadObj {
	var payloadObj PayloadObj

	payloadObj.Record = payload.Record.ToRecordObj()
	payloadObj.Signatures = payload.Signatures

	return payloadObj
}
