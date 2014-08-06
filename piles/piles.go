// Package piles implements basic functions around
// abelian sandpiles using a 2D Image.
package piles

import (
	"image"
	"image/color"
)

type Connectivity int

const (
	Four Connectivity = iota
	Eight
)

var deltaMap = map[Connectivity][]image.Point{
	Four:  []image.Point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}},
	Eight: []image.Point{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}},
}

func contains(r image.Rectangle, p image.Point) bool {
	return p.X >= r.Min.X && p.X < r.Max.X && p.Y >= r.Min.Y && p.Y < r.Max.Y
}

type Config struct {
	Wrap         bool
	Connectivity Connectivity
	deltas       []image.Point
	height       uint16
}

func (c *Config) Deltas() []image.Point {
	if len(c.deltas) == 0 {
		c.deltas = deltaMap[c.Connectivity]
	}
	return c.deltas
}

func (c *Config) Height() uint16 {
	if c.height == 0 {
		c.height = uint16(len(c.Deltas()))
	}
	return c.height
}

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
			switch {
			case nbr.X < bounds.Min.X:
				nbr.X = bounds.Max.X - 1
			case nbr.X >= bounds.Max.X:
				nbr.X = bounds.Min.X
			case nbr.Y < bounds.Min.Y:
				nbr.Y = bounds.Max.Y - 1
			case nbr.Y >= bounds.Max.Y:
				nbr.Y = bounds.Min.Y
			}
		}
		nbrValue, _, _, _ := img.At(nbr.X, nbr.Y).RGBA()
		img.SetGray16(nbr.X, nbr.Y, color.Gray16{uint16(nbrValue) + 1})
	}
	return img, true
}

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
