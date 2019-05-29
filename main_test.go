package main

import (
	"bufio"
	"os"
	"testing"
)

func TestCsvSetup(t *testing.T) {
	expected := "CSV setup complete."
	got := csvSetup("TestCSV")
	if got != expected {
		t.Errorf("csvSetup file name error. Expected: %v, Got: %v", expected, got)
	}

	f, _ := os.Open("TestCSV.csv")
	defer os.Remove("TestCSV.csv")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	expected = "URL,top_color1,top_color2,top_color3"
	got = scanner.Text()
	if got != expected {
		t.Errorf("csvSetup header error.\nExpected: %v\n Got:\t%v\n", expected, got)
	}
}

func TestExtractURLs(t *testing.T) {
	expected := "Process complete."
	got := extractURLs("input_test.txt")
	if got != expected {
		t.Errorf("URL extraction error. Expected: %v, Got: %v", expected, got)
	}
}

// func TestExtractTopColors(t *testing.T) {

// }
