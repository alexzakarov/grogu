package utils

import (
	"reflect"
	"strings"
	"sync"
)

type FieldMetadata struct {
	FieldName string
	JsonName  string
	Index     int
}

type StructMetadata struct {
	Fields       []FieldMetadata
	JsonNames    []string
	FieldIndices []int
}

type MetadataCache struct {
	cache sync.Map
}

var globalCache = &MetadataCache{}

func (mc *MetadataCache) get(structType reflect.Type) (*StructMetadata, bool) {
	if value, ok := mc.cache.Load(structType); ok {
		return value.(*StructMetadata), true
	}
	return nil, false
}

func (mc *MetadataCache) set(structType reflect.Type, metadata *StructMetadata) {
	mc.cache.Store(structType, metadata)
}

func buildMetadata(structType reflect.Type) *StructMetadata {
	var fields []FieldMetadata
	var jsonNames []string
	var fieldIndices []int

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		jsonTag := field.Tag.Get("json")

		if strings.Contains(jsonTag, "sub") {
			continue
		}

		switch jsonTag {
		case "-", "":
			continue
		default:
			parts := strings.Split(jsonTag, ",")
			name := parts[0]
			if name == "" {
				name = jsonTag
			}

			fields = append(fields, FieldMetadata{
				FieldName: field.Name,
				JsonName:  name,
				Index:     i,
			})
			jsonNames = append(jsonNames, name)
			fieldIndices = append(fieldIndices, i)
		}
	}

	return &StructMetadata{
		Fields:       fields,
		JsonNames:    jsonNames,
		FieldIndices: fieldIndices,
	}
}

func getOrBuildMetadata(structType reflect.Type) *StructMetadata {
	if metadata, found := globalCache.get(structType); found {
		return metadata
	}

	metadata := buildMetadata(structType)
	globalCache.set(structType, metadata)
	return metadata
}
