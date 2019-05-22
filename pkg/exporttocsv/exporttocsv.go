package exporttocsv

import (
	"colorchallenge/pkg/errorlogging"
	"encoding/csv"
	"fmt"
	"os"
)

// Export serves as a wrapper to append a record to the given CSV file.
func Export(filename string, record []string) error {

	if filename[len(filename)-4:] != ".csv" {
		filename = fmt.Sprintf("%v.csv", filename)
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errorlogging.ErrorCheck("Could not open CSV file:", err)
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.Write(record)
	errorlogging.ErrorCheck("Could not write to CSV file:", err)
	w.Flush()

	return err
}

// CreateCSV creates a CSV file in the main project directory.
func CreateCSV(filename string) string {

	filename = fmt.Sprintf("%v.csv", filename)
	f, err := os.Create(filename)
	errorlogging.ErrorCheck("Could not create CSV file:", err)
	defer f.Close()

	return f.Name()
}
