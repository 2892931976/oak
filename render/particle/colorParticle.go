package particle

import (
	"image"
	"image/color"
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

// A ColorParticle is a particle with a defined color and size
type ColorParticle struct {
	*baseParticle
	startColor color.Color
	endColor   color.Color
	size       int
}

// Draw redirectes to DrawOffset
func (cp *ColorParticle) Draw(buff draw.Image) {
	cp.DrawOffset(buff, 0, 0)
}

// DrawOffset redirectes to DrawOffsetGen
func (cp *ColorParticle) DrawOffset(buff draw.Image, xOff, yOff float64) {
	cp.DrawOffsetGen(cp.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

// DrawOffsetGen draws a particle with it's generator's variables
func (cp *ColorParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {
	gen := generator.(*ColorGenerator)

	r, g, b, a := cp.startColor.RGBA()
	r2, g2, b2, a2 := cp.endColor.RGBA()
	progress := cp.Life / cp.totalLife
	c := color.RGBA64{
		uint16OnScale(r, r2, progress),
		uint16OnScale(g, g2, progress),
		uint16OnScale(b, b2, progress),
		uint16OnScale(a, a2, progress),
	}

	img := image.NewRGBA64(image.Rect(0, 0, cp.size, cp.size))

	for i := 0; i < cp.size; i++ {
		for j := 0; j < cp.size; j++ {
			if gen.Shape(i, j, cp.size) {
				img.SetRGBA64(i, j, c)
			}
		}
	}

	halfSize := float64(cp.size / 2)

	render.ShinyDraw(buff, img, int((xOff+cp.Pos.X)-halfSize), int((yOff+cp.Pos.Y)-halfSize))
}

// GetPos returns the middle of a color particle
func (cp *ColorParticle) GetPos() physics.Vector {
	fSize := float64(cp.size)
	return physics.NewVector(cp.Pos.X-fSize/2, cp.Pos.Y-fSize/2)
}

// GetDims returns the color particle's size, twice
func (cp *ColorParticle) GetDims() (int, int) {
	return cp.size, cp.size
}
