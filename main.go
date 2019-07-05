package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/jeffvswanson/colorchallenge/exporttocsv"

	log "github.com/jeffvswanson/colorchallenge/errorlogging"
)

type colorCode struct {
	Red, Green, Blue uint8
}

type colorNode struct {
	Color       colorCode
	Occurrences int
}

// A colorHeap is a max-heap of the colors found from an image.
type colorHeap []colorNode

func (c colorHeap) Len() int           { return len(c) }
func (c colorHeap) Less(i, j int) bool { return c[i].Occurrences < c[j].Occurrences }
func (c colorHeap) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c *colorHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the
	// slices's length not just its contents.
	*c = append(*c, x.(colorNode))
}

func (c *colorHeap) Pop() interface{} {
	old := *c
	n := len(old)
	x := old[n-1]
	*c = old[0 : n-1]
	return x
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
	log.WriteToLog("Info", "Beginning setup", nil)

	// Grab the URLs to parse
	status := extractURLs(inputFilename, csvfile)
	log.WriteToLog("Info", status, nil)
}

// extractURLs pulls the URLs from the given file for image processing.
func extractURLs(inFilename string, csv *os.File) string {

	f, err := os.Open(inFilename)
	if err != nil {
		log.WriteToLog("Fatal", "URL extraction failed during setup.", err)
	}
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
	if err != nil {
		log.WriteToLog("Fatal", "Error scanning: ", err)
	}
	close(urlChan)

	wg.Wait()

	return "Process complete."
}

// imageData extracts the image from a given URL.
func imageData(url string, csv *os.File) {

	resp, err := http.Get(url)
	if err != nil {
		log.WriteToLog("Warn", "http.Get failure", err)
		return
	}
	defer resp.Body.Close()

	// Extract the image information.
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.WriteToLog("Warn", fmt.Sprintf("%v image decode error", url), err)
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
	output := heapify(timesAppeared)

	return output
}

// heapify turns the color set into a max-heap data structure
func heapify(timesAppeared map[colorCode]int) []string {

	c := make(colorHeap, 0, len(timesAppeared))

	for color, appeared := range timesAppeared {
		// Multiply by -1 since standard heap is a min-heap, this makes
		// it a max-heap.
		c = append(c, colorNode{color, appeared * -1})
	}

	h := &c
	heap.Init(h)

	return extractTopColors(h)
}

// extractTopColors pulls out the top 3 top colors in the image and
// returns them in hexadecimal format.
func extractTopColors(c *colorHeap) []string {

	topColors := make([]string, 4)
	for i := 1; i < 4; i++ {
		color := heap.Pop(c).(colorNode)
		hexColor := fmt.Sprintf("#%.2X%.2X%.2X", color.Color.Red,
			color.Color.Green, color.Color.Blue)
		topColors[i] = hexColor
	}
	return topColors
}
