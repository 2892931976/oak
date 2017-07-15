package shape

import (
	"errors"

	"bitbucket.org/oakmoundstudio/oak/alg/intgeom"
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

// ToOutline returns the set of points along the input shape's outline, if
// one exists.
func ToOutline(shape Shape) func(...int) ([]intgeom.Point, error) {
	return func(sizes ...int) ([]intgeom.Point, error) {
		w := sizes[0]
		h := sizes[0]
		if len(sizes) > 1 {
			h = sizes[1]
		}
		//TODO: use width and height deltas so that more shapes are valid
		//First decrement on diagonal to find start of outline
		startX := 0
		for !shape.In(startX, startX, sizes...) {
			startX++
			if startX == w || startX == h {
				return []intgeom.Point{}, errors.New("Could not find any valid space on the shapes diagonal... Assuming that it is not valid for outlines")
			}
		}

		startY := startX
		for startY >= 0 && shape.In(startX, startY, sizes...) {
			startY--
		}
		startY++

		//Here we have found a point on the outline
		x := startX
		y := startY

		outline := []intgeom.Point{intgeom.NewPoint(startX, startY)}

		direction := topright

		x += xyMods[direction*2]
		y += xyMods[direction*2+1]

		for direction != top &&
			(!inBounds(x, y, w, h) ||
				!shape.In(x, y, sizes...)) {

			x += pointDeltas[direction*2]
			y += pointDeltas[direction*2+1]
			direction = (direction + 1) % lastdirection
		}
		if direction == top {
			return outline, nil
		}

		//Follow the outline point by point
		for x != startX || y != startY {
			outline = append(outline, intgeom.NewPoint(x, y))
			direction -= 2
			if direction < 0 {
				direction += lastdirection
			}
			x += xyMods[direction*2]
			y += xyMods[direction*2+1]
			//From a point on the outline look clockwise around for next direction
			for !inBounds(x, y, w, h) ||
				!shape.In(x, y, sizes...) {
				x += pointDeltas[direction*2]
				y += pointDeltas[direction*2+1]
				direction = (direction + 1) % lastdirection
			}
		}

		return outline, nil
	}
}

func inBounds(x, y, w, h int) bool {
	return x < w && x >= 0 && y < h && y >= 0
}
