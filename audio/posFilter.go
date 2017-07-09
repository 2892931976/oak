package audio

import (
	"fmt"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/audio/filter"
	"github.com/200sc/klangsynthese/audio/filter/supports"
)

type SupportsPos interface {
	supports.Encoding
	GetX() *float64
	GetY() *float64
}

var (
	_ audio.Filter = Pos(func(SupportsPos) {})
)

type Pos func(SupportsPos)

func (xp Pos) Apply(a audio.Audio) (audio.Audio, error) {
	fmt.Println(a)
	if sxp, ok := a.(SupportsPos); ok {
		xp(sxp)
		return a, nil
	}
	fmt.Println("Doesn't support SupportsPos")
	return a, supports.NewUnsupported([]string{"XPan"})
}

func PosFilter(e *Ears) Pos {
	return func(sp SupportsPos) {
		fmt.Println("PosFilter", e, sp)
		x, y := sp.GetX(), sp.GetY()
		if x != nil {
			filter.Pan(e.CalculatePan(*x))
			if y != nil {
				filter.Volume(e.CalculateVolume(physics.NewVector(*x, *y)))
			}
		}
	}
}
