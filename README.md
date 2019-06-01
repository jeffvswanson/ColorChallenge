# ColorChallenge

[Background](#background)  
[Use](#use)  
[Installation](#installation)  
[Unexplored Optimizations](#unexplored-optimizations)  
[Failed Optimizations](#failed-optimizations)  
[Contributing](#contributing)  
[Contact](#contact)  
[License](#license)  

## Background
**Given:** A list of 1000 URLs leading to images.  
**Find:** The 3 most prevalent colors in the RGB scheme in hexadecimal format (#000000 - #FFFFFF) in each image and write the result into a CSV file in a form url,color,color,color.  
**Constraints:**
- Focus on speed and resources.
- Solution should be able to handle input files with more than a billion URLs.
- Limited resources, for instance, 1 CPU and 512 MB RAM. Max out usage of resources.  

**Facts:**
- No limit on the execution time.
- Order of return of the URLs does not matter.
- Order of the top 3 colors returned does not matter.  

**Assumptions:**  
- Solution will need to be deployed in a containerized environment.
- There will always be at least 3 colors in the images.
- While there are only 40 websites repeated in the given list, assume production list will all be unique websites so treat repeated URLs as unique websites.

## Use
This program takes the URLs from the input.txt file and produces a CSV in the format of url,color,color,color.  
  
Additionally, the colorchallenge program logs any informational notifications, warnings, and errors as they occur to a dated log file in the form of time, error level, error message.  

## Installation

Installation requires Go version 1.8 or higher and can be deployed to an individual machine or a [Docker](https://www.docker.com/) container.

To run the tests:
1. Open a command line terminal
2. Change directory to the location you cloned the colorchallenge repository.
3. On the command line type: `go test ./...` This will run all tests in the repository.

### Docker Container

1. To deploy to a docker container you must have Docker installed.
2. Navigate to the downloaded repository.
3. On the command line run: `docker build -t color-challenge .`
4. Once the image has been built successfully you can run the container: `docker run -m 512m color-challenge` The option `-m` limits the maximum memory the container can use.
5. The generated container will exit when the process completes. You can check how the container is performing while running by opening another command line terminal and running: `docker container inspect container_name` where `container_name` can be found by running `docker ps -a`
6. To export the container's files at any time: 
   1. Change directories into where you want the container snapshot located.
   2. Give the command `docker export container_name > contents.tar`
   3. The `.tar` file can be unzipped with Linux tools or Windows 7zip. 

### Local Machine

1. With Go you have two options. You can run `go run main.go` when you're in the colorchallenge directory or you can run `go build .`
   1. `go run` will deposit the log and CSV output file in the current directory
   2. `go build` will allow you to place the generated executable file into a new directory and have the resulting files in the new directory.

## Future Optimizations
1. Implement a data structure that does not require mapping, but still gives constant time lookup on average. I think this could be the biggest savings given the CPU profile map.
2. Have a pre-built slice associated with each possible combination of the RGB values (16,974,593 or 8<sup>3</sup>), but consider, will need an in-memory object of about 68 MB since each int is allocated 4 bytes. Increment each index when the RGB value is found O(1), search for the top 3 values in the slice when done with image O(n). Do not sort! Indices indicate RGB value. Copy the top 3 to a new slice to release the searching slice and proceed.

## Failed Optimizations
1. Use sync.Waitgroups instead of channels. sync.Waitgroups ended up taking longer during benchmarking. 14.6 sec versus 14.1 sec.
2. Take a sample of 50% of the rows. Speed of operations reduced by about 43%, accuracy reduced from 100% to about 92.5%.
3. After profiling the program a lot of time is being spent on map assignment. Unfortunately, when attempting to reduce reliance on map assignment in search of an alternative data structure, I essentially re-implemented a map.
4. Have a set of rows fan out as goroutines. Increased concurrency increased runtime due to increased communication requirements.
5. Convert the image from YCbCr format to RGB without going pixel-by-pixel, then checking pixels. Increased runtime about 7% from pixel-by-pixel level conversion, 15.1 sec versus 14.1 sec.

## Contributing

Pull requests are welcome. For major changes, that is: add a function, change output, etc.; please open an issue to discuss what you would like to change.

Update tests as appropriate.

## Contact

<jeff.v.swanson@gmail.com>

## License

[GNU GPL v3.0](https://choosealicense.com/licenses/gpl-3.0/)