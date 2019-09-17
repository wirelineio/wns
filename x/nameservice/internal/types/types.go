package types

import (
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MinNamePrice is Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// NewWhois returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}

// ID for records.
type ID string

// Record represents a registry record that can be serialized from/to YAML.
type Record struct {
	ID         ID                     `json:"id"`
	Type       string                 `json:"type"`
	Owner      string                 `json:"owner"`
	Attributes map[string]interface{} `json:"attributes"`
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

// RecordObj represents a registry record.
type RecordObj struct {
	ID         ID     `json:"id"`
	Type       string `json:"type"`
	Owner      string `json:"owner"`
	Attributes []byte `json:"attributes"`
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
	resourceObj.Type = record.Type
	resourceObj.Owner = record.Owner
	resourceObj.Attributes = MarshalMapToJSONBytes(record.Attributes)

	return resourceObj
}

// MarshalLinksToJSONBytes converts []map[string]interface{} to bytes.
func MarshalLinksToJSONBytes(val []map[string]interface{}) (bytes []byte) {
	bytes, err := json.Marshal(val)
	if err != nil {
		panic("Marshal error.")
	}

	return
}

// UnMarshalLinksFromJSONBytes converts bytes to []map[string]interface{}.
func UnMarshalLinksFromJSONBytes(bytes []byte) []map[string]interface{} {
	var val []map[string]interface{}
	err := json.Unmarshal(bytes, &val)

	if err != nil {
		panic("Marshal error.")
	}

	return val
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

// MarshalSliceToJSONBytes converts map[string]interface{} to bytes.
func MarshalSliceToJSONBytes(val []interface{}) (bytes []byte) {
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

// UnMarshalSliceFromJSONBytes converts bytes to map[string]interface{}.
func UnMarshalSliceFromJSONBytes(bytes []byte) []interface{} {
	var val []interface{}
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
	record.Type = resourceObj.Type
	record.Owner = resourceObj.Owner
	record.Attributes = UnMarshalMapFromJSONBytes(resourceObj.Attributes)

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
