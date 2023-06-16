package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogObject interface {
	// LogName returns the name of the object to log
	LogName() string
	// LogProperties returns the fields to log
	// NOTE: Avoid using struct as a value as it will trigger a reflection field
	// only use if it's mandatory
	LogProperties() map[string]any
}

// MapObject create LogObject with data
func MapObject(logName string, data map[string]any) LogObject {
	return &wrapper{
		Name:   logName,
		Fields: data,
	}
}

// wrapper wraps any object as a loggable object
type wrapper struct {
	Name   string
	Fields map[string]any
}

// LogName returns the name of the object to log
func (w *wrapper) LogName() string {
	return w.Name
}

// LogProperties returns the fields to log
func (w *wrapper) LogProperties() map[string]any {
	return w.Fields
}

// ErrObject returns an LogObject containing error information
func ErrObject(err error) LogObject {
	fields := map[string]any{
		"error": "empty error",
	}
	if err != nil {
		fields["error"] = err.Error()
	}

	return &wrapper{
		Name:   "errorMsg",
		Fields: fields}
}

// InfoObject returns an LogObject containing information
func InfoObject(msg string) LogObject {
	return &wrapper{
		Name: "infoMsg",
		Fields: map[string]any{
			"msg": msg,
		},
	}
}

// InfoObject returns an LogObject containing information
func WarnObject(err error) LogObject {
	fields := map[string]any{
		"error": "empty error",
	}
	if err != nil {
		fields["error"] = err.Error()
	}

	return &wrapper{
		Name:   "warnMsg",
		Fields: fields,
	}
}

func mapObjectsToZapFields(objects []LogObject) []zapcore.Field {
	fields := []zapcore.Field{}
	for _, object := range objects {
		fields = append(fields, zap.Object(object.LogName(), logPropertiesMarshaler(object.LogProperties())))
	}

	return fields
}

type logPropertiesMarshaler map[string]any

// MarshalLogObject handles every field to be added and ensures it's of the correct type, using reflect as the last resourse
func (e logPropertiesMarshaler) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for key, value := range e {
		switch v := value.(type) {
		case string:
			if v != "" {
				enc.AddString(key, v)
			}
		case int:
			if v != 0 {
				enc.AddInt(key, v)
			}
		case bool:
			enc.AddBool(key, v)
		case time.Time:
			enc.AddTime(key, v)
		case float64:
			if v != 0 {
				enc.AddFloat64(key, v)
			}
		case float32:
			if v != 0 {
				enc.AddFloat32(key, v)
			}
		case int64:
			if v != 0 {
				enc.AddInt64(key, v)
			}
		case int32:
			if v != 0 {
				enc.AddInt32(key, v)
			}
		case int16:
			if v != 0 {
				enc.AddInt16(key, v)
			}
		case int8:
			if v != 0 {
				enc.AddInt8(key, v)
			}
		case []byte:
			if len(v) != 0 {
				enc.AddBinary(key, v)
			}
		case []rune:
			if len(v) != 0 {
				enc.AddByteString(key, []byte(string(v)))
			}
		case complex128:
			if v != 0 {
				enc.AddComplex128(key, v)
			}
		case complex64:
			if v != 0 {
				enc.AddComplex64(key, v)
			}
		case time.Duration:
			if v != 0 {
				enc.AddDuration(key, v)
			}
		case uint:
			if v != 0 {
				enc.AddUint(key, v)
			}
		case uint64:
			if v != 0 {
				enc.AddUint64(key, v)
			}
		case uint32:
			if v != 0 {
				enc.AddUint32(key, v)
			}
		case uint16:
			if v != 0 {
				enc.AddUint16(key, v)
			}
		case uint8:
			if v != 0 {
				enc.AddUint8(key, v)
			}
		case uintptr:
			if v != 0 {
				enc.AddUintptr(key, v)
			}
		// Avoid this case, it's used as the last resort
		default:
			err := enc.AddReflected(key, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
