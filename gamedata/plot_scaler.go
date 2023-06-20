package gamedata

import (
	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

type PlotScaler struct {
	Factor float64
	Offset gmath.Vec
}

func (s *PlotScaler) TranslateTiledPos(tileset *tiled.Tileset, pos gmath.Vec) gmath.Vec {
	return gmath.Vec{
		X: pos.X + s.Offset.X + (tileset.TileWidth / 2),
		Y: pos.Y - (tileset.TileHeight / 2),
	}
}

func (s *PlotScaler) ScaleXY(x, y float64) gmath.Vec {
	return s.ScalePos(gmath.Vec{X: x, Y: y})
}

func (s *PlotScaler) ScalePos(pos gmath.Vec) gmath.Vec {
	pos = gmath.Vec{
		X: pos.X * s.Factor,
		Y: -(pos.Y * s.Factor),
	}
	return pos.Add(s.Offset)
}
