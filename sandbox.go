package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	_ "image/jpeg"
)

var (
	source     = flag.String("source", "", "Source image for the computation")
	outputFile = flag.String("output_file", "", "Write the output file to this location")
)

// load reads in the given file as an image.
// jpeg and png are currently supported.
func load(fn string) (image.Image, error) {
	reader, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// writePNG writes the given image to the given filename in PNG format.
func writePNG(img image.Image, fn string) error {
	w, err := os.Create(fn)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}

// toGray converts the given image to Gray16
func toGray(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	gray := image.NewGray16(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c := img.At(x, y)
			gray.Set(x, y, color.Gray16Model.Convert(c))
		}
	}
	return gray, nil
}

// process converts the input source image into a final output image.
func process(img image.Image) (image.Image, error) {
	gray, err := toGray(img)
	if err != nil {
		return nil, err
	}
	return gray
}

// main simply calls load -> process -> write.
func main() {
	flag.Parse()

	img, err := load(*source)
	if err != nil {
		log.Fatalf("unable to load %s:\n%s", *source, err)
	}

	img, err = process(img)
	if err != nil {
		log.Fatalf("unable to process image:\n%s", err)
	}

	if err = writePNG(img, *outputFile); err != nil {
		log.Fatalf("unable to write image %s:\n%s", *outputFile, err)
	}
}
