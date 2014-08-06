package img

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"strconv"
)

func NewFromStrings(strings []string) *image.Gray16 {
	height := len(strings)
	width := len(strings[0])
	img := image.NewGray16(image.Rect(0, 0, width, height))
	for y, row := range strings {
		for x := range row {
			v, err := strconv.Atoi(row[x : x+1])
			if err != nil {
				return nil
			}
			img.SetGray16(x, y, color.Gray16{uint16(v)})
		}
	}
	return img
}

func ToSimpleString(img *image.Gray16) string {
	var buffer bytes.Buffer
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			v, _, _, _ := img.At(x, y).RGBA()
			buffer.WriteString(fmt.Sprintf("%d", v))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
