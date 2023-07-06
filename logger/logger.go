package logger

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ErrEmptyService = errors.New("empty service name error")

	mu sync.Mutex

	defaultCfg = zap.Config{
		Encoding:          "console",
		Level:             zap.NewAtomicLevelAt(zapcore.Level(zapcore.DebugLevel)),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     true,
		DisableStacktrace: true,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			MessageKey:     "event",
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
	}
)

type Logger struct {
	ServiceName    string
	HasProbability bool
	Probability    float64
	Output         *zap.Logger
}

// New return a new Logger instance
// add probability if the log shouldn't be logged always
// the probability should be 1 based. e.g: 0.1 -> 10%
// TODO: Verify that this logger is not blocking any thread
// Note: if you find anything weird after the implementation reach out as soon as you can
func New(service string, hasProbability bool, probability float64, config *zap.Config) (*Logger, error) {
	if service == "" {
		return nil, ErrEmptyService
	}

	var zapLogger *zap.Logger
	var err error

	if config == nil {
		zapLogger, err = defaultCfg.Build()
	} else {
		zapLogger, err = config.Build()
	}

	if err != nil {
		return nil, err
	}

	zap.NewProductionConfig()
	return &Logger{
		ServiceName:    service,
		Probability:    probability,
		HasProbability: hasProbability,
		Output:         zapLogger,
	}, nil
}

// log generate the log with the given parameters
func (logger *Logger) log(eventName string, level zapcore.Level, properties []zapcore.Field) {
	floatNumber := rand.Float64()
	if logger.HasProbability && floatNumber <= logger.Probability {
		return
	}

	now := time.Now()

	LogProperties := []zapcore.Field{
		zap.String("service", logger.ServiceName),
		zap.String("time", now.Format(time.RFC3339Nano)),
	}

	LogProperties = append(LogProperties, properties...)

	logger.writeToConsole(level, eventName, LogProperties)
	// Future idea: maybe send the data directly to our log handler
}

// Info logs an info event
func (logger *Logger) Info(eventName string, objects []LogObject) {
	logger.log(eventName, zapcore.InfoLevel, mapObjectsToZapFields(objects))
}

// Warning logs an warn event
func (logger *Logger) Warning(eventName string, objects []LogObject) {
	logger.log(eventName, zapcore.WarnLevel, mapObjectsToZapFields(objects))
}

// Error logs an error event
func (logger *Logger) Error(eventName string, objects []LogObject) {
	logger.log(eventName, zapcore.ErrorLevel, mapObjectsToZapFields(objects))
}

// Fatal logs a fatal event
func (logger *Logger) Fatal(eventName string, objects []LogObject) {
	logger.log(eventName, zapcore.FatalLevel, mapObjectsToZapFields(objects))
}

func (logger *Logger) writeToConsole(level zapcore.Level, eventName string, properties []zapcore.Field) {
	mu.Lock()
	defer mu.Unlock()

	logger.Output.Log(level, eventName, properties...)
}
