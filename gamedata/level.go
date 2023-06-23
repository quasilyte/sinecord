package gamedata

import (
	"errors"
	"fmt"
	"strings"

	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

type LevelBonusObjectives struct {
	MaxInstruments int

	ForbiddenFuncs []string

	AllTargets   bool
	AvoidPenalty bool
}

type LevelData struct {
	Name string

	Targets []Target

	ActNumber     int
	MissionNumber int

	MaxInstruments int
	Description    string

	Solution Track

	Bonus LevelBonusObjectives
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

	var systemLayer *tiled.MapLayer
	var targetsLayer *tiled.MapLayer
	for i := range m.Layers {
		l := &m.Layers[i]
		switch l.Name {
		case "targets":
			targetsLayer = l
		case "system":
			systemLayer = l
		}
	}

	if targetsLayer == nil {
		return nil, errors.New("can't find targets layer")
	}
	if systemLayer == nil {
		return nil, errors.New("can't find system layer")
	}

	instrumentMap := map[string]InstrumentKind{}
	for kind := BassInstrument; kind <= AnyInstrument; kind++ {
		key := strings.ToLower(InstrumentShape(kind).String()) + "_target"
		instrumentMap[key] = kind
	}

	for _, o := range systemLayer.Objects {
		id := o.GID - int64(ref.FirstGID)
		t := tileset.TileByID(int(id))
		switch t.Class {
		case "level_settings":
			result.Name = o.GetStringProp("name", "")
			result.MaxInstruments = o.GetIntProp("max_instruments", 2)
			result.Description = o.GetStringProp("description", "")

			result.Bonus.AllTargets = o.GetBoolProp("bonus_all_targets", false)
			result.Bonus.AvoidPenalty = o.GetBoolProp("bonus_avoid_penalty", false)
			result.Bonus.MaxInstruments = o.GetIntProp("bonus_max_instruments", 1)
			for _, fn := range strings.Split(o.GetStringProp("bonus_forbidden_funcs", ""), ",") {
				fn = strings.TrimSpace(fn)
				if fn != "" {
					result.Bonus.ForbiddenFuncs = append(result.Bonus.ForbiddenFuncs, fn)
				}
			}
		}
	}

	if result.Name == "" {
		return nil, errors.New("a level *name* can't be empty")
	}
	if result.Description == "" {
		return nil, errors.New("a level *description* can't be empty")
	}
	if result.MaxInstruments == 0 {
		return nil, errors.New("a *max_instruments* can't be zero")
	}

	for _, o := range targetsLayer.Objects {
		id := o.GID - int64(ref.FirstGID)
		t := tileset.TileByID(int(id))
		outline := false
		class := t.Class
		if strings.HasPrefix(class, "outline_") {
			outline = true
			class = class[len("outline_"):]
		}
		if _, ok := instrumentMap[class]; !ok {
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
			Instrument: instrumentMap[class],
			Size:       size,
			Outline:    outline,
		})
	}

	return &result, nil
}
