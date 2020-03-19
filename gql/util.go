//
// Copyright 2019 Wireline, Inc.
//

package gql

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/Masterminds/semver"
	"github.com/mitchellh/mapstructure"
	"github.com/wirelineio/wns/x/bond"
	"github.com/wirelineio/wns/x/nameservice"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VersionAttributeName denotes the version attribute name in a record.
const VersionAttributeName = "version"

// VersionMatchAll represents a special value to match all versions.
const VersionMatchAll = "*"

// VersionMatchLatest represents a special value to match only the latest version of each record.
const VersionMatchLatest = "latest"

// OwnerAttributeName denotes the owner attribute name for a bond.
const OwnerAttributeName = "owner"

// BondIDAttributeName denotes the record bond ID.
const BondIDAttributeName = "bondId"

// ExpiryTimeAttributeName denotes the record expiry time.
const ExpiryTimeAttributeName = "expiryTime"

func GetGQLRecord(ctx context.Context, resolver QueryResolver, record *nameservice.Record) (*Record, error) {
	// Nil record.
	if record == nil || record.Deleted {
		return nil, nil
	}

	attributes, err := getAttributes(record)
	if err != nil {
		return nil, err
	}

	extension, err := getExtension(record)
	if err != nil {
		return nil, err
	}

	references, err := getReferences(ctx, resolver, record)
	if err != nil {
		return nil, err
	}

	return &Record{
		ID:         string(record.ID),
		Type:       record.Type(),
		Name:       record.Name(),
		Version:    record.Version(),
		BondID:     record.GetBondID(),
		CreateTime: record.GetCreateTime(),
		ExpiryTime: record.GetExpiryTime(),
		Owners:     record.GetOwners(),
		Attributes: attributes,
		References: references,
		Extension:  extension,
	}, nil
}

func getReferences(ctx context.Context, resolver QueryResolver, r *nameservice.Record) ([]*Record, error) {
	var ids []string

	for _, value := range r.Attributes {
		switch value.(type) {
		case interface{}:
			if obj, ok := value.(map[string]interface{}); ok {
				if typeAttr, ok := obj["type"]; ok && typeAttr.(string) == WrnTypeReference {
					ids = append(ids, obj["id"].(string))
				}
			}
		}
	}

	return resolver.GetRecordsByIds(ctx, ids)
}

func getAttributes(r *nameservice.Record) ([]*KeyValue, error) {
	return mapToKeyValuePairs(r.Attributes)
}

func getExtension(r *nameservice.Record) (ext Extension, err error) {
	switch r.Type() {
	case WnsTypeProtocol:
		var protocol Protocol
		err := mapstructure.Decode(r.Attributes, &protocol)
		return protocol, err
	case WnsTypeBot:
		var bot Bot
		err := mapstructure.Decode(r.Attributes, &bot)
		return bot, err
	case WnsTypePad:
		var pad Pad
		err := mapstructure.Decode(r.Attributes, &pad)
		return pad, err
	default:
		var unknown UnknownExtension
		err := mapstructure.Decode(r.Attributes, &unknown)
		return unknown, err
	}
}

func mapToKeyValuePairs(attrs map[string]interface{}) ([]*KeyValue, error) {
	kvPairs := []*KeyValue{}

	trueVal := true
	falseVal := false

	for key, value := range attrs {

		kvPair := &KeyValue{
			Key: key,
		}

		switch val := value.(type) {
		case nil:
			kvPair.Value.Null = &trueVal
		case int:
			kvPair.Value.Int = &val
		case float64:
			kvPair.Value.Float = &val
		case string:
			kvPair.Value.String = &val
		case bool:
			kvPair.Value.Boolean = &val
		case interface{}:
			if obj, ok := value.(map[string]interface{}); ok {
				if valueType, ok := obj["type"]; ok && valueType.(string) == WrnTypeReference {
					kvPair.Value.Reference = &Reference{
						ID: obj["id"].(string),
					}
				} else {
					bytes, err := json.Marshal(obj)
					if err != nil {
						return nil, err
					}

					jsonStr := string(bytes)
					kvPair.Value.String = &jsonStr
				}
			}
		}

		if kvPair.Value.Null == nil {
			kvPair.Value.Null = &falseVal
		}

		valueType := reflect.ValueOf(value)
		if valueType.Kind() == reflect.Slice {
			// TODO(ashwin): Handle arrays.
		}

		kvPairs = append(kvPairs, kvPair)
	}

	return kvPairs, nil
}

func matchOnRecordField(record *nameservice.Record, attr *KeyValueInput) (fieldFound bool, matched bool) {
	fieldFound = false
	matched = true

	switch attr.Key {
	case BondIDAttributeName:
		{
			fieldFound = true
			if attr.Value.String == nil || record.GetBondID() != *attr.Value.String {
				matched = false
				return
			}
		}
	case ExpiryTimeAttributeName:
		{
			fieldFound = true
			if attr.Value.String == nil || record.GetExpiryTime() != *attr.Value.String {
				matched = false
				return
			}
		}
	}

	return
}

