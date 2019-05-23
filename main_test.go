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
		t.Logf("csvSetup file name error. Expected: %v, Got: %v", expected, got)
	}
	f, _ := os.Open("TestCSV.csv")
	defer os.Remove("TestCSV.csv")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	expected = "URL,top_color1,top_color2,top_color3"
	got = scanner.Text()
	if got != expected {
		t.Logf("csvSetup header error.\nExpected: %v\n Got:\t%v\n", expected, got)
	}
}

func TestExtractURLs(t *testing.T) {
	// Test first, middle, and last values for expected results
	expected := []string{
		"http://i.imgur.com/FApqk3D.jpg",
		"http://i.imgur.com/TKLs9lo.jpg",
		"https://i.redd.it/lsuw4p2ncyny.jpg",
	}
	got, _ := extractURLs("input.txt")

	for _, url := range expected {
		_, ok := got[url]
		if !ok {
			t.Logf("URL not present in imgColorPrevalence map: %v", url)
		}
	}
	expectedLength := 1000
	gotLength := len(got)
	if gotLength != expectedLength {
		t.Logf("Not all URLs extracted. Expected: %d URLs, Got %d", expectedLength, gotLength)
	}
}
