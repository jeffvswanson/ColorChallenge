package errorlogging

import (
	"errors"
	"fmt"
	"os"
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

func TestWriteToLog(t *testing.T) {

	var logTests = []struct {
		Level, Message string
		TestErr        error
	}{
		{"Trace", "Trace Notification", errors.New("test trace error")},
		{"Debug", "Debug Notification", errors.New("test debug error")},
		{"Info", "Info Notification", errors.New("test info error")},
		{"Warn", "Warning", errors.New("test warn error")},
		{"Error", "Error! Error!", errors.New("test error")},
		// Have not found a way to test around logrus call to os.Exit on fatal error.
		// {"Fatal", "Fatal Error!", errors.New("test fatal error")},
		{"Panic", "Panic!", errors.New("test panic error")},
	}

	logrus.SetLevel(6)       // Log down to Trace level
	defer logrus.SetLevel(4) // Return logging level to info after test

	for _, tt := range logTests {
		// Need a deferred recover otherwise the check on Panic will cause the test
		// to crash.
		if tt.Level == "Panic" {
			defer func() { recover() }()
		}
		writeToLog(tt.Level, tt.Message, tt.TestErr)
		// Check if appropriately named logfile exists.
		filename := fmt.Sprintf("%v.log", time.Now().Format("20060102"))
		f, err := os.Open(filename)
		if err != nil {
			t.Errorf("writeToLog failure, log file does not exist.")
		}
		defer f.Close()
	}
}
