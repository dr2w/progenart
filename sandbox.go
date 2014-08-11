package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"dr2w.com/progenart/piles"

	_ "image/jpeg"
)

var (
	source     = flag.String("source", "", "Source image for the computation")
	outputFile = flag.String("output_file", "", "Write the output file to this location")
	width      = flag.Int("width", 1024, "Image width to use when not loading an existing image.")
	height     = flag.Int("height", 768, "Image height to use when not loading an existing image.")
	grains     = flag.Int("grains", 100*1000, "The number of grains to drop onto the image.")
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
func toGray(img image.Image) (*image.Gray16, error) {
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
	bounds := img.Bounds()
	mid := bounds.Max.Sub(bounds.Min).Div(2)
	gray.SetGray16(mid.X, mid.Y, color.Gray16{uint16(*grains)})
	config := &piles.Config{
		Wrap:         false,
		Connectivity: piles.Four,
	}
	config.Resolve(gray)
	return amplify(gray, 10000), nil
}

// amplify multiples the pixels of a Gray16 by the given factor.
func amplify(img *image.Gray16, factor uint16) *image.Gray16 {
	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c, ok := img.At(x, y).(color.Gray16)
			if !ok {
				log.Fatalf("Bad color cast in amplify")
			}
			img.Set(x, y, color.Gray16{c.Y * factor})
		}
	}
	return img
}

// main simply calls load -> process -> write.
func main() {
	flag.Parse()

	var img image.Image
	if *source != "" {
		var err error
		img, err = load(*source)
		if err != nil {
			log.Fatalf("unable to load %s:\n%s", *source, err)
		}
	} else {
		img = image.NewGray16(image.Rect(0, 0, *width, *height))
	}

	result, err := process(img)
	if err != nil {
		log.Fatalf("unable to process image:\n%s", err)
	}

	if err = writePNG(result, *outputFile); err != nil {
		log.Fatalf("unable to write image %s:\n%s", *outputFile, err)
	}
}
