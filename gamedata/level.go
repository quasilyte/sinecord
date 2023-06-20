package gamedata

import (
	"errors"
	"fmt"
	"strings"

	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

type LevelData struct {
	Targets []Target
}

func ParseTileset(jsonData []byte) (*tiled.Tileset, error) {
	return tiled.UnmarshalTileset(jsonData)
}

func ParseLevel(tileset *tiled.Tileset, scaler *PlotScaler, jsonData []byte) (*LevelData, error) {
	var result LevelData

	m, err := tiled.UnmarshalMap(jsonData)
	if err != nil {
		return nil, err
	}
	if len(m.Tilesets) != 1 {
		return nil, errors.New("expected exactly 1 tileset")
	}

	ref := m.Tilesets[0]

	var targetsLayer *tiled.MapLayer
	for i := range m.Layers {
		l := &m.Layers[i]
		if l.Name == "targets" {
			targetsLayer = l
			break
		}
	}

	if targetsLayer == nil {
		return nil, errors.New("can't find targets layer")
	}

	instrumentMap := map[string]InstrumentKind{}
	for kind := BassInstrument; kind <= OtherInstrument; kind++ {
		key := strings.ToLower(InstrumentShape(kind).String()) + "_target"
		instrumentMap[key] = kind
	}

	for _, o := range targetsLayer.Objects {
		id := o.GID - int64(ref.FirstGID)
		t := tileset.TileByID(int(id))
		if _, ok := instrumentMap[t.Class]; !ok {
			return nil, fmt.Errorf("unexpected kind target: %q", t.Class)
		}

		x := float64(o.X)
		y := float64(o.Y)
		var size TargetSize
		switch o.Width {
		case 23:
			size = SmallTarget
			x -= 13
			y += 13
		case 46:
			size = NormalTarget
		case 69:
			size = BigTarget
			x += 13
			y -= 13
		default:
			return nil, fmt.Errorf("unexpected target size: %dx%d", o.Width, o.Height)
		}

		result.Targets = append(result.Targets, Target{
			Pos:        scaler.TranslateTiledPos(tileset, gmath.Vec{X: x, Y: y}),
			Instrument: instrumentMap[t.Class],
			Size:       size,
		})
	}

	return &result, nil
}
