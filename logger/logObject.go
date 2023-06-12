package logger

type LogObject interface {
	LogName() string
	LogProperties() map[string]interface{}
}

// MapObject create LogObject with data
func MapObject(logName string, data map[string]interface{}) LogObject {
	return &Wrapper{
		Name:   logName,
		Fields: data,
	}
}

// Wrapper wraps any object as a loggable object
type Wrapper struct {
	Name   string
	Fields map[string]interface{}
}

// LogName returns the name of the object to log
func (w *Wrapper) LogName() string {
	return w.Name
}

// LogProperties returns the fields to log
func (w *Wrapper) LogProperties() map[string]interface{} {
	return w.Fields
}

// ErrObject returns an LogObject containing error information
func ErrObject(err error) LogObject {
	return &Wrapper{
		Name: "error",
		Fields: map[string]interface{}{
			"error": err.Error(),
		},
	}
}

// InfoObject returns an LogObject containing information
func InfoObject(err error) LogObject {
	return &Wrapper{
		Name: "info",
		Fields: map[string]interface{}{
			"error": err.Error(),
		},
	}
}
