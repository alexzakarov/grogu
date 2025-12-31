package utils

import (
	"errors"
	"reflect"
	"strings"
)

func getStructType(struc interface{}) reflect.Type {
	sType := reflect.TypeOf(struc)
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}

	return sType
}

func Convert(struc interface{}) ([]string, []interface{}, error) {
	var returnMap []interface{}
	var returnJsons []string

	sType := getStructType(struc)

	if sType.Kind() != reflect.Struct {
		return returnJsons, returnMap, errors.New("variable given is not a struct or a pointer to a struct")
	}

	metadata := getOrBuildMetadata(sType)

	structVal := reflect.ValueOf(struc)
	if structVal.Kind() == reflect.Ptr {
		structVal = structVal.Elem()
	}

	for _, fieldMeta := range metadata.Fields {
		fieldValue := structVal.Field(fieldMeta.Index)
		returnMap = append(returnMap, fieldValue.Interface())
		returnJsons = append(returnJsons, fieldMeta.JsonName)
	}

	return returnJsons, returnMap, nil
}

// CheckStringIfContains check a string if contains given param
func CheckStringIfContains(input_text string, search_text string) bool {
	CheckContains := strings.Contains(input_text, search_text)
	return CheckContains
}

func ConvertStatus(status int64, status_type string) interface{} {
	var data interface{}
	if status_type == "bool" {
		switch status {
		case 1:
			data = true
		case 2:
			data = false
		}
	} else if status_type == "int" {
		data = status
	}
	return data
}
