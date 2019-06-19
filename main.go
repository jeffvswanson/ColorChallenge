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
	"sync"

	"github.com/jeffvswanson/colorchallenge/exporttocsv"

	log "github.com/jeffvswanson/colorchallenge/errorlogging"
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

var wg sync.WaitGroup
var logfile, csvfile *os.File

func init() {
	// Specifically limited to 1 CPU
	runtime.GOMAXPROCS(1)

	logfile = log.FormatLog()

	csvfile = exporttocsv.CreateCSV("ColorChallengeOutput")
	headerRecord := []string{"URL", "top_color1", "top_color2", "top_color3"}
	exporttocsv.Export(csvfile, headerRecord)
}

func main() {

	defer logfile.Close()
	defer csvfile.Close()
	inputFilename := "input.txt"

	// Setup
	status := logInfo{
		Level:   "Info",
		Message: "Beginning setup.",
	}
	log.WriteToLog(status.Level, status.Message, nil)

	// Grab the URLs to parse
	status = logInfo{
		Level:   "Info",
		Message: extractURLs(inputFilename, csvfile),
	}
	log.WriteToLog(status.Level, status.Message, nil)
}

// extractURLs pulls the URLs from the given file for image processing.
func extractURLs(inFilename string, csv *os.File) string {

	f, err := os.Open(inFilename)
	log.ErrorCheck("Fatal", "URL extraction failed during setup.", err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Spawn workers to prevent running out of memory.
	urlChan := make(chan string)
	for i := 0; i < 7; i++ {
		wg.Add(1)
		go func() {
			for url := range urlChan {
				imageData(url, csv)
			}
			wg.Done()
		}()
	}

	for scanner.Scan() {
		urlChan <- scanner.Text()
	}
	close(urlChan)

	wg.Wait()

	return "Process complete."
}

// imageData extracts the image from a given URL.
func imageData(url string, csv *os.File) {

	resp, err := http.Get(url)
	if log.ErrorCheck("Warn", "http.Get failure", err) {
		return
	}
	defer resp.Body.Close()
	// Extract the image information.
	img, _, err := image.Decode(resp.Body)
	if log.ErrorCheck("Warn", fmt.Sprintf("%v image decode error", url), err) {
		return
	}
	// Get the output string into url,color,color,color format.
	output := countColors(img)
	output[0] = url
	exporttocsv.Export(csv, output)
}

// countColors finds pixel color mapping of an image in RGB format.
func countColors(img image.Image) []string {

	timesAppeared := make(map[colorCode]int)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			// img returns in YCbCr format, need to convert to RGB 8-bit
			rgb := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			timesAppeared[colorCode{rgb.R, rgb.G, rgb.B}]++
		}
	}
	// Sort colors in descending order.
	output := sortColors(timesAppeared)

	return output
}

// sortColors sorts from the most common color to the least common color.
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

// extractTopColors pulls out the top 3 top colors in the image and
// prints them in hexadecimal format.
func extractTopColors(xColors []kv) []string {

	topColors := make([]string, 4)
	for i := 0; i < 3; i++ {
		// Convert RGB color codes to hexadecimal, #000000 - #FFFFFF
		hexColor := fmt.Sprintf("#%.2X%.2X%.2X", xColors[i].Key.Red,
			xColors[i].Key.Green, xColors[i].Key.Blue)
		topColors[i+1] = hexColor
	}
	return topColors
}
