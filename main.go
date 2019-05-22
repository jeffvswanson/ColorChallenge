package main

import (
	"bufio"
	"colorchallenge/pkg/errorlogging"
	"colorchallenge/pkg/exporttocsv"
	"fmt"
	"os"
)

type rgb struct {
	Red, Green, Blue int
}

func main() {

	status := "Beginning setup."
	fmt.Println(status)
	status = csvSetup()
	fmt.Println(status)
	imgColorPrevalence, status := extractURLs("input.txt")

	_, ok := imgColorPrevalence["http://i.imgur.com/FApqk3D.jpg"]
	if !ok {
		fmt.Println("Key not found.")
	} else {
		fmt.Println("Key found.")
	}

	fmt.Println(status)
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
Q: Are they interested in the whole image, that is, including background color, or do we want the dominant focus of the image?

Q: The requirement is for the three most prevalent colors, do they have to be in any sort of order or just list the three color
hexadecimal color codes?

Q: Do the urls have to stay in the same order of appearance in the CSV file as the input.txt?
*/

/*
Ideas:

1. Make sections of the code supporting packages. For example, not all the code needs to be in one main file. The CSV handler could be a package and called into main.

2. Nothing says I'm explicitly limited, just that I may be limited.

3. Given list 1000 urls to an image, to simulate a larger number keep looping around the list. Will this cause a denial of service?

4. Take a wide sample, say 1000 pixels apart. If the pixels are the same value assume all pixels have the same value in between. If not, cut
the sample in half to find where the pixels would be the same.

5. Create a struct to hold RGB values.
	type RGB struct {R int, G int, B int}
*/

/*

 */

func csvSetup() string {

	filename := exporttocsv.CreateCSV("ColorChallengeOutput")
	headerRecord := []string{"URL", "top_color1", "top_color2", "top_color3"}
	exporttocsv.Export(filename, headerRecord)

	return "CSV setup complete."
}

func extractURLs(filename string) (map[string]map[rgb]int, string) {
	f, err := os.Open(filename)
	errorlogging.ErrorCheck("URL extraction failed during setup.", err)
	defer f.Close()

	// Continue to think on data structure
	// URL is the key
	// URL represents an image with RGB color codes
	// Color codes are a key
	// The number of times a color code appears
	imgColorPrevalence := make(map[string]map[rgb]int)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		imgColorPrevalence[scanner.Text()] = make(map[rgb]int)
	}

	return imgColorPrevalence, "URLs extracted."
}
