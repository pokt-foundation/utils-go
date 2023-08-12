package environment

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MustGetInt64(t *testing.T) {
	tests := []struct {
		name          string
		envVarName    string
		envVarValue   string
		shouldPanic   bool
		expectedValue int64
	}{
		{
			name:          "Should return correct int64 value",
			envVarName:    "TEST_INT64",
			envVarValue:   "42",
			shouldPanic:   false,
			expectedValue: 42,
		},
		{
			name:          "Should panic on invalid int64 value",
			envVarName:    "TEST_INT64",
			envVarValue:   "invalid",
			shouldPanic:   true,
			expectedValue: 0,
		},
		{
			name:          "Should panic on missing int64 value",
			envVarName:    "MISSING_INT64",
			envVarValue:   "",
			shouldPanic:   true,
			expectedValue: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.envVarName != "" {
				os.Setenv(test.envVarName, test.envVarValue)
			} else {
				os.Unsetenv(test.envVarName)
			}

			if test.shouldPanic {
				assert.Panics(t, func() {
					MustGetInt64(test.envVarName)
				})
			} else {
				result := MustGetInt64(test.envVarName)
				assert.Equal(t, test.expectedValue, result)
			}
		})
	}
}

func Test_GetInt64(t *testing.T) {
	tests := []struct {
		name          string
		envVarName    string
		envVarValue   string
		defaultValue  int64
		expectedValue int64
	}{
		{
			name:          "Should return correct int64 value",
			envVarName:    "TEST_INT64",
			envVarValue:   "42",
			defaultValue:  0,
			expectedValue: 42,
		},
		{
			name:          "Should return default int64 value on invalid value",
			envVarName:    "TEST_INT64",
			envVarValue:   "invalid",
			defaultValue:  0,
			expectedValue: 0,
		},
		{
			name:          "Should return default int64 value on missing value",
			envVarName:    "MISSING_INT64",
			envVarValue:   "",
			defaultValue:  0,
			expectedValue: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envVarName, test.envVarValue)
			result := GetInt64(test.envVarName, test.defaultValue)
			assert.Equal(t, test.expectedValue, result)
		})
	}
}

func Test_MustGetFloat64(t *testing.T) {
	tests := []struct {
		name        string
		envVarName  string
		envVarValue string
		shouldPanic bool
	}{
		{
			name:        "Should return correct float64 value",
			envVarName:  "TEST_FLOAT64",
			envVarValue: "42.2",
			shouldPanic: false,
		},
		{
			name:        "Should panic on invalid float64 value",
			envVarName:  "TEST_FLOAT64",
			envVarValue: "invalid",
			shouldPanic: true,
		},
		{
			name:        "Should panic on missing float64 value",
			envVarName:  "MISSING_FLOAT64",
			envVarValue: "",
			shouldPanic: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envVarName, test.envVarValue)

			if test.shouldPanic {
				assert.Panics(t, func() {
					MustGetFloat64(test.envVarName)
				})
			} else {
				result := MustGetFloat64(test.envVarName)
				expectedValue, _ := strconv.ParseFloat(test.envVarValue, 64)
				assert.Equal(t, expectedValue, result)
			}
		})
	}
}

func Test_GetFloat64(t *testing.T) {
	tests := []struct {
		name          string
		envVarName    string
		envVarValue   string
		defaultValue  float64
		expectedValue float64
	}{
		{
			name:          "Should return correct float64 value",
			envVarName:    "TEST_FLOAT64",
			envVarValue:   "42.2",
			defaultValue:  0,
			expectedValue: 42.2,
		},
		{
			name:          "Should return default float64 value on invalid value",
			envVarName:    "TEST_FLOAT64",
			envVarValue:   "invalid",
			defaultValue:  0,
			expectedValue: 0,
		},
		{
			name:          "Should return default float64 value on missing value",
			envVarName:    "MISSING_FLOAT64",
			envVarValue:   "",
			defaultValue:  0,
			expectedValue: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envVarName, test.envVarValue)
			result := GetFloat64(test.envVarName, test.defaultValue)
			assert.Equal(t, test.expectedValue, result)
		})
	}
}

func Test_MustGetString(t *testing.T) {
	tests := []struct {
		name        string
		envVarName  string
		envVarValue string
		shouldPanic bool
	}{
		{
			name:        "Should return correct string value",
			envVarName:  "TEST_STRING",
			envVarValue: "value",
			shouldPanic: false,
		},
		{
			name:        "Should panic on missing string value",
			envVarName:  "MISSING_STRING",
			envVarValue: "",
			shouldPanic: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envVarName, test.envVarValue)

			if test.shouldPanic {
				assert.Panics(t, func() {
					MustGetString(test.envVarName)
				})
			} else {
				result := MustGetString(test.envVarName)
				assert.Equal(t, test.envVarValue, result)
			}
		})
	}
}

