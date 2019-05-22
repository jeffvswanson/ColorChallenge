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
