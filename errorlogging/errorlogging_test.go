package errorlogging

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestCreateLogfileName(t *testing.T) {
	expected := fmt.Sprintf("%v.log", time.Now().Format("20060102"))
	got := createLogfileName()
	if got != expected {
		t.Errorf("Created logfile name incorrect.\nExpected: %v, Got: %v", expected, got)
	}
}

func TestLogrus(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.Error("Test error")

	expectedLength := 1
	gotLength := len(hook.Entries)
	if gotLength != expectedLength {
		t.Errorf("logger write fail. Expected: %d, Got: %d", expectedLength, gotLength)
	}

	expectedLevel := logrus.ErrorLevel
	gotLevel := hook.LastEntry().Level
	if gotLevel != expectedLevel {
		t.Errorf("logger error level fail. Expected: %v, Got: %v", expectedLevel, gotLevel)
	}

	gotMessage := hook.LastEntry().Message
	if gotMessage != "Test error" {
		t.Errorf("logger message fail. Expected: %v, Got: %v", "Test error", gotMessage)
	}

	hook.Reset()
}

func TestWrite(t *testing.T) {

	type logInput struct {
		Level   string
		TestErr error
	}

	var inputs = []logInput{
		{"Trace Notification", errors.New("test trace error")},
		{"Debug Notification", errors.New("test debug error")},
		{"Info Notification", errors.New("test info error")},
		{"Warning", errors.New("test warn error")},
		{"Error! Error!", errors.New("test error")},
		{"Fatal Error!", errors.New("test fatal error")},
		{"Panic!", errors.New("test panic error")},
	}

	tests := map[string]struct {
		input logInput
	}{
		"Trace":   {input: inputs[0]},
		"Debug":   {input: inputs[1]},
		"Info":    {input: inputs[2]},
		"Warning": {input: inputs[3]},
		"Error":   {input: inputs[4]},
		"Panic":   {input: inputs[5]},
		// Fatal works, but causes a crash included if a workaround becomes available.
		// "Fatal":   {input: inputs[6]},
	}

	logrus.SetLevel(6)       // Log down to Trace level
	defer logrus.SetLevel(4) // Return logging level to info after test

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Need a deferred recover otherwise the check on Panic will cause the test
			// to crash.
			if name == "Panic" {
				defer func() { recover() }()
			}
			got := captureOutput(func() {
				Write(name, tc.input.Level, tc.input.TestErr)
			})
			expected := fmt.Sprintf("time=\"%v\" level=%v msg=\"%v - %v\"\n", time.Now().Format(time.RFC3339), strings.ToLower(name), tc.input.Level, tc.input.TestErr)
			if got != expected {
				t.Errorf("\nExpected: %v != \nGot: %v", expected, got)
			}
		})
	}
}

// captureOutput is a helper function to capture the output stream for testing.
func captureOutput(f func()) string {
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	f()
	logrus.SetOutput(os.Stderr)
	return buf.String()
}