func Test_GetString(t *testing.T) {
	tests := []struct {
		name          string
		envVarName    string
		envVarValue   string
		defaultValue  string
		expectedValue string
	}{
		{
			name:          "Should return correct string value",
			envVarName:    "TEST_STRING",
			envVarValue:   "value",
			defaultValue:  "default",
			expectedValue: "value",
		},
		{
			name:          "Should return default string value when missing",
			envVarName:    "MISSING_STRING",
			envVarValue:   "",
			defaultValue:  "default",
			expectedValue: "default",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envVarName, test.envVarValue)
			result := GetString(test.envVarName, test.defaultValue)
			assert.Equal(t, test.expectedValue, result)
		})
	}
}

func Test_MustGetBool(t *testing.T) {
	tests := []struct {
		name        string
		envVarName  string
		envVarValue string
		shouldPanic bool
	}{
		{
			name:        "Should return correct bool value",
			envVarName:  "TEST_BOOL",
			envVarValue: "true",
			shouldPanic: false,
		},
		{
			name:        "Should panic on invalid bool value",
			envVarName:  "TEST_BOOL",
			envVarValue: "invalid",
			shouldPanic: true,
		},
		{
			name:        "Should panic on missing bool value",
			envVarName:  "MISSING_BOOL",
			envVarValue: "",
			shouldPanic: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envVarName, test.envVarValue)

			if test.shouldPanic {
				assert.Panics(t, func() {
					MustGetBool(test.envVarName)
				})
			} else {
				result := MustGetBool(test.envVarName)
				expectedValue, _ := strconv.ParseBool(test.envVarValue)
				assert.Equal(t, expectedValue, result)
			}
		})
	}
}

func Test_GetBool(t *testing.T) {
	tests := []struct {
		name          string
		envVarName    string
		envVarValue   string
		defaultValue  bool
		expectedValue bool
	}{
		{
			name:          "Should return correct bool value",
			envVarName:    "TEST_BOOL",
			envVarValue:   "true",
			defaultValue:  false,
			expectedValue: true,
		},
		{
			name:          "Should return default bool value when invalid",
			envVarName:    "TEST_BOOL",
			envVarValue:   "invalid",
			defaultValue:  false,
			expectedValue: false,
		},
		{
			name:          "Should return default bool value when missing",
			envVarName:    "MISSING_BOOL",
			envVarValue:   "",
			defaultValue:  false,
			expectedValue: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envVarName, test.envVarValue)
			result := GetBool(test.envVarName, test.defaultValue)
			assert.Equal(t, test.expectedValue, result)
		})
	}
}

func Test_MustGetStringMap(t *testing.T) {
	tests := []struct {
		name          string
		envVarName    string
		separator     string
		envVarValue   string
		expectedValue map[string]bool
	}{
		{
			name:          "Should return correct string map",
			envVarName:    "TEST_STRING_MAP",
			separator:     ",",
			envVarValue:   "key1,key2,key3",
			expectedValue: map[string]bool{"key1": true, "key2": true, "key3": true},
		},
		{
			name:          "Should panic on missing string map",
			envVarName:    "MISSING_STRING_MAP",
			separator:     ",",
			envVarValue:   "",
			expectedValue: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envVarName, test.envVarValue)

			if test.expectedValue == nil {
				assert.Panics(t, func() {
					MustGetStringMap(test.envVarName, test.separator)
				})
			} else {
				result := MustGetStringMap(test.envVarName, test.separator)
				assert.Equal(t, test.expectedValue, result)
			}
		})
	}
}

func Test_GetStringMap(t *testing.T) {
	tests := []struct {
		name          string
		envVarName    string
		defaultValue  string
		separator     string
		envVarValue   string
		expectedValue map[string]bool
	}{
		{
			name:          "Should return correct string map",
			envVarName:    "TEST_STRING_MAP",
			defaultValue:  "defaultKey",
			separator:     ",",
			envVarValue:   "key1,key2,key3",
			expectedValue: map[string]bool{"key1": true, "key2": true, "key3": true},
		},
		{
			name:          "Should return default string map when missing",
			envVarName:    "MISSING_STRING_MAP",
			defaultValue:  "defaultKey",
			separator:     ",",
			envVarValue:   "",
			expectedValue: map[string]bool{"defaultKey": true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envVarName, test.envVarValue)
			result := GetStringMap(test.envVarName, test.defaultValue, test.separator)
			assert.Equal(t, test.expectedValue, result)
		})
	}
}
