package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/pokt-foundation/utils-go/environment"
)

const (
	logLevel        = "LOG_LEVEL"
	defaultLogLevel = "info" // Default log level if no environment variable is set
)

// Logger wraps the underlying slog.Logger and keeps track of the current log level.
type Logger struct {
	*slog.Logger
	logLevel logLevelStr
}

type logLevelStr string

// logLevelMap maps log levels as strings to their corresponding slog.Level values.
var logLevelMap = map[logLevelStr]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// isValid checks if a log level string is a valid log level.
func (l logLevelStr) isValid() bool {
	switch l {
	case "debug", "info", "warn", "error":
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
func New() *Logger {
	logLevelVar := logLevelStr(environment.GetString(logLevel, defaultLogLevel))
	if !logLevelVar.isValid() {
		log.Printf("invalid LOG_LEVEL env: %s, using info level default", logLevelVar)
		logLevelVar = defaultLogLevel
	}

	programLevel := new(slog.LevelVar)
	handlerOptions := &slog.HandlerOptions{Level: programLevel}
	textHandler := slog.NewTextHandler(os.Stderr, handlerOptions)

	slogger := slog.New(textHandler)

	// Configure logger - logs levels below the set level will be ignored (default is info)
	logLevel := logLevelMap[logLevelVar]
	programLevel.Set(logLevel)

	return &Logger{Logger: slogger, logLevel: logLevelVar}
}

// LogLevel returns the current log level as a string.
func (l *Logger) LogLevel() string {
	return string(l.logLevel)
}
