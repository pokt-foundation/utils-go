// Package environment is a set of functions to get env values or their default
// It has the autoload from .env files
package environment

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	// autoload env vars
	_ "github.com/joho/godotenv/autoload"
)

// MustGetInt64 gets the required env var as an int and panics if it is not present
func MustGetInt64(varName string) int64 {
	val, ok := os.LookupEnv(varName)
	if !ok {
		panic(fmt.Sprintf("environment error (int64): required env var %s not found", varName))
	}

	iVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("environment error (int64): unable to parse value %s to required env var %s", val, varName))
	}

	return iVal
}

// GetInt64 gets the env var as an int
func GetInt64(varName string, defaultValue int64) int64 {
	val, ok := os.LookupEnv(varName)
	if !ok {
		return defaultValue
	}

	iVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultValue
	}

	return iVal
}

// GetFloat64 gets the env var as a float64
func GetFloat64(varName string, defaultValue float64) float64 {
	val, ok := os.LookupEnv(varName)
	if !ok {
		return defaultValue
	}

	iVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultValue
	}

	return iVal
}

// MuastGetFloat64 gets the required env var as a float64 and panics if is not present
func MustGetFloat64(varName string) float64 {
	val, ok := os.LookupEnv(varName)
	if !ok {
		panic(fmt.Sprintf("environment error (float64): required env var %s not found", varName))
	}

	iVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(fmt.Sprintf("environment error (float64): unable to parse: %v. env name: %s", err, varName))
	}

	return iVal
}

// MustGetString gets the required environment var as a string and panics if it is not present
func MustGetString(varName string) string {
	val, _ := os.LookupEnv(varName)
	if val == "" {
		panic(fmt.Sprintf("environment error (string): required env var %s not found", varName))
	}

	return val
}

// GetString gets the environment var as a string
func GetString(varName string, defaultValue string) string {
	val, _ := os.LookupEnv(varName)
	if val == "" {
		return defaultValue
	}

	return val
}

// MustGetBool gets the required environment var as a bool and panics if it is not present
func MustGetBool(varName string) bool {
	val, _ := os.LookupEnv(varName)
	if val == "" {
		panic(fmt.Sprintf("environment error (bool): required env var %s not found", varName))
	}

	iVal, err := strconv.ParseBool(val)
	if err != nil {
		panic(fmt.Sprintf("environment error (bool): unable to parse value %s to required env var %s", val, varName))
	}

	return iVal
}

// GetBool gets the environment var as a bool
func GetBool(varName string, defaultValue bool) bool {
	val, _ := os.LookupEnv(varName)
	if val == "" {
		return defaultValue
	}

	iVal, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return iVal
}

// MustGetStringMap gets the required environment var as a string map and panics if it is not present
func MustGetStringMap(varName, separator string) map[string]bool {
	rawString := MustGetString(varName)
	stringSlice := strings.Split(rawString, separator)

	stringMap := make(map[string]bool, len(stringSlice))

	for _, singleString := range stringSlice {
		stringMap[singleString] = true
	}

	return stringMap
}

// GetStringMap gets the environment var as a string map
func GetStringMap(varName, defaultValue, separator string) map[string]bool {
	rawString := GetString(varName, defaultValue)
	stringSlice := strings.Split(rawString, separator)

	stringMap := make(map[string]bool, len(stringSlice))

	for _, singleString := range stringSlice {
		stringMap[singleString] = true
	}

	return stringMap
}
