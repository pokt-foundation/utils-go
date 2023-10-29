package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Logger_JSON(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name          string
		envLogLevel   logLevelStr
		expectedOut   []map[string]interface{}
		defaultToJSON bool
	}{
		{
			name:          "Should log at debug level AS json if env var is not set (use default)",
			defaultToJSON: true,
			envLogLevel:   logLevelDebug,
			expectedOut: []map[string]interface{}{
				{"level": "DEBUG", "msg": "Debug message"},
				{"level": "INFO", "msg": "Info message"},
				{"level": "WARN", "msg": "Warn message"},
				{"level": "ERROR", "msg": "Error message"},
			},
		},
		{
			name:        "Should log at debug level",
			envLogLevel: logLevelDebug,
			expectedOut: []map[string]interface{}{
				{"level": "DEBUG", "msg": "Debug message"},
				{"level": "INFO", "msg": "Info message"},
				{"level": "WARN", "msg": "Warn message"},
				{"level": "ERROR", "msg": "Error message"},
			},
		},
		{
			name:        "Should log at info level",
			envLogLevel: logLevelInfo,
			expectedOut: []map[string]interface{}{
				{"level": "INFO", "msg": "Info message"},
				{"level": "WARN", "msg": "Warn message"},
				{"level": "ERROR", "msg": "Error message"},
			},
		},
		{
			name:        "Should log at warn level",
			envLogLevel: logLevelWarn,
			expectedOut: []map[string]interface{}{
				{"level": "WARN", "msg": "Warn message"},
				{"level": "ERROR", "msg": "Error message"},
			},
		},
		{
			name:        "Should log at error level",
			envLogLevel: logLevelError,
			expectedOut: []map[string]interface{}{
				{"level": "ERROR", "msg": "Error message"},
			},
		},
		{
			name:        "Should default to logging at info level if invalid log level is set",
			envLogLevel: "what_is_this",
			expectedOut: []map[string]interface{}{
				{"level": "INFO", "msg": "Info message"},
				{"level": "WARN", "msg": "Warn message"},
				{"level": "ERROR", "msg": "Error message"},
			},
		},
		{
			name:        "Should default to logging at info level if no log level is set",
			envLogLevel: "",
			expectedOut: []map[string]interface{}{
				{"level": "INFO", "msg": "Info message"},
				{"level": "WARN", "msg": "Warn message"},
				{"level": "ERROR", "msg": "Error message"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.defaultToJSON {
				// Set log handler to json for the test
				err := os.Setenv(logHandler, string(logHandlerJSON))
				c.NoError(err)
			}

			// Set environment variable for the test
			err := os.Setenv(logLevel, string(test.envLogLevel))
			c.NoError(err)

			// Create a pipe to capture standard error output
			r, w, err := os.Pipe()
			c.NoError(err)
			originalStderr := os.Stderr
			os.Stderr = w

			// Create the logger using the New function
			logger := New()

			// Log the messages at different levels
			logger.Debug("Debug message")
			logger.Info("Info message")
			logger.Warn("Warn message")
			logger.Error("Error message")

			// Restore original stderr
			os.Stderr = originalStderr
			w.Close()

			now := time.Now()

			// Read the captured output
			var buffer bytes.Buffer
			_, err = io.Copy(&buffer, r)
			c.NoError(err)

			actualOutput := buffer.String()
			actualLines := strings.Split(strings.TrimSpace(actualOutput), "\n")

			c.Equal(len(test.expectedOut), len(actualLines))

			for i, line := range actualLines {
				var actualMap map[string]interface{}
				err := json.Unmarshal([]byte(line), &actualMap)
				c.NoError(err)

				c.Equal(test.expectedOut[i]["level"], actualMap["level"])
				c.Equal(test.expectedOut[i]["msg"], actualMap["msg"])

				actualTimestamp := actualMap["time"].(string)

				// Try parsing the timestamp in the first format
				parsedTime, err := time.Parse("2006-01-02T15:04:05.000-07:00", actualTimestamp)
				if err != nil {
					// If the first format fails, try the second format
					parsedTime, err = time.Parse(time.RFC3339, actualTimestamp)
				}
				c.NoError(err, "Timestamp is not in the expected format")

				c.True(now.Sub(parsedTime) < 100*time.Millisecond, "Timestamp is not within 100ms of current time")
			}
		})
	}
}

