package exporttocsv

import (
	"bufio"
	"os"
	"testing"
)

// TestCreateCSV creates a CSV file, checks if it exists, and cleans up.
func TestCreateCSV(t *testing.T) {

	filename := "TestCSV"

	f := CreateCSV(filename)

	err := f.Close()
	if err != nil {
		t.Errorf("%v\n", err)
	}
	err = os.Remove(f.Name())
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestExport(t *testing.T) {

	var err error
	filename := "TestCSV"
	f := CreateCSV(filename)
	filename = f.Name()

	record := []string{"item 1", "item2", "item3"}

	expected := 10
	for i := 0; i < expected; i++ {
		Export(f, record)
	}
	f.Close()

	csv, err := os.Open(filename)
	if err != nil {
		t.Errorf("%v\n", err)
	}
	defer func() {
		err = csv.Close()
		if err != nil {
			t.Errorf("%v\n", err)
		}
		err = os.Remove(filename)
		if err != nil {
			t.Errorf("%v\n", err)
		}
	}()

	scanner := bufio.NewScanner(csv)
	got := 0
	for scanner.Scan() {
		got++
	}
	if got != expected {
		t.Errorf("Expected: %d lines in the CSV file.\n", expected)
		t.Errorf("Got: %d lines.\n", got)
	}

}
