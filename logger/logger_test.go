package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Logger(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name        string
		envLogLevel string
		expectedOut string
	}{
		{
			name:        "Should log at debug level",
			envLogLevel: "debug",
			expectedOut: "time=<CURRENT TIMESTAMP> level=DEBUG msg=\"Debug message\"\ntime=<CURRENT TIMESTAMP> level=INFO msg=\"Info message\"\ntime=<CURRENT TIMESTAMP> level=WARN msg=\"Warn message\"\ntime=<CURRENT TIMESTAMP> level=ERROR msg=\"Error message\"",
		},
		{
			name:        "Should log at info level",
			envLogLevel: "info",
			expectedOut: "time=<CURRENT TIMESTAMP> level=INFO msg=\"Info message\"\ntime=<CURRENT TIMESTAMP> level=WARN msg=\"Warn message\"\ntime=<CURRENT TIMESTAMP> level=ERROR msg=\"Error message\"",
		},
		{
			name:        "Should log at warn level",
			envLogLevel: "warn",
			expectedOut: "time=<CURRENT TIMESTAMP> level=WARN msg=\"Warn message\"\ntime=<CURRENT TIMESTAMP> level=ERROR msg=\"Error message\"",
		},
		{
			name:        "Should log at error level",
			envLogLevel: "error",
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
			// Set environment variable for the test
			err := os.Setenv(logLevel, test.envLogLevel)
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

				c.True(now.Sub(parsedTime) < 100*time.Millisecond, "Timestamp is not within 100ms of current time")

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
