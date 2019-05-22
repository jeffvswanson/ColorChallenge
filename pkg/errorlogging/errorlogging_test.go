package errorlogging

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestErrorCheckNoError(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	var message string
	var err error
	ErrorCheck(message, err)
	t.Log(buf.String())
}

func TestErrorCheckWithError(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	message := "Test error:"
	err := fmt.Errorf("this is a test error")
	ErrorCheck(message, err)
	t.Log(buf.String())
}