func Test_Logger_Text(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name        string
		envLogLevel logLevelStr
		expectedOut string
	}{
		{
			name:        "Should log at debug level",
			envLogLevel: logLevelDebug,
			expectedOut: "time=<CURRENT TIMESTAMP> level=DEBUG msg=\"Debug message\"\ntime=<CURRENT TIMESTAMP> level=INFO msg=\"Info message\"\ntime=<CURRENT TIMESTAMP> level=WARN msg=\"Warn message\"\ntime=<CURRENT TIMESTAMP> level=ERROR msg=\"Error message\"",
		},
		{
			name:        "Should log at info level",
			envLogLevel: logLevelInfo,
			expectedOut: "time=<CURRENT TIMESTAMP> level=INFO msg=\"Info message\"\ntime=<CURRENT TIMESTAMP> level=WARN msg=\"Warn message\"\ntime=<CURRENT TIMESTAMP> level=ERROR msg=\"Error message\"",
		},
		{
			name:        "Should log at warn level",
			envLogLevel: logLevelWarn,
			expectedOut: "time=<CURRENT TIMESTAMP> level=WARN msg=\"Warn message\"\ntime=<CURRENT TIMESTAMP> level=ERROR msg=\"Error message\"",
		},
		{
			name:        "Should log at error level",
			envLogLevel: logLevelError,
			expectedOut: "time=<CURRENT TIMESTAMP> level=ERROR msg=\"Error message\"",
		},
		{
			name:        "Should default to logging at info level if invalid log level is set",
			envLogLevel: "what_is_this",
			expectedOut: "time=<CURRENT TIMESTAMP> level=INFO msg=\"Info message\"\ntime=<CURRENT TIMESTAMP> level=WARN msg=\"Warn message\"\ntime=<CURRENT TIMESTAMP> level=ERROR msg=\"Error message\"",
		},
		{
			name:        "Should default to logging at info level if no log level is set",
			envLogLevel: "",
			expectedOut: "time=<CURRENT TIMESTAMP> level=INFO msg=\"Info message\"\ntime=<CURRENT TIMESTAMP> level=WARN msg=\"Warn message\"\ntime=<CURRENT TIMESTAMP> level=ERROR msg=\"Error message\"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set log handler to text for the test
			err := os.Setenv(logHandler, string(logHandlerText))
			c.NoError(err)
			// Set environment variable for the test
			err = os.Setenv(logLevel, string(test.envLogLevel))
			c.NoError(err)

			// Create a pipe to capture standard error output
			r, w, err := os.Pipe()
			c.NoError(err)
			originalStderr := os.Stderr
			os.Stderr = w

			// Create the logger using the New function
			logger := New()

			// Log the messages at different levels
			logger.Debug("Debug message")
			logger.Info("Info message")
			logger.Warn("Warn message")
			logger.Error("Error message")

			// Restore original stderr
			os.Stderr = originalStderr
			w.Close()

			now := time.Now()

			// Replace the placeholder for the current timestamp
			expectedOut := strings.ReplaceAll(test.expectedOut, "<CURRENT TIMESTAMP>", now.Format("2006-01-02T15:04:05.000-07:00"))

			// Read the captured output
			var buffer bytes.Buffer
			_, err = io.Copy(&buffer, r)
			c.NoError(err)

			actualOutput := buffer.String()
			actualLines := strings.Split(strings.TrimSpace(actualOutput), "\n")
			expectedLines := strings.Split(strings.TrimSpace(expectedOut), "\n")

			// Compare the number of lines
			c.Equal(len(expectedLines), len(actualLines))

			// Compare the timestamps and the rest of each line
			for i := 0; i < len(expectedLines); i++ {
				actualTimestamp := strings.SplitN(actualLines[i], " ", 2)[0][5:]

				// Try parsing the timestamp in the first format
				parsedTime, err := time.Parse("2006-01-02T15:04:05.000-07:00", actualTimestamp)
				if err != nil {
					// If the first format fails, try the second format
					parsedTime, err = time.Parse(time.RFC3339, actualTimestamp)
				}
				c.NoError(err, "Timestamp is not in the expected format")

				c.True(now.Sub(parsedTime) < 1000*time.Millisecond, "Timestamp is not within 1000ms of current time")

				// Convert the expected line to use "Z" if the timezone is UTC
				expectedLine := expectedLines[i]
				if parsedTime.UTC().Equal(parsedTime) {
					expectedLine = strings.ReplaceAll(expectedLine, "+00:00", "Z")
				}

				// Compare the rest of the line
				c.Equal(expectedLine, actualLines[i])
			}
		})
	}
}

func Test_NewTestLogger(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name         string
		logHandler   logHandlerStr
		logMessages  []string
		expectedLogs []string
	}{
		{
			name:       "Should log multiple levels in JSON",
			logHandler: logHandlerJSON,
			logMessages: []string{
				"Debug message",
				"Info message",
				"Warn message",
				"Error message",
			},
			expectedLogs: []string{
				"Debug message",
				"Info message",
				"Warn message",
				"Error message",
			},
		},
		{
			name:       "Should log single level in JSON",
			logHandler: logHandlerJSON,
			logMessages: []string{
				"Debug message",
			},
			expectedLogs: []string{
				"Debug message",
			},
		},
		{
			name:       "Should log multiple levels in Text",
			logHandler: logHandlerText,
			logMessages: []string{
				"Debug message",
				"Info message",
				"Warn message",
				"Error message",
			},
			expectedLogs: []string{
				"Debug message",
				"Info message",
				"Warn message",
				"Error message",
			},
		},
		{
			name:       "Should log single level in Text",
			logHandler: logHandlerText,
			logMessages: []string{
				"Debug message",
			},
			expectedLogs: []string{
				"Debug message",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set the environment variable for log handler
			os.Setenv("LOG_HANDLER", string(test.logHandler))

			logger, readOutput, cleanup := NewTestLogger()
			defer cleanup()

			for _, msg := range test.logMessages {
				logger.Info(msg)
			}

			// Small delay to give the goroutine time to capture the logs
			time.Sleep(100 * time.Millisecond)

			logLines := readOutput()

			c.Equal(len(test.expectedLogs), len(logLines), "Number of logged lines should match expected")

			for i, logLine := range logLines {
				c.Equal(test.expectedLogs[i], logLine, "Logged message should match expected")
			}

			// Clean up the environment variable
			os.Unsetenv("LOG_HANDLER")
		})
	}
}
