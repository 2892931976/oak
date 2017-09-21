package particle

import (
	"image/draw"

	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

// A SpriteParticle is a particle that has an amount of sprite rotation
type SpriteParticle struct {
	*baseParticle
	rotation float32
}

// Draw redirectes to DrawOffset
func (sp *SpriteParticle) Draw(buff draw.Image) {
	sp.DrawOffset(buff, 0, 0)
}

// DrawOffset redirects to DrawOffsetGen
func (sp *SpriteParticle) DrawOffset(buff draw.Image, xOff, yOff float64) {
	sp.DrawOffsetGen(sp.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

// DrawOffsetGen draws a particle with it's generator's variables
func (sp *SpriteParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {

	sp.rotation += sp.rotation
	gen := generator.(*SpriteGenerator)
	rgba := gen.Base.Copy().Modify(mod.Rotate(sp.rotation)).GetRGBA()
	render.ShinyDraw(buff, rgba, int(sp.X()+xOff), int(sp.Y()+yOff))
}

// GetParticleSize returns the size of the sprite that the generator generates
func (sg *SpriteGenerator) GetParticleSize() (float64, float64, bool) {

	bounds := sg.Base.GetRGBA().Rect.Max

	return float64(bounds.X), float64(bounds.Y), false
}
