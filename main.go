package main

import (
	"github.com/nfnt/resize"
	"image"
	"log"
	"os"
)

func main() {
	sampleFile := SampleGIF
	width := uint(140)
	height := uint(140)

	file, _ := sampleFile.Open("sample.gif")
	tx := func(m image.Image) image.Image {
		return resize.Resize(width, height, m, resize.NearestNeighbor)
	}

	f, err := os.Create("result.gif")
	if err != nil {
		log.Fatal(err)
	}
	err = processImage(f, file, tx)
	if err != nil {
		log.Fatalf("error processing image: %v", err)
	}

}
