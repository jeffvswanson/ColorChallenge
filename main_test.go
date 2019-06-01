package main

import (
	"bufio"
	"bytes"
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
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
	got := extractURLs("input_test.txt", "ColorChallengeOutput")
	if got != expected {
		t.Errorf("URL extraction error. Expected: %v, Got: %v", expected, got)
	}
}

func TestExtractURLsBigInput(t *testing.T) {
	// Test just checks if it can access an input file with 1 billion URLs
	go extractURLs("big_input_test.txt", "ColorChallengeOutput")
}

func TestExtractTopColors(t *testing.T) {
	xColorAppearance := []kv{
		{colorCode{0, 0, 0}, 5},       // black
		{colorCode{255, 255, 255}, 4}, // white
		{colorCode{255, 0, 0}, 3},     // red
		{colorCode{0, 255, 0}, 2},     // green
		{colorCode{0, 0, 255}, 1},     // blue
	}
	expected := []string{"", "#000000", "#FFFFFF", "#FF0000"}
	got := extractTopColors(xColorAppearance)

	for idx, value := range expected {
		if got[idx] != value {
			t.Errorf("Top color extraction error. Expected: %v at index %d, Got: %v", value, idx, got[idx])
		}
	}
}

func TestSortColors(t *testing.T) {
	timesAppeared := map[colorCode]int{
		colorCode{0, 0, 0}:       2, // black
		colorCode{255, 255, 255}: 4, // white
		colorCode{255, 0, 0}:     1, // red
		colorCode{0, 255, 0}:     5, // green
		colorCode{0, 0, 255}:     3, // blue
	}
	expected := []string{"", "#00FF00", "#FFFFFF", "#0000FF"}
	got := sortColors(timesAppeared)

	for idx, value := range expected {
		if got[idx] != value {
			t.Errorf("sortcolors error. Expected: %v at index %d, Got: %v", value, idx, got[idx])
		}
	}
}

func TestCountColors(t *testing.T) {
	// Create an evenly partitioned 3-color test image
	width := 90
	height := 90

	topLeft := image.Point{0, 0}
	bottomRight := image.Point{width, height}

	img := image.NewNRGBA(image.Rectangle{topLeft, bottomRight})

	r := color.NRGBA{255, 0, 0, 0}
	g := color.NRGBA{0, 255, 0, 0}
	b := color.NRGBA{0, 0, 255, 0}

	// Build the even color partitions
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			switch {
			case x < width/3: // left third
				img.Set(x, y, r)
			case x >= width/3 && x < 2*width/3: // middle third
				img.Set(x, y, g)
			case x >= 2*width/3: // right third
				img.Set(x, y, b)
			}
		}
	}

	// Use a map to check if the values are in the returned []string
	// since we can't guarantee the order of the slice.
	expected := map[string]int{"": 0, "#FF0000": 0, "#00FF00": 0, "#0000FF": 0}
	got := countColors(img)
	for _, s := range got {
		if _, ok := expected[s]; !ok {
			t.Errorf("countColors error. Expected: %v, Got: %v", expected, s)
		}
	}
}

func TestGetImageData(t *testing.T) {
	// Testing errorChecks only as happy path is covered by TestExtractURLs.
	ch := make(chan []string)
	var getImageTests = []struct {
		URL      string
		testChan chan []string
	}{
		// Not found error
		{"https://malformedURL", ch},
		// Incorrect file format
		{"https://github.com/edent/SuperTinyIcons/blob/master/images/svg/stackoverflow.svg", ch},
	}

	// Capture the log stream for testing
	var buf bytes.Buffer
	logrus.SetOutput(&buf)

	for _, tt := range getImageTests {
		getImageData(tt.URL, tt.testChan)
		if buf.String() == "" {
			t.Log("Expected an error string. Got an empty string.")
		}
	}
}

func BenchmarkExtractURLs(b *testing.B) {
	for n := 0; n < b.N; n++ {
		extractURLs("input_test.txt", "ColorChallengeOutput")
	}
	// Channel benchmark:	14133294000 ns/op	1590616752 B/op	198078710 allocs/op
	// Mutex on map: 		16588424800 ns/op	1592626208 B/op	198085264 allocs/op
	// Channel on map:    	After waiting five minutes, killed process
	// Image conversion:  	68627798300 ns/op	6488996328 B/op	222952381 allocs/op
	// Image conversion 2:  15187969200 ns/op	2037514560 B/op	198079323 allocs/op
}
