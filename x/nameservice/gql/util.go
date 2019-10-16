//
// Copyright 2019 Wireline, Inc.
//

package gql

import (
	"context"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/wirelineio/wns/x/nameservice/internal/types"
)

func getGQLRecord(ctx context.Context, resolver *queryResolver, record types.Record) (*Record, error) {
	attributes, err := getAttributes(&record)
	if err != nil {
		return nil, err
	}

	extension, err := getExtension(&record)
	if err != nil {
		return nil, err
	}

	references, err := getReferences(ctx, resolver, &record)
	if err != nil {
		return nil, err
	}

	return &Record{
		ID:         string(record.ID),
		Type:       record.Type(),
		Name:       record.Name(),
		Version:    record.Version(),
		Owners:     record.GetOwners(),
		Attributes: attributes,
		References: references,
		Extension:  extension,
	}, nil
}

func getReferences(ctx context.Context, resolver *queryResolver, r *types.Record) ([]*Record, error) {
	var ids []string

	for _, value := range r.Attributes {
		switch value.(type) {
		case interface{}:
			if obj, ok := value.(map[string]interface{}); ok {
				if obj["type"].(string) == WrnTypeReference {
					ids = append(ids, obj["id"].(string))
				}
			}
		}
	}

	return resolver.GetRecordsByIds(ctx, ids)
}

func getAttributes(r *types.Record) ([]*KeyValue, error) {
	return mapToKeyValuePairs(r.Attributes)
}

func getExtension(r *types.Record) (ext Extension, err error) {
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
				if obj["type"].(string) == WrnTypeReference {
					kvPair.Value.Reference = &Reference{
						ID: obj["id"].(string),
					}
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

func matchesOnAttributes(record *types.Record, attributes []*KeyValueInput) bool {
	recAttrs := record.Attributes

	for _, attr := range attributes {
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
			if !ok || *attr.Value.String != recAttrValString {
				return false
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
