package render

import (
	"errors"
	"image"
	"image/draw"
	"sync"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render/mod"
)

// The Compound type is intended for use to easily swap between multiple
// renderables that are drawn at the same position on the same layer.
// A common use case for this would be a character entitiy who switches
// their animation based on how they are moving or what they are doing.
//
// The Compound type removes the need to repeatedly draw and undraw elements
// of a character, which has a tendency to leave nothing drawn for a draw frame.
type Compound struct {
	LayeredPoint
	subRenderables map[string]Modifiable
	curRenderable  string
	lock           sync.RWMutex
}

// NewCompound creates a new compound from a map of names to modifiables
func NewCompound(start string, m map[string]Modifiable) *Compound {
	return &Compound{
		LayeredPoint:   NewLayeredPoint(0, 0, 0),
		subRenderables: m,
		curRenderable:  start,
		lock:           sync.RWMutex{},
	}
}

// Add makes a new entry in the Compounds map
func (c *Compound) Add(k string, v Modifiable) (err error) {
	if _, ok := c.subRenderables[k]; ok {
		err = errors.New("Key already defined. Overwriting")
	}
	c.lock.Lock()
	c.subRenderables[k] = v
	c.lock.Unlock()
	return err
}

// Set sets the current renderable to the one specified
func (c *Compound) Set(k string) error {
	c.lock.RLock()
	if _, ok := c.subRenderables[k]; !ok {
		return errors.New("Unknown renderable for string " + k + " on compound")
	}
	c.lock.RUnlock()
	c.curRenderable = k
	return nil
}

// GetSub returns a keyed Modifiable from this compound's map
func (c *Compound) GetSub(s string) Modifiable {
	c.lock.RLock()
	m := c.subRenderables[s]
	c.lock.RUnlock()
	return m
}

// Get returns the Compound's current key
func (c *Compound) Get() string {
	return c.curRenderable
}

// SetOffsets sets the logical offset for the specified key
func (c *Compound) SetOffsets(k string, offsets physics.Vector) {
	c.lock.RLock()
	if r, ok := c.subRenderables[k]; ok {
		r.SetPos(offsets.X(), offsets.Y())
	}
	c.lock.RUnlock()
}

// Copy creates a copy of the Compound
func (c *Compound) Copy() Modifiable {
	newC := new(Compound)
	newC.LayeredPoint = c.LayeredPoint.Copy()
	newSubRenderables := make(map[string]Modifiable)
	c.lock.RLock()
	for k, v := range c.subRenderables {
		newSubRenderables[k] = v.Copy()
	}
	c.lock.RUnlock()
	newC.subRenderables = newSubRenderables
	newC.curRenderable = c.curRenderable
	newC.lock = sync.RWMutex{}
	return newC
}

//GetRGBA returns the current renderables rgba
func (c *Compound) GetRGBA() *image.RGBA {
	c.lock.RLock()
	rgba := c.subRenderables[c.curRenderable].GetRGBA()
	c.lock.RUnlock()
	return rgba
}

// Modify performs the input modifications on all elements of the Compound
func (c *Compound) Modify(ms ...mod.Mod) Modifiable {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		r.Modify(ms...)
	}
	c.lock.RUnlock()
	return c
}

// Filter filters all elements of the compound with fs
func (c *Compound) Filter(fs ...mod.Filter) {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		r.Filter(fs...)
	}
	c.lock.RUnlock()
}

//DrawOffset draws the Compound at an offset from its logical location
func (c *Compound) DrawOffset(buff draw.Image, xOff float64, yOff float64) {
	c.lock.RLock()
	c.subRenderables[c.curRenderable].DrawOffset(buff, c.X()+xOff, c.Y()+yOff)
	c.lock.RUnlock()
}

//Draw draws the Compound at its logical location
func (c *Compound) Draw(buff draw.Image) {
	c.lock.RLock()
	c.subRenderables[c.curRenderable].DrawOffset(buff, c.X(), c.Y())
	c.lock.RUnlock()
}

// ShiftPos shifts the Compounds logical position
func (c *Compound) ShiftPos(x, y float64) {
	c.SetPos(c.X()+x, c.Y()+y)
}

// GetDims gets the current Renderables dimensions
func (c *Compound) GetDims() (int, int) {
	c.lock.RLock()
	w, h := c.subRenderables[c.curRenderable].GetDims()
	c.lock.RUnlock()
	return w, h
}

// Pause stops the current Renderable if possible
func (c *Compound) Pause() {
	c.lock.RLock()
	if cp, ok := c.subRenderables[c.curRenderable].(CanPause); ok {
		cp.Pause()
	}
	c.lock.RUnlock()
}

// Unpause tries to unpause the current Renderable if possible
func (c *Compound) Unpause() {
	c.lock.RLock()
	if cp, ok := c.subRenderables[c.curRenderable].(CanPause); ok {
		cp.Unpause()
	}
	c.lock.RUnlock()
}

// IsInterruptable returns whether the current renderable is interruptable
func (c *Compound) IsInterruptable() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if i, ok := c.subRenderables[c.curRenderable].(NonInterruptable); ok {
		return i.IsInterruptable()
	}
	return true
}

// IsStatic returns whether the current renderable is static
func (c *Compound) IsStatic() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if s, ok := c.subRenderables[c.curRenderable].(NonStatic); ok {
		return s.IsStatic()
	}
	return true
}

// SetTriggerID sets the ID AnimationEnd will trigger on for animating subtypes.
// Todo: standardize this with the other interface Set functions so that it
// also only acts on the current subRenderable, or the other way around, or
// somehow offer both options
func (c *Compound) SetTriggerID(cid event.CID) {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		if t, ok := r.(Triggerable); ok {
			t.SetTriggerID(cid)
		}
	}
	c.lock.RUnlock()
}

// Revert will revert all parts of this compound that can be reverted
func (c *Compound) Revert(mod int) {
	c.lock.RLock()
	for _, v := range c.subRenderables {
		switch t := v.(type) {
		case *Reverting:
			t.Revert(mod)
		}
	}
	c.lock.RUnlock()
}

// RevertAll will revert all parts of this compound that can be reverted, back
// to their original state.
func (c *Compound) RevertAll() {
	c.lock.RLock()
	for _, v := range c.subRenderables {
		switch t := v.(type) {
		case *Reverting:
			t.RevertAll()
		}
	}
	c.lock.RUnlock()
}

func (c *Compound) update() {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		if u, ok := r.(updates); ok {
			u.update()
		}
	}
	c.lock.RUnlock()
}
