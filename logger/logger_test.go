package logger

import (
	"bytes"
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

func TestLoggerInfo(t *testing.T) {
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
	}{
		{
			name:      "Info log with single object",
			eventName: "testEvent",
			demoStruct: testStruct{
				StringValue: "lala",
				IntValue:    78,
			},
			objects: []LogObject{
				InfoObject("info error: something went wrong but it's ok tho"),
			},
			wantOutput: []string{
				`testEvent`,
				`"infoMsg": {"msg": "info error: something went wrong but it's ok tho"}`,
				`"demoTest-struct": {"stringProp": "lala", "intProp": 78, "boolProp": false}}`,
				`"service": "test-service"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.objects = append(tt.objects, tt.demoStruct)
			logger.Info(tt.eventName, tt.objects)

			// Read and assert the captured console output
			capturedOutput := mockConsoleWriter.String()
			outputMatched := false
			for _, want := range tt.wantOutput {
				if strings.Contains(capturedOutput, want) {
					outputMatched = true
					break
				}
			}
			if !outputMatched {
				t.Errorf("Expected console output to contain one of the following: %v\nGot: '%s'", tt.wantOutput, capturedOutput)
			}

			// Reset the captured console output for the next test case
			mockConsoleWriter.captured.Reset()
		})
	}
}
