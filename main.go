package main

import (
	"bufio"
	"container/heap"
	"errors"
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

type imageInfo struct {
	p   *http.Response
	URL string
}

type colorNode struct {
	Color       color.Color
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

var logfile, csvfile *os.File

func init() {
	logfile = log.FormatLog()

	csvfile = exporttocsv.CreateCSV("ColorChallengeOutput")
	headerRecord := []string{"URL", "top_color1", "top_color2", "top_color3"}
	exporttocsv.Export(csvfile, headerRecord)
}

func main() {

	defer logfile.Close()
	defer csvfile.Close()
	inputFilename := "input_test.txt"

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

	var wg sync.WaitGroup

	urlChan := make(chan string)
	defer close(urlChan)

	images := make(chan imageInfo)
	defer close(images)

	// Spawn workers to prevent saturating bandwidth.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for url := range urlChan {
				extractImageData(url, images)
			}
			wg.Done()
		}()
	}

	// While there may only be 1 processor, maybe we'll get lucky.
	workers := runtime.GOMAXPROCS(-1)

	// Spawn workers to prevent running out of memory.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			for image := range images {
				countColors(image, csv)
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

	wg.Wait()

	return "Process complete."
}

// extractImageData extracts the image from a given URL.
func extractImageData(url string, images chan<- imageInfo) {

	resp, err := http.Get(url)
	if err != nil {
		log.WriteToLog("Warn", "http.Get failure - ", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.WriteToLog("Warn", fmt.Sprintf("%s http status not ok", url), errors.New(resp.Status))
		return
	}
	img := imageInfo{
		resp,
		url,
	}

	images <- img
}

// countColors finds pixel color mapping of an image in RGB format.
func countColors(i imageInfo, csv *os.File) {

	// Extract the image information.
	defer i.p.Body.Close()
	img, _, err := image.Decode(i.p.Body)
	if err != nil {
		log.WriteToLog("Warn", fmt.Sprintf("%s image decode error", i.URL), err)
		return
	}

	timesAppeared := make(map[color.Color]int)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			// img returns in YCbCr format will convert to RGB 8-bit on top 3 color return
			yCbCr := img.At(x, y).(color.Color)
			timesAppeared[yCbCr.(color.Color)]++
		}
	}
	// Sort colors in descending order.
	output := heapify(timesAppeared)

	// Get the output string into url,color,color,color format.
	output[0] = i.URL
	exporttocsv.Export(csv, output)
}

// heapify turns the color set into a max-heap data structure
func heapify(timesAppeared map[color.Color]int) []string {

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
func extractTopColors(h *colorHeap) []string {

	topColors := make([]string, 4)
	for i := 1; i < 4; i++ {
		c := heap.Pop(h).(colorNode)
		// Convert the YCbCr format to RGB
		rgb := color.NRGBAModel.Convert(c.Color).(color.NRGBA)
		hexColor := fmt.Sprintf("#%.2X%.2X%.2X", rgb.R, rgb.G, rgb.B)
		topColors[i] = hexColor
	}
	return topColors
}
