package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"runtime"
	"sort"

	"github.com/jeffvswanson/colorchallenge/pkg/exporttocsv"

	log "github.com/jeffvswanson/colorchallenge/pkg/errorlogging"
)

type colorCode struct {
	Red, Green, Blue uint8
}

type kv struct {
	Key   colorCode
	Value int
}

type logInfo struct {
	Level, Message string
	ErrorMessage   error
}

func init() {
	// Specifically limited to 1 CPU
	runtime.GOMAXPROCS(1)
	log.FormatLog()
}

func main() {

	inputFilename := "input.txt"
	outputFilename := "ColorChallengeOutput"

	// Setup
	status := logInfo{
		Level:   "Info",
		Message: "Beginning setup.",
	}
	log.WriteToLog(status.Level, status.Message, nil)

	// CSV file setup
	status = logInfo{
		Level:   "Info",
		Message: csvSetup(outputFilename),
	}
	log.WriteToLog(status.Level, status.Message, nil)

	// Grab the URLs to parse
	status = logInfo{
		Level:   "Info",
		Message: extractURLs(inputFilename, outputFilename),
	}
	log.WriteToLog(status.Level, status.Message, nil)
}

func csvSetup(filename string) string {

	filename = exporttocsv.CreateCSV(filename)
	headerRecord := []string{"URL", "top_color1", "top_color2", "top_color3"}
	exporttocsv.Export(filename, headerRecord)

	return "CSV setup complete."
}

func extractURLs(inFilename, outFilename string) string {

	f, err := os.Open(inFilename)
	log.ErrorCheck("Fatal", "URL extraction failed during setup.", err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Default behavior is to scan line-by-line
	ch := make(chan []string)
	for scanner.Scan() {
		// We're not interested in keeping the URL and color mapping in
		// memory, just extracting the color mapping.
		go getImageData(scanner.Text(), ch)
		exporttocsv.Export(fmt.Sprintf("%v.csv", outFilename), <-ch)
	}
	return "Process complete."
}

func getImageData(url string, ch chan []string) {
	resp, err := http.Get(url)
	if log.ErrorCheck("Warn", "http.Get failure", err) {
		return
	}
	defer resp.Body.Close()
	// Extract the image information.
	imgData, _, err := image.Decode(resp.Body)
	if log.ErrorCheck("Warn", fmt.Sprintf("%v image decode error", url), err) {
		return
	}

	// Get the output string into url,color,color,color format.
	output := countColors(imgData)
	output[0] = url
	ch <- output
}

func countColors(img image.Image) []string {
	// Find pixel color mapping
	timesAppeared := make(map[colorCode]int)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			// imgData returns in YCbCr format, need to convert to RGB 8-bit
			rgb := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			timesAppeared[colorCode{rgb.R, rgb.G, rgb.B}]++
		}
	}

	// Sort colors in descending order.
	output := sortColors(timesAppeared)

	return output
}

func sortColors(timesAppeared map[colorCode]int) []string {
	// Struct to extract colorCode, key, and times it appeared, value,
	// from the map.
	// Only stable for Go 1.8 and higher

	// Sort the colors from largest value to smallest value.
	var sortAppearances []kv

	for color, appeared := range timesAppeared {
		sortAppearances = append(sortAppearances, kv{color, appeared})
	}
	sort.Slice(sortAppearances, func(i, j int) bool {
		return sortAppearances[i].Value > sortAppearances[j].Value
	})

	output := extractTopColors(sortAppearances)

	return output
}

func extractTopColors(xColors []kv) []string {
	// Extract the top 3 colors
	topColors := make([]string, 4)
	for i := 0; i < 3; i++ {
		// Convert RGB color codes to hexadecimal, #000000 - #FFFFFF
		hexColor := fmt.Sprintf("#%.2X%.2X%.2X", xColors[i].Key.Red,
			xColors[i].Key.Green, xColors[i].Key.Blue)
		topColors[i+1] = hexColor
	}
	return topColors
}

// Start with the end in mind.

// Result written to CSV file in the form of url, top_color1, top_color2, top_color3. O(n) Key = url, value = string of top 3 hexadecimal values

// Convert RGB color scheme (0 - 255, 0 - 255, 0 - 255) (256 bytes or 2^8) to hexadecimal format (#000000 - #FFFFFF) O(1) due to only needing
// to deal with 3 colors.

// Utilize quicksort to sort colors into ascending/descending order and slice off top 3. O(n lg n)

// 1st approach to get colors from image
// Scan image pixel by pixel and increment a counter relating to each color found in the image. Best way in a map. Key = color code,
// value = number of times color found.

// Navigate to the image

// Keep a counter of what image we're on

// Allocate a map to hold the url and it's index of colors. As the algorithm progresses the slice holding the colors will be converted to

// Load in the input.txt file line by line and send off a gofunc for as many lines are as possible while staying within memory and CPU constraints.

// Constant to set max memory used. 512 MB. I'm guessing their using a Docker container or something similar.

// Setup function to initialize log file and csv file to write to.

/*
Ideas:

1. Make sections of the code supporting packages. For example, not all the code needs to be in one main file. The CSV handler could be a package and called into main.

2. Nothing says I'm explicitly limited, just that I may be limited.

3. Given list 1000 urls to an image, to simulate a larger number keep looping around the list. Will this cause a denial of service?

4. Take a wide sample, say 1000 pixels apart. If the pixels are the same value assume all pixels have the same value in between. If not, cut
the sample in half to find where the pixels would be the same.

5. Benchmark how running different numbers of goroutines would affect performance.
	Should the goroutine start after the URLs get extracted or part of the
	extraction process?
	a. Launch a goroutine for each URL
	b. Launch a goroutine for every 10 URLs
	c. Launch a goroutine for every 100 URLs
	d. Launch a goroutine for every 1000 URLs

6. Once program runs dockerize it.

7. Have pointers to the errors passed to errorlogging.

8. Create data structure other than map to support color mapping
*/

/*
Errors encountered:

1. "Get https://i.redd.it/fyqzavufvjwy.jpg: dial tcp: lookup i.redd.it: no such host"
Approach: This is a fatal error if it's more than a few. It means there's no data connection.

*/

/*
Tasks 29 May 2019
- Test benchmark
- Branch using waitgroups instead of channels
- Benchmark waitgroup results
- Branch sampling every other pixel, compare to baseline
*/
