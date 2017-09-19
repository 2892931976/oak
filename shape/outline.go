package shape

import (
	"errors"

	"github.com/oakmound/oak/alg/intgeom"
)

const (
	top = iota
	topright
	right
	bottomright
	bottom
	bottomleft
	left
	topleft
	lastdirection
)

var (
	xyMods = []int{
		0, -1,
		1, -1,
		1, 0,
		1, 1,
		0, 1,
		-1, 1,
		-1, 0,
		-1, -1,
	}
	pointDeltas = []int{
		1, 0,
		0, 1,
		0, 1,
		-1, 0,
		-1, 0,
		0, -1,
		0, -1,
		1, 0,
	}
)

func outlineSizes(sizes ...int) (w, h int) {
	w = sizes[0]
	h = sizes[0]
	if len(sizes) > 1 {
		h = sizes[1]
	}
	return
}

// ToOutline returns the set of points along the input shape's outline, if
// one exists.
func ToOutline(shape Shape) func(...int) ([]intgeom.Point, error) {
	return func(sizes ...int) ([]intgeom.Point, error) {
		w, h := outlineSizes(sizes...)

		//First decrement on diagonal to find start of outline
		startX := 0.0
		startY := 0.0
		fw := float64(w)
		fh := float64(h)
		maxDim := fw
		if h > w {
			maxDim = fh
		}
		xDelta := fw / maxDim
		yDelta := fh / maxDim
		for !shape.In(int(startX), int(startY), sizes...) {
			startX += xDelta
			startY += yDelta
			if startX >= fw || startY >= fh {
				return []intgeom.Point{}, errors.New("Could not find any valid space on the shapes diagonal... Assuming that it is not valid for outlines")
			}
		}

		for startY >= 0 && shape.In(int(startX), int(startY), sizes...) {
			startY--
		}
		startY++

		//Here we have found a point on the outline
		sx := int(startX)
		sy := int(startY)
		x := sx
		y := sy

		outline := []intgeom.Point{intgeom.NewPoint(x, y)}

		direction := topright

		x += xyMods[direction*2]
		y += xyMods[direction*2+1]

		for direction != top && !inOutline(shape, x, y, w, h) {
			x += pointDeltas[direction*2]
			y += pointDeltas[direction*2+1]
			direction = (direction + 1) % lastdirection
		}
		if direction == top {
			return outline, nil
		}

		//Follow the outline point by point
		for x != sx || y != sy {
			outline = append(outline, intgeom.NewPoint(x, y))
			direction -= 2
			if direction < 0 {
				direction += lastdirection
			}
			x += xyMods[direction*2]
			y += xyMods[direction*2+1]
			//From a point on the outline look clockwise around for next direction
			for !inOutline(shape, x, y, w, h) {
				x += pointDeltas[direction*2]
				y += pointDeltas[direction*2+1]
				direction = (direction + 1) % lastdirection
			}
		}

		return outline, nil
	}
}

func inOutline(s Shape, x, y, w, h int) bool {
	return (x < w && x >= 0 && y < h && y >= 0) && s.In(x, y, w, h)
}
