package logger

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MockConsoleWriter struct {
	captured bytes.Buffer
}

func (m *MockConsoleWriter) Write(p []byte) (n int, err error) {
	m.captured.Write(p)
	return len(p), nil
}

func (m *MockConsoleWriter) String() string {
	return m.captured.String()
}

type testStruct struct {
	StringValue   string
	IntValue      int
	BoolValue     bool
	IgnoredField1 string
	IgnoredField2 int
}

func (t testStruct) LogProperties() map[string]any {
	return map[string]any{
		"stringProp": t.StringValue,
		"intProp":    t.IntValue,
		"boolProp":   t.BoolValue,
	}
}

func (t testStruct) LogName() string {
	return "demoTest-struct"
}

func TestLogger(t *testing.T) {
	mockConsoleWriter := &MockConsoleWriter{}
	logger := &Logger{
		ServiceName:    "test-service",
		HasProbability: false,
		Probability:    0.0,
		Output: zap.New(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
				zapcore.AddSync(mockConsoleWriter),
				zapcore.DebugLevel,
			)),
	}

	tests := []struct {
		name       string
		eventName  string
		objects    []LogObject
		demoStruct testStruct
		wantOutput []string
		logMethod  func(logger *Logger, eventName string, objects []LogObject)
	}{
		{
			name:      "Info log with single object",
			eventName: "testEvent",
			demoStruct: testStruct{
				StringValue: "infoTest",
				IntValue:    78,
			},
			objects: []LogObject{
				InfoObject("info error: something went wrong but it's ok tho"),
			},
			wantOutput: []string{
				`testEvent`,
				`"infoMsg": {"msg": "info error: something went wrong but it's ok tho"}`,
				`"demoTest-struct"`,
				`"stringProp": "infoTest"`,
				`"intProp": 78`,
				`"boolProp": false`,
			},
			logMethod: (*Logger).Info,
		},
		{
			name:      "Warning log with multiple objects",
			eventName: "warningEvent",
			demoStruct: testStruct{
				StringValue: "warnTest",
				BoolValue:   true,
			},
			objects: []LogObject{
				WarnObject(errors.New("Warning: app not found in this region")),
			},
			wantOutput: []string{
				`warningEvent`,
				`"warnMsg": {"error": "Warning: app not found in this region"}`,
				`"demoTest-struct"`,
				`"stringProp": "warnTest"`,
				`"boolProp": true`,
			},
			logMethod: (*Logger).Warning,
		},
		{
			name:      "Error log with multiple objects",
			eventName: "errorEvent",
			demoStruct: testStruct{
				StringValue: "errTest",
			},
			objects: []LogObject{
				ErrObject(errors.New("error occurred")),
				MapObject("extraData", map[string]interface{}{
					"key1": "value1",
					"key2": 123,
				}),
				ErrObject(nil),
			},
			wantOutput: []string{
				`errorEvent`,
				`"errorMsg": {"error": "error occurred"}`,
				`"extraData"`,
				`"key1": "value1"`,
				`"key2": 123`,
				`"demoTest-struct"`,
				`"stringProp": "errTest"`,
				`"boolProp": false`,
			},
			logMethod: (*Logger).Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the captured console output for the next test case
			mockConsoleWriter.captured.Reset()
			tt.objects = append(tt.objects, tt.demoStruct)
			tt.logMethod(logger, tt.eventName, tt.objects)

			// Read and assert the captured console output
			capturedOutput := mockConsoleWriter.String()

			for _, want := range tt.wantOutput {
				if !strings.Contains(capturedOutput, want) {
					t.Errorf("Expected console output to contain '%s', got: '%s'", want, capturedOutput)
				}
			}
		})
	}
}
