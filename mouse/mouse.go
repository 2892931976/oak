// Package mouse handles the propagation of mouse events
// though clickable regions.
// It has a lot of functions which are equivalent to those in the collision package.
package mouse

import "bitbucket.org/oakmoundstudio/oak/collision"

// Propagate triggers direct mouse events on entities which are clicked
func Propagate(eventName string, me Event) {
	mouseLoc := collision.NewUnassignedSpace(float64(me.X), float64(me.Y), 0.01, 0.01)
	hits := rt.SearchIntersect(mouseLoc.Bounds())
	for _, v := range hits {
		sp := v.(*collision.Space)
		sp.CID.Trigger(eventName, me)
	}
}

// GetMouseButton is a utitilty function which translates
// integer values of mouse keys from golang's event/mouse library
// into strings.
func GetMouseButton(i int32) (s string) {
	switch i {
	case 1:
		s = "LeftMouse"
	case 2:
		s = "MiddleMouse"
	case 3:
		s = "RightMouse"
	case -1:
		s = "ScrollUpMouse"
	case -2:
		s = "ScrollDownMouse"
	default:
		s = ""
	}
	return
}
