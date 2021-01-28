package main

import (
	"flag"
	"fmt"
	"image"
	"lego-mosaic/cmd/mosaic"
	"os"
	"time"
)

func main() {
	fromFile := flag.String("from", "examples/ll.jpg", "path to file")
	toFile := flag.String("to", "./mosaic.jpeg", "save file")
	paletteFile := flag.String("palette", "./examples/legoPalette.json", "palette file")
	size := flag.Int("size", 1, "image size")
	enableDithering := flag.Bool("dithering", true, "enable quantizer.go")
	flag.Parse()

	fmt.Printf("Loading image %s\n", *fromFile)
	originalImage, err := mosaic.LoadImage(*fromFile)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loading palette %s...\n", *paletteFile)
	mosaicPalette, err := mosaic.LoadColorPalette(*paletteFile)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Resizing image...\n")
	err, resizedImage := mosaic.ResizeImage(originalImage, image.Point{
		X: *size,
		Y: *size,
	})
	if (err != nil) {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	now := time.Now()
	fmt.Printf("Generating mosaic...\n")
	lMosaic, err := mosaic.QuantizeMosaic(mosaic.NewLegoMosaic(resizedImage, mosaicPalette), *enableDithering)
	since := time.Since(now)

	if (err != nil) {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	fmt.Printf("Saving mosaic as %s...\n", *toFile)
	aaa := lMosaic.ToImage()
	err = mosaic.SaveImage(*toFile, aaa)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	fmt.Printf("Render took %.2fs\n", since.Seconds())
	fmt.Printf("Number of lego parts: %d\n", lMosaic.LegoPartsCount())
	fmt.Printf("\nLego parts that needs to be used:\n")

	for s, i := range lMosaic.LegoParts {
		fmt.Printf("Color \"%s\" = %d\n", s, i)
	}

	os.Exit(0)
}
