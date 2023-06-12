package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type (
	// LogLevel for info, warning, error and fatal type of level
	LogLevel string
)

var (
	ErrEmptyService = errors.New("empty service name")

	mu sync.Mutex
)

const (
	Info    LogLevel = "info"
	Warning LogLevel = "warning"
	Error   LogLevel = "error"
	Fatal   LogLevel = "fatal"
)

type Logger struct {
	ServiceName    string
	HasProbability bool
	Probability    float64
	Output         io.Writer
}

func mapObjectsToProperties(objects []LogObject) map[string]interface{} {
	properties := map[string]interface{}{}

	for _, object := range objects {
		properties[object.LogName()] = object.LogProperties()
	}

	return properties
}

// New return a new Logger instance
// add probability if the log shouldn't be logged always
// the probability should be in the form of 1-probability, e.g: probability= 0.90 -> 10% to be logged
func New(service string, hasProbability bool, probability float64) *Logger {
	if service == "" {
		panic(ErrEmptyService)
	}

	output := os.Stderr

	return &Logger{
		ServiceName:    service,
		Probability:    probability,
		HasProbability: hasProbability,
		Output:         output,
	}
}

// Log generate the log with the given parameters y return in stderr a json.
func (logger *Logger) Log(eventName string, level LogLevel, properties map[string]interface{}) {
	floatNumber := rand.Float64()
	if logger.HasProbability && floatNumber >= logger.Probability {
		return
	}

	now := time.Now()

	output := map[string]interface{}{
		"event":      eventName,
		"level":      level,
		"service":    logger.ServiceName,
		"properties": properties,
		"time":       now.Format(time.RFC3339Nano),
	}

	jsonData, err := json.Marshal(output)
	if err != nil {
		log.Println("error marshaling data:", err)
		return
	}

	logger.writeToDev(level, jsonData)
}

// Info logs an info event
// logger.Info("something", []Object{toBeLoggedObject})
func (logger *Logger) Info(eventName string, maxRetentionPeriod time.Duration, objects []LogObject) {
	logger.Log(eventName, Info, mapObjectsToProperties(objects))
}

func (logger *Logger) Warning(eventName string, maxRetentionPeriod time.Duration, objects []LogObject) {
	logger.Log(eventName, Warning, mapObjectsToProperties(objects))
}

// Error logs an error event
func (logger *Logger) Error(eventName string, maxRetentionPeriod time.Duration, objects []LogObject) {
	logger.Log(eventName, Error, mapObjectsToProperties(objects))
}

// Critical logs a fatal event
func (logger *Logger) Critical(eventName string, maxRetentionPeriod time.Duration, objects []LogObject) {
	logger.Log(eventName, Fatal, mapObjectsToProperties(objects))
}

func (logger *Logger) writeToDev(level LogLevel, jsonData []byte) {
	outputDev := logger.Output
	if level != Fatal && outputDev == io.Discard {
		return
	}

	if level == Fatal && outputDev == io.Discard {
		outputDev = os.Stderr
	}

	mu.Lock()
	defer mu.Unlock()

	_, _ = fmt.Fprintln(outputDev, string(jsonData))
}
