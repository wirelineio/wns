package types

import (
	"encoding/json"
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

// RecordObj represents a WNS record.
type RecordObj struct {
	ID         ID       `json:"id,omitempty"`
	Owners     []string `json:"owners,omitempty"`
	Attributes []byte   `json:"attributes"`
	Extension  []byte   `json:"extension"`
}

// Payload represents a signed record payload that can be serialized from/to YAML.
type Payload struct {
	Record     Record      `json:"record"`
	Signatures []Signature `json:"signatures"`
}

// RecordToRecordObj convers Record to RecordObj.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func RecordToRecordObj(record Record) RecordObj {
	var resourceObj RecordObj

	resourceObj.ID = record.ID
	resourceObj.Owners = record.Owners
	resourceObj.Attributes = MarshalMapToJSONBytes(record.Attributes)
	resourceObj.Extension = MarshalMapToJSONBytes(record.Extension)

	return resourceObj
}

// PayloadToPayloadObj converts Payload to PayloadObj object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func PayloadToPayloadObj(payload Payload) PayloadObj {
	var payloadObj PayloadObj

	payloadObj.Record = RecordToRecordObj(payload.Record)
	payloadObj.Signatures = payload.Signatures

	return payloadObj
}

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

// RecordObjToRecord convers RecordObj to Record.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func RecordObjToRecord(resourceObj RecordObj) Record {
	var record Record

	record.ID = resourceObj.ID
	record.Owners = resourceObj.Owners
	record.Attributes = UnMarshalMapFromJSONBytes(resourceObj.Attributes)
	record.Extension = UnMarshalMapFromJSONBytes(resourceObj.Extension)

	return record
}

// PayloadObjToPayload converts Payload to PayloadObj object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func PayloadObjToPayload(payloadObj PayloadObj) Payload {
	var payload Payload

	payload.Record = RecordObjToRecord(payloadObj.Record)
	payload.Signatures = payloadObj.Signatures

	return payload
}
