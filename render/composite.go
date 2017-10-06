package render

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render/mod"
)

// Composite Types, distinct from Compound Types,
// Display all of their parts at the same time,
// and respect the positions and layers of their
// parts.
type Composite struct {
	LayeredPoint
	rs []Modifiable
}

// NewComposite creates a Composite
func NewComposite(sl ...Modifiable) *Composite {
	cs := new(Composite)
	cs.LayeredPoint = NewLayeredPoint(0, 0, 0)
	cs.rs = sl
	return cs
}

// AppendOffset adds a new offset modifiable to the composite
func (cs *Composite) AppendOffset(r Modifiable, p floatgeom.Point2) {
	r.SetPos(p.X(), p.Y())
	cs.Append(r)
}

// Append adds a renderable as is to the composite
func (cs *Composite) Append(r Modifiable) {
	cs.rs = append(cs.rs, r)
}

// Prepend adds a new renderable to the front of the CompositeR.
func (cs *Composite) Prepend(r Modifiable) {
	cs.rs = append([]Modifiable{r}, cs.rs...)
}

// SetIndex places a renderable at a certain point in the composites renderable slice
func (cs *Composite) SetIndex(i int, r Modifiable) {
	cs.rs[i] = r
}

// Len returns the number of renderables in this composite.
func (cs *Composite) Len() int {
	return len(cs.rs)
}

// AddOffset offsets all renderables in the composite by a vector
func (cs *Composite) AddOffset(i int, p floatgeom.Point2) {
	if i < len(cs.rs) {
		cs.rs[i].SetPos(p.X(), p.Y())
	}
}

// SetOffsets applies the initial offsets to the entire Composite
func (cs *Composite) SetOffsets(vs ...floatgeom.Point2) {
	for i, v := range vs {
		if i < len(cs.rs) {
			cs.rs[i].SetPos(v.X(), v.Y())
		}
	}
}

// Get returns a renderable at the given index within the composite
func (cs *Composite) Get(i int) Modifiable {
	return cs.rs[i]
}

// DrawOffset draws the Composite with some offset from its logical position (and therefore sub renderables logical positions).
func (cs *Composite) DrawOffset(buff draw.Image, xOff, yOff float64) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X()+xOff, cs.Y()+yOff)
	}
}

// Draw draws the Composite at its logical position
func (cs *Composite) Draw(buff draw.Image) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X(), cs.Y())
	}
}

// UnDraw stops the composite from being drawn
func (cs *Composite) UnDraw() {
	cs.layer = Undraw
	for _, c := range cs.rs {
		c.UnDraw()
	}
}

// GetRGBA does not work on a composite and therefore returns nil
func (cs *Composite) GetRGBA() *image.RGBA {
	return nil
}

// Modify applies mods to the composite
func (cs *Composite) Modify(ms ...mod.Mod) Modifiable {
	for _, r := range cs.rs {
		r.Modify(ms...)
	}
	return cs
}

// Filter filters each component part of this composite by all of the inputs.
func (cs *Composite) Filter(fs ...mod.Filter) {
	for _, r := range cs.rs {
		r.Filter(fs...)
	}
}

// Copy makes a new Composite with the same renderables
func (cs *Composite) Copy() Modifiable {
	cs2 := new(Composite)
	cs2.layer = cs.layer
	cs2.Vector = cs.Vector
	cs2.rs = make([]Modifiable, len(cs.rs))
	for i, v := range cs.rs {
		cs2.rs[i] = v.Copy()
	}
	return cs2
}

// CompositeR keeps track of a set of renderables at a location
type CompositeR struct {
	LayeredPoint
	toPush []Renderable
	rs     []Renderable
}

// NewCompositeR creates a new CompositeR from a slice of renderables
func NewCompositeR(sl ...Renderable) *CompositeR {
	cs := new(CompositeR)
	cs.LayeredPoint = NewLayeredPoint(0, 0, 0)
	cs.toPush = make([]Renderable, 0)
	cs.rs = sl
	return cs
}

// AppendOffset adds a new renderable to CompositeR with an offset
func (cs *CompositeR) AppendOffset(r Renderable, p floatgeom.Point2) {
	r.SetPos(p.X(), p.Y())
	cs.Append(r)
}

