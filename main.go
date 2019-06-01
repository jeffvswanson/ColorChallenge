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
	img, _, err := image.Decode(resp.Body)
	if log.ErrorCheck("Warn", fmt.Sprintf("%v image decode error", url), err) {
		return
	}
	// Get the output string into url,color,color,color format.
	output := countColors(img)
	output[0] = url
	ch <- output
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
