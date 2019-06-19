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
var logfile *os.File

func init() {
	// Specifically limited to 1 CPU
	runtime.GOMAXPROCS(1)
	logfile = log.FormatLog()
}

func main() {

	defer logfile.Close()
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

	outFilename = fmt.Sprintf("%s.csv", outFilename)

	scanner := bufio.NewScanner(f)

	// Spawn workers to prevent running out of memory.
	urlChan := make(chan string)
	for i := 0; i < 7; i++ {
		wg.Add(1)
		go func() {
			for url := range urlChan {
				getImageData(url, outFilename)
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

func getImageData(url, csv string) {

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

func countColors(img image.Image) []string {
	// Find pixel color mapping
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