func MatchOnAttributes(record *nameservice.Record, attributes []*KeyValueInput) bool {
	// Filter deleted records.
	if record.Deleted {
		return false
	}

	recAttrs := record.Attributes

	for _, attr := range attributes {
		// First try matching on record struct fields.
		fieldFound, matched := matchOnRecordField(record, attr)
		if fieldFound {
			if !matched {
				return false
			}

			continue
		}

		recAttrVal, recAttrFound := recAttrs[attr.Key]
		if !recAttrFound {
			return false
		}

		if attr.Value.Int != nil {
			recAttrValInt, ok := recAttrVal.(int)
			if !ok || *attr.Value.Int != recAttrValInt {
				return false
			}
		}

		if attr.Value.Float != nil {
			recAttrValFloat, ok := recAttrVal.(float64)
			if !ok || *attr.Value.Float != recAttrValFloat {
				return false
			}
		}

		if attr.Value.String != nil {
			recAttrValString, ok := recAttrVal.(string)
			if !ok {
				return false
			}

			// Special handling for version attribute.
			if attr.Key == VersionAttributeName {
				if !matchOnVersionAttribute(*attr.Value.String, recAttrValString) {
					return false
				}
			} else {
				if *attr.Value.String != recAttrValString {
					return false
				}
			}
		}

		if attr.Value.Boolean != nil {
			recAttrValBool, ok := recAttrVal.(bool)
			if !ok || *attr.Value.Boolean != recAttrValBool {
				return false
			}
		}

		if attr.Value.Reference != nil {
			obj, ok := recAttrVal.(map[string]interface{})
			if !ok || obj["type"].(string) != WrnTypeReference {
				return false
			}
			recAttrValRefID := obj["id"].(string)
			if recAttrValRefID != attr.Value.Reference.ID {
				return false
			}
		}

		// TODO(ashwin): Handle arrays.
	}

	return true
}

func matchOnVersionAttribute(querySemverStr string, recordVersionStr string) bool {
	if querySemverStr == VersionMatchAll || querySemverStr == VersionMatchLatest {
		return true
	}

	querySemverConstraint, err := semver.NewConstraint(querySemverStr)
	if err != nil {
		// Handle constraint not being parsable.
		return false
	}

	recordVersion, err := semver.NewVersion(recordVersionStr)
	if err != nil {
		return false
	}

	return querySemverConstraint.Check(recordVersion)
}

func RequestedLatestVersionsOnly(attributes []*KeyValueInput) bool {
	for _, attr := range attributes {
		if attr.Key == VersionAttributeName && attr.Value.String != nil {
			if *attr.Value.String == VersionMatchAll {
				return false
			}

			if *attr.Value.String == VersionMatchLatest {
				return true
			}
		}
	}

	return true
}

// Used to filter records and retain only the latest versions.
type bestMatch struct {
	version *semver.Version
	record  *nameservice.Record
}

// Only return the latest version of each record.
func GetLatestVersions(records []*nameservice.Record) []*nameservice.Record {
	baseWrnBestMatch := make(map[string]bestMatch)
	for _, record := range records {
		baseWrn := record.BaseWRN()
		recordVersion, _ := semver.NewVersion(record.Version())

		currentBestMatch, exists := baseWrnBestMatch[baseWrn]
		if !exists || recordVersion.GreaterThan(currentBestMatch.version) {
			// Update current best match.
			baseWrnBestMatch[baseWrn] = bestMatch{recordVersion, record}
		}
	}

	var matches = make([]*nameservice.Record, len(baseWrnBestMatch))
	var index int
	for _, match := range baseWrnBestMatch {
		matches[index] = match.record
		index++
	}

	return matches
}

func getGQLCoins(coins sdk.Coins) []Coin {
	gqlCoins := make([]Coin, len(coins))
	for index, coin := range coins {
		gqlCoins[index] = Coin{
			Type:     coin.Denom,
			Quantity: BigUInt(coin.Amount.Int64()),
		}
	}

	return gqlCoins
}

func getGQLBond(ctx context.Context, resolver *queryResolver, bondObj *bond.Bond) (*Bond, error) {
	// Nil record.
	if bondObj == nil {
		return nil, nil
	}

	return &Bond{
		ID:      string(bondObj.ID),
		Owner:   bondObj.Owner,
		Balance: getGQLCoins(bondObj.Balance),
	}, nil
}

func matchBondOnAttributes(bondObj *bond.Bond, attributes []*KeyValueInput) bool {
	for _, attr := range attributes {
		switch attr.Key {
		case OwnerAttributeName:
			{
				if attr.Value.String == nil || bondObj.Owner != *attr.Value.String {
					return false
				}
			}
		default:
			{
				// Only attributes explicitly listed in the switch are queryable.
				return false
			}
		}
	}

	return true
}
