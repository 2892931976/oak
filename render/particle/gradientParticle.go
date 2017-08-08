package particle

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/render"
)

// A GradientParticle has a gradient from one color to another
type GradientParticle struct {
	ColorParticle
	startColor2 color.Color
	endColor2   color.Color
}

// Draw redirectes to DrawOffset
func (gp *GradientParticle) Draw(buff draw.Image) {
	gp.DrawOffset(buff, 0, 0)
}

// DrawOffset redirectes to DrawOffsetGen
func (gp *GradientParticle) DrawOffset(buff draw.Image, xOff, yOff float64) {
	gp.DrawOffsetGen(gp.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

// DrawOffsetGen draws a particle with it's generator's variables
func (gp *GradientParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {

	gen := generator.(*GradientGenerator)
	progress := gp.Life / gp.totalLife
	c1 := render.GradientColorAt(gp.startColor, gp.endColor, progress)
	c2 := render.GradientColorAt(gp.startColor2, gp.endColor2, progress)

	img := image.NewRGBA64(image.Rect(0, 0, gp.size, gp.size))

	for i := 0; i < gp.size; i++ {
		for j := 0; j < gp.size; j++ {
			if gen.Shape.In(i, j, gp.size) {
				progress := gen.ProgressFunction(i, j, gp.size, gp.size)
				c := render.GradientColorAt(c1, c2, progress)
				img.SetRGBA64(i, j, c)
			}
		}
	}

	halfSize := float64(gp.size / 2)

	render.ShinyDraw(buff, img, int((xOff+gp.X())-halfSize), int((yOff+gp.Y())-halfSize))
}