// AddOffset adds an offset to a given renderable of the slice
func (cs *CompositeR) AddOffset(i int, p floatgeom.Point2) {
	if i < len(cs.rs) {
		cs.rs[i].SetPos(p.X(), p.Y())
	}
}

// Append adds a new renderable to the end of the CompositeR.
func (cs *CompositeR) Append(r Renderable) {
	cs.rs = append(cs.rs, r)
}

// Prepend adds a new renderable to the front of the CompositeR.
func (cs *CompositeR) Prepend(r Renderable) {
	cs.rs = append([]Renderable{r}, cs.rs...)
}

// Len returns the number of renderables in this composite.
func (cs *CompositeR) Len() int {
	return len(cs.rs)
}

// SetIndex places a renderable at a certain point in the composites renderable slice
func (cs *CompositeR) SetIndex(i int, r Renderable) {
	cs.rs[i] = r
}

// SetOffsets sets all renderables in CompositeR to the passed in Vector positions positions
func (cs *CompositeR) SetOffsets(ps ...floatgeom.Point2) {
	for i, p := range ps {
		if i < len(cs.rs) {
			cs.rs[i].SetPos(p.X(), p.Y())
		}
	}
}

// DrawOffset Draws the CompositeR with an offset from its logical location.
func (cs *CompositeR) DrawOffset(buff draw.Image, xOff, yOff float64) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X()+xOff, cs.Y()+yOff)
	}
}

// Draw draws the CompositeR at its logical location and therefore its consituent renderables as well
func (cs *CompositeR) Draw(buff draw.Image) {
	for _, c := range cs.rs {
		c.DrawOffset(buff, cs.X(), cs.Y())
	}
}

// UnDraw undraws the CompositeR and its consituent renderables
func (cs *CompositeR) UnDraw() {
	cs.layer = Undraw
	for _, c := range cs.rs {
		c.UnDraw()
	}
}

// GetRGBA always returns nil from Composites
func (cs *CompositeR) GetRGBA() *image.RGBA {
	return nil
}

// Get returns renderable from a given index in CompositeR
func (cs *CompositeR) Get(i int) Renderable {
	return cs.rs[i]
}

// Add stages a renderable to be added to the Composite at the next PreDraw
func (cs *CompositeR) Add(r Renderable, _ ...int) Renderable {
	cs.toPush = append(cs.toPush, r)
	return r
}

// Replace updates a renderable in the CompositeR to the new Renderable
func (cs *CompositeR) Replace(r1, r2 Renderable, i int) {
	cs.Add(r2, i)
	r1.UnDraw()
}

// PreDraw updates the CompositeR with the new renderables to add. This helps keep consistency and mitigates the threat of unsafe operations.
func (cs *CompositeR) PreDraw() {
	push := cs.toPush
	cs.toPush = []Renderable{}
	cs.rs = append(cs.rs, push...)
}

// Copy returns a new composite with the same length slice of renderables but no actual renderables...
// CompositeRs cannot have their internal elements copied,
// as renderables cannot be copied.
func (cs *CompositeR) Copy() Stackable {
	cs2 := new(CompositeR)
	cs2.LayeredPoint = cs.LayeredPoint
	cs2.rs = make([]Renderable, len(cs.rs))
	return cs2
}

func (cs *CompositeR) draw(world draw.Image, viewPos image.Point, screenW, screenH int) {
	realLength := len(cs.rs)
	for i := 0; i < realLength; i++ {
		r := cs.rs[i]
		for (r == nil || r.GetLayer() == Undraw) && realLength > i {
			cs.rs[i], cs.rs[realLength-1] = cs.rs[realLength-1], cs.rs[i]
			realLength--
			r = cs.rs[i]
		}
		if realLength == i {
			break
		}
		x := int(r.X())
		y := int(r.Y())
		x2 := x
		y2 := y
		w, h := r.GetDims()
		x += w
		y += h
		if x > viewPos.X && y > viewPos.Y &&
			x2 < viewPos.X+screenW && y2 < viewPos.Y+screenH {

			if InDrawPolygon(x, y, x2, y2) {
				r.DrawOffset(world, float64(-viewPos.X), float64(-viewPos.Y))
			}
		}
	}
	cs.rs = cs.rs[0:realLength]
}
