// Package errorlogging provides a clean function to allow for error
// checking throughout the module instead of the if err != nil
// pattern.
package errorlogging

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// FormatLog sets up the logging format.
func FormatLog() {
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05.0000"
	Formatter.FullTimestamp = true
	logrus.SetFormatter(Formatter)
}

// ErrorCheck writes an error message and the error to a log.
func ErrorCheck(level, message string, err error) bool {

	var isError bool
	if err != nil {
		WriteToLog(level, message, err)
		isError = true
	}
	return isError
}

// WriteToLog writes messages to a log.
func WriteToLog(level, message string, reportedErr error) {
	filename := createLogfileName()
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	FormatLog()
	if err != nil {
		// Cannot open log file, defaulting to stderr.
		fmt.Println(err)
	} else {
		logrus.SetOutput(f)
	}

	logMessage := fmt.Sprintf("%v - %v", message, reportedErr)
	switch level {
	case "Trace":
		logrus.Trace(logMessage)
	case "Debug":
		logrus.Debug(logMessage)
	case "Info":
		logrus.Info(logMessage)
	case "Warn":
		logrus.Warn(logMessage)
	case "Error":
		logrus.Error(logMessage)
	case "Fatal":
		logrus.Fatal(logMessage)
	case "Panic":
		logrus.Panic(logMessage)
	}
}

func createLogfileName() string {
	return fmt.Sprintf("%v.log", time.Now().Format("20060102"))
}
