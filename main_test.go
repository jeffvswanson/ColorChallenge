package main

import (
	"image/color"
	"testing"
)

func TestExtractURLs(t *testing.T) {
	expected := "Process complete."
	got := extractURLs("input_test.txt", csvfile)
	if got != expected {
		t.Errorf("URL extraction error. Expected: %v, Got: %v", expected, got)
	}
}

// TestExtractTopColors also tests heapify.
func TestExtractTopColors(t *testing.T) {
	timesAppeared := map[color.Color]int{
		color.YCbCr{0, 128, 128}:   2, // black, #000000
		color.YCbCr{255, 128, 128}: 4, // white, #FFFFFF
		color.YCbCr{75, 84, 255}:   1, // red, #FF0000
		color.YCbCr{149, 43, 21}:   5, // green, #00FF00
		color.YCbCr{29, 255, 107}:  3, // blue, #0000FE, can't quite get close enough to #0000FF due to conversion loss
	}
	expected := []string{"", "#00FF00", "#FFFFFF", "#0000FE"}
	got := heapify(timesAppeared)

	for idx, value := range expected {
		if got[idx] != value {
			t.Errorf("Top color extraction error. Expected: %v at index %d, Got: %v", value, idx, got[idx])
		}
	}
}

// func TestCountColors(t *testing.T) {
// 	// Create an evenly partitioned 3-color test image
// 	width := 90
// 	height := 90

// 	topLeft := image.Point{0, 0}
// 	bottomRight := image.Point{width, height}

// 	img := image.NewNRGBA(image.Rectangle{topLeft, bottomRight})

// 	r := color.NRGBA{255, 0, 0, 0}
// 	g := color.NRGBA{0, 255, 0, 0}
// 	b := color.NRGBA{0, 0, 255, 0}

// 	// Build the even color partitions
// 	for y := 0; y < height; y++ {
// 		for x := 0; x < width; x++ {
// 			switch {
// 			case x < width/3: // left third
// 				img.Set(x, y, r)
// 			case x >= width/3 && x < 2*width/3: // middle third
// 				img.Set(x, y, g)
// 			case x >= 2*width/3: // right third
// 				img.Set(x, y, b)
// 			}
// 		}
// 	}

// 	// Use a map to check if the values are in the returned []string
// 	// since we can't guarantee the order of the slice.
// 	expected := map[string]int{"": 0, "#FF0000": 0, "#00FF00": 0, "#0000FF": 0}
// 	got := countColors(img)
// 	for _, s := range got {
// 		if _, ok := expected[s]; !ok {
// 			t.Errorf("countColors error. Expected: %v, Got: %v", expected, s)
// 		}
// 	}
// }

func BenchmarkExtractURLs(b *testing.B) {
	for n := 0; n < b.N; n++ {
		extractURLs("input_test.txt", csvfile)
	}
}
