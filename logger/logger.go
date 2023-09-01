package logger

import (
	"bufio"
	"log"
	"log/slog"
	"os"
	"regexp"
	"sync"

	"github.com/pokt-foundation/utils-go/environment"
)

const (
	logLevel          = "LOG_LEVEL"
	defaultLogLevel   = "info" // Default log level if no environment variable is set
	logHandler        = "LOG_HANDLER"
	defaultLogHandler = "json" // Default log handler if no environment variable is set

	// Log levels as strings
	logLevelDebug = "debug"
	logLevelInfo  = "info"
	logLevelWarn  = "warn"
	logLevelError = "error"

	// Log handlers as strings
	logHandlerJSON = "json"
	logHandlerText = "text"
)

// logLevelMap maps log levels as strings to their corresponding slog.Level values.
var logLevelMap = map[logLevelStr]slog.Level{
	logLevelDebug: slog.LevelDebug,
	logLevelInfo:  slog.LevelInfo,
	logLevelWarn:  slog.LevelWarn,
	logLevelError: slog.LevelError,
}

// Logger wraps the underlying slog.Logger and keeps track of the current log level.
type (
	Logger struct {
		*slog.Logger
		logLevel logLevelStr
	}

	logLevelStr   string
	logHandlerStr string
)

// isValid checks if a log level string is a valid log level.
func (l logLevelStr) isValid() bool {
	switch l {
	case logLevelDebug, logLevelInfo, logLevelWarn, logLevelError:
		return true
	default:
		return false
	}
}

// isValid checks if a log handler is valid.
func (l logHandlerStr) isValid() bool {
	switch l {
	case logHandlerJSON, logHandlerText:
		return true
	default:
		return false
	}
}

// New creates a new Logger instance.
// It reads the LOG_LEVEL environment variable to set the log level for the new logger.
// Valid log levels are "debug", "info", "warn", and "error".
// If an invalid or missing value is provided, it falls back to the default log level "info".
// The LOG_LEVEL environment variable allows dynamic control over logging verbosity.
// The LOG_HANDLER environment variable allows setting output to JSON or text (default is JSON).
func New() *Logger {
	logLevelVar := logLevelStr(environment.GetString(logLevel, defaultLogLevel))
	if !logLevelVar.isValid() {
		log.Printf("invalid LOG_LEVEL env: %s, using info level default", logLevelVar)
		logLevelVar = defaultLogLevel
	}
	logHandlerVar := logHandlerStr(environment.GetString(logHandler, defaultLogHandler))
	if !logHandlerVar.isValid() {
		log.Printf("invalid LOG_HANDLER env: %s, using json default", logHandlerVar)
		logHandlerVar = defaultLogHandler
	}

	programLevel := new(slog.LevelVar)
	handlerOptions := &slog.HandlerOptions{Level: programLevel}

	// Allow configuration of log handler. default is to use JSON.
	var handler slog.Handler
	switch logHandlerVar {
	case logHandlerText: // If LOG_HANDLER var set to "text", logger will use text output
		handler = slog.NewTextHandler(os.Stderr, handlerOptions)
	default: // If no LOG_HANDLER var set, logger will use JSON output
		handler = slog.NewJSONHandler(os.Stderr, handlerOptions)
	}

	slogger := slog.New(handler)

	// Configure logger - logs levels below the set level will be ignored (default is info)
	logLevel := logLevelMap[logLevelVar]
	programLevel.Set(logLevel)

	return &Logger{Logger: slogger, logLevel: logLevelVar}
}

// LogLevel returns the current log level as a string.
func (l *Logger) LogLevel() string {
	return string(l.logLevel)
}

// Fatal logs an Error level log and exits the program using os.Exit(1).
func (l *Logger) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}

// NewTestLogger creates a new Logger instance and a reader to capture its output.
// It returns a pointer to the logger, a function to read the logged messages'
// `msg=` field as a slice, and a function to clean up resources.
func NewTestLogger() (*Logger, func() []string, func()) {
	// Create a pipe to capture standard error output
	r, w, _ := os.Pipe()
	originalStderr := os.Stderr // Keep track of original stderr
	os.Stderr = w               // Redirect stderr to the write end of the pipe

	var logs []string
	var logsMu sync.Mutex

	// Create the logger using the New function
	logger := New()

	// Compile the regular expression to capture msg value
	re := regexp.MustCompile(`msg="([^"]+)"`)

	// Run a goroutine to capture logs into a slice
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			text := scanner.Text()
			matches := re.FindStringSubmatch(text)
			if len(matches) > 1 {
				logsMu.Lock()
				logs = append(logs, matches[1])
				logsMu.Unlock()
			}
		}
	}()

	// Function to read captured output as a slice
	readOutput := func() []string {
		logsMu.Lock()
		defer logsMu.Unlock()
		clone := make([]string, len(logs))
		copy(clone, logs)
		return clone
	}

	// Function to clean up and restore original stderr
	cleanup := func() {
		os.Stderr = originalStderr
		w.Close()
	}

	return logger, readOutput, cleanup
}
