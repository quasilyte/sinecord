package stage

import (
	"github.com/quasilyte/gmath"
)

type pos32 struct {
	x float32
	y float32
}

func toPos32(p gmath.Vec) pos32 {
	return pos32{x: float32(p.X), y: float32(p.Y)}
}

func rotate(x, y float32, angle gmath.Rad) pos32 {
	sine := float32(angle.Sin())
	cosi := float32(angle.Cos())
	return pos32{x*cosi - y*sine, x*sine + y*cosi}
}
