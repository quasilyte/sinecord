package session

import (
	"fmt"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/langs"
	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/gamedata"
)

type State struct {
	UIResources *eui.Resources

	LevelTileset *tiled.Tileset

	LevelsByAct [][]*gamedata.LevelData

	PlotScaler *gamedata.PlotScaler

	Input *input.Handler

	Persistent PersistentData

	EffectiveVolume float64
}

type PersistentData struct {
	LevelsCompleted []LevelCompletionInfo `json:"levels_completed"`

	VolumeLevel int `json:"volume_level"`
}

type LevelCompletionInfo struct {
	Name  string `json:"name"`
	Bonus bool   `json:"bonus"`
}

type LevelCompletionStatus int

const (
	LevelNotCompleted LevelCompletionStatus = iota
	LevelCompleted
	LevelCompletedWithBonus
)

func (d *PersistentData) UpdateLevelCompletion(level *gamedata.LevelData, bonus bool) {
	var existingEntry *LevelCompletionInfo
	for i := range d.LevelsCompleted {
		l := &d.LevelsCompleted[i]
		if l.Name == level.Name {
			existingEntry = l
			break
		}
	}
	if existingEntry == nil {
		d.LevelsCompleted = append(d.LevelsCompleted, LevelCompletionInfo{
			Name:  level.Name,
			Bonus: bonus,
		})
		return
	}
	if existingEntry.Bonus {
		return
	}
	existingEntry.Bonus = bonus
}

func (d *PersistentData) GetLevelCompletionStatus(level *gamedata.LevelData) LevelCompletionStatus {
	for _, entry := range d.LevelsCompleted {
		if entry.Name == level.Name {
			if entry.Bonus {
				return LevelCompletedWithBonus
			}
			return LevelCompleted
		}
	}
	return LevelNotCompleted
}

func ReloadLanguage(ctx *ge.Context, language string) {
	var id resource.RawID
	switch language {
	case "en":
		id = assets.RawEnDictBase
	default:
		panic(fmt.Sprintf("unsupported lang: %q", language))
	}
	dict, err := langs.ParseDictionary(language, 4, ctx.Loader.LoadRaw(id).Data)
	if err != nil {
		panic(err)
	}
	ctx.Dict = dict
}
