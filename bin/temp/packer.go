package main

import (
	"compress/flate"
	"io"
	"os"

	//"github.com/klauspost/compress/zlib"
)

func main() {
	// Open the input binary file
	inputFile, err := os.Open("shell")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	// Create the output packed file
	outputFile, err := os.Create("shell_packed")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Create a flate compressor
	compressor, err := flate.NewWriter(outputFile, flate.BestCompression)
	if err != nil {
		panic(err)
	}
	defer compressor.Close()

	// Copy the input file to the compressor
	_, err = io.Copy(compressor, inputFile)
	if err != nil {
		panic(err)
	}

	// Flush and close the compressor
	err = compressor.Flush()
	if err != nil {
		panic(err)
	}
}
