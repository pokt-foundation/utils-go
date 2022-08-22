package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseIntegerFromPayload parses a string (or integer if any) value as an int from a JSON payload
// Keys must be given in the form: "key1.nestedKey1.nestedKey2"
// Should this be refactored to also support strings?
func ParseIntegerFromPayload(r io.Reader, key string) (int64, error) {
	// TODO: Parse nested fields
	res := map[string]any{} // got'em 'any' haters

	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return 0, errors.New("error decoding payload: " + err.Error())
	}

	value, err := nestedMapLookup(res, key)
	if err != nil {
		return 0, err
	}

	var valueDecimal int64
	switch v := value.(type) {
	case string:
		valueDecimal, err = strconv.ParseInt(v, 0, 64)
	case int64:
		valueDecimal = int64(v)
	case float64:
		valueDecimal = int64(v)
	default:
		err = fmt.Errorf("invalid type for payload: %T", v)
	}

	if err != nil {
		return 0, fmt.Errorf("error parsing field %s: %s", key, err.Error())
	}

	return valueDecimal, nil
}

// ParseIntegerJSONString parses a string (or integer if any) value as an int from a JSON string.
// Keys must be given in the form: "key1.nestedKey1.nestedKey2"
// Should this be refactored to also support strings?
func ParseIntegerJSONString(r string, key string) (int64, error) {
	// TODO: Parse nested fields
	res := map[string]any{} // got'em 'any' haters

	err := json.Unmarshal([]byte(r), &res)
	if err != nil {
		return 0, errors.New("error decoding payload: " + err.Error())
	}

	value, err := nestedMapLookup(res, key)
	if err != nil {
		return 0, err
	}

	var valueDecimal int64
	switch v := value.(type) {
	case string:
		valueDecimal, err = strconv.ParseInt(v, 0, 64)
	case int64:
		valueDecimal = int64(v)
	case float64:
		valueDecimal = int64(v)
	default:
		err = fmt.Errorf("invalid type for payload: %T", v)
	}

	if err != nil {
		return 0, fmt.Errorf("error parsing field %s: %s", key, err.Error())
	}

	return valueDecimal, nil
}

// nestedMapLookup returns the value of a map with an arbitrary number of nested
// fields. Keys must be given in the form: "key1.nestedKey1.nestedKey2"
func nestedMapLookup(m map[string]any, keys string) (any, error) {
	ks := strings.Split(keys, ".")
	if len(ks) == 0 {
		return nil, fmt.Errorf("no keys given to find on map")
	}

	value, ok := m[ks[0]]
	if !ok {
		return nil, fmt.Errorf("key %s not found", ks[0])
	}
	if len(ks) == 1 {
		return value, nil
	}
	nestedMap, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("nested key is not of type map[string]any")
	}

	return nestedMapLookup(nestedMap, strings.Join(ks[1:], "."))
}
