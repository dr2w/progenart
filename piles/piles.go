// Package piles implements basic functions around
// abelian sandpiles using a 2D Image.
package piles

import (
	"image"
	"image/color"
	"log"
)

// Connectivity defines the number of "neighbors" each grid space has.
type Connectivity int

const (
	Four Connectivity = iota
	Eight
)

// deltaMap maps Connectivity to a set of offsets used to generate neighbors.
var deltaMap = map[Connectivity][]image.Point{
	Four:  []image.Point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}},
	Eight: []image.Point{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}},
}

// contains returns true iff the given point falls within the given rectangle.
func contains(r image.Rectangle, p image.Point) bool {
	return p.X >= r.Min.X && p.X < r.Max.X && p.Y >= r.Min.Y && p.Y < r.Max.Y
}

// Config exposes the Wrap and Connectivity options, and is used to store deltas
// and height during computation.
type Config struct {
	Wrap         bool
	Connectivity Connectivity
	deltas       []image.Point
	height       uint16
}

// Deltas caches and returns the set of deltas required to generate all neighbors.
func (c *Config) Deltas() []image.Point {
	if len(c.deltas) == 0 {
		var ok bool
		c.deltas, ok = deltaMap[c.Connectivity]
		if !ok {
			log.Fatalf("Invalid connectivity (%d) in piles.Config", c.Connectivity)
		}
	}
	return c.deltas
}

// Height caches and returns the height which will trigger an overflow event.
func (c *Config) Height() uint16 {
	if c.height == 0 {
		c.height = uint16(len(c.Deltas()))
	}
	return c.height
}

// spill takes an image and a point, and performs the basic spill operation if
// necessary. It returns the potentially modified image, along with a boolean
// indicating whether or not the image was modified.
func (c *Config) spill(p image.Point, img *image.Gray16) (*image.Gray16, bool) {
	bounds := img.Bounds()
	height := c.Height()
	v, _, _, _ := img.At(p.X, p.Y).RGBA()
	value := uint16(v)
	if value < height {
		return img, false
	}
	img.SetGray16(p.X, p.Y, color.Gray16{value - height})
	for _, d := range c.Deltas() {
		nbr := p.Add(d)
		if !contains(bounds, nbr) {
			if !c.Wrap {
				continue
			}
			if nbr.X < bounds.Min.X {
				nbr.X = bounds.Max.X - 1
			}
			if nbr.X >= bounds.Max.X {
				nbr.X = bounds.Min.X
			}
			if nbr.Y < bounds.Min.Y {
				nbr.Y = bounds.Max.Y - 1
			}
			if nbr.Y >= bounds.Max.Y {
				nbr.Y = bounds.Min.Y
			}
		}
		nbrValue, _, _, _ := img.At(nbr.X, nbr.Y).RGBA()
		img.SetGray16(nbr.X, nbr.Y, color.Gray16{uint16(nbrValue) + 1})
	}
	return img, true
}

// step performs a single "spill iteration" on the given image, spilling each
// point once and returning the modified image, along with a boolean indicating
// whether or not any modifications were required.
func (c *Config) step(img *image.Gray16) (*image.Gray16, bool) {
	var changed bool
	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			i, changes := c.spill(image.Pt(x, y), img)
			changed = changed || changes
			img = i
		}
	}
	return img, changed
}

// Resolve takes an image and steps until no changes have been made,then returns
// the final image. Math says that this should always converge.
func (c *Config) Resolve(img *image.Gray16) *image.Gray16 {
	changed := true
	count := 0
	for changed {
		img, changed = c.step(img)
		count++
		if count%10 == 0 {
			log.Printf("Iteration %d", count)
		}
	}
	return img
}
