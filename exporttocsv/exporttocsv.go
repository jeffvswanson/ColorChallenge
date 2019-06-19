// Package exporttocsv provides functions to create and append to CSV
// files.
package exporttocsv

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jeffvswanson/colorchallenge/errorlogging"
)

// Export serves as a wrapper to append a record to the given CSV file.
func Export(f *os.File, record []string) error {

	w := csv.NewWriter(f)
	err := w.Write(record)
	errorlogging.ErrorCheck("Fatal", "Could not write to CSV file:", err)
	w.Flush()

	return err
}

// CreateCSV creates a CSV file in the main project directory.
func CreateCSV(filename string) (f *os.File) {

	filename = fmt.Sprintf("%s.csv", filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errorlogging.ErrorCheck("Fatal", "Could not create CSV file:", err)

	return f
}
