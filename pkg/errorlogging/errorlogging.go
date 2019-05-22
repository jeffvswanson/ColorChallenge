// Package errorlogging provides a clean function to allow for error
// checking throughout the module instead of the if err != nil
// pattern.
package errorlogging

import (
	"log"
)

// ErrorCheck writes an error message and the error to a log.
func ErrorCheck(message string, err error) {
	if err != nil {
		log.Println(message, err)
	}
}

/*
TODO: Implement logging such that the errors are written to an output file.

The output file goes into a directory:
	logs/
		YYYY/
			MM-Month/
				YYYYMMMDD.log

Format of log entry:
	Prefix = YYYY-MM-DD HH:MM:SS,ms UTC - package name - ErrorLevel - message
*/
