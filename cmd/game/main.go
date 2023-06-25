package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	ebiteninput "github.com/ebitenui/ebitenui/input"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/controls"
	"github.com/quasilyte/sinecord/eui"
	"github.com/quasilyte/sinecord/gamedata"
	"github.com/quasilyte/sinecord/scenes"
	"github.com/quasilyte/sinecord/session"
)

func main() {
	ctx := ge.NewContext(ge.ContextConfig{})
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "sinecord"
	ctx.WindowTitle = "Sinecord"
	ctx.WindowWidth = 1920
	ctx.WindowHeight = 1080
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = assets.MakeOpenAssetFunc(ctx)
	assets.RegisterResources(ctx)

	keymap := input.Keymap{
		controls.ActionBack:       {input.KeyEscape},
		controls.ActionTogglePlay: {input.KeySpace},
	}
	state := &session.State{
		Persistent: getDefaultData(),
		Input:      ctx.Input.NewHandler(0, keymap),
		PlotScaler: &gamedata.PlotScaler{
			Factor: 46,
			Offset: gmath.Vec{
				X: 4,
				Y: 46 * 3,
			},
		},
		UIResources: eui.PrepareResources(ctx.Loader),
		Track: gamedata.Track{
			Name: "session",
			Instruments: []gamedata.InstrumentSettings{
				{Function: "sin(x/2) - 0.2", PeriodFunction: "pi/9", Volume: 1.0, InstrumentName: "Synth Bass 1", Enabled: true},
				{Function: "pi/4", PeriodFunction: "pi/3", Volume: 1.0, InstrumentName: "Synth Bass 1", Enabled: true},
			},
		},
	}

	if err := ctx.LoadGameData("save", &state.Persistent); err != nil {
		fmt.Printf("can't load game data: %v", err)
		state.Persistent = getDefaultData()
		ctx.SaveGameData("save", state.Persistent)
	}

	session.ReloadLanguage(ctx, "en")

	tileset, err := gamedata.ParseTileset(ctx.Loader.LoadRaw(assets.RawLevelTilesetJSON).Data)
	if err != nil {
		panic(err)
	}
	state.LevelTileset = tileset

	{
		allLevels := assets.ReadLevelsData()
		maxAct := 1
		for _, raw := range allLevels {
			if raw.Act > maxAct {
				maxAct = raw.Act
			}
		}
		levelNames := map[string]struct{}{}
		levelsByAct := make([][]*gamedata.LevelData, maxAct+1)
		for _, raw := range allLevels {
			parsed, err := gamedata.ParseLevel(state.LevelTileset, state.PlotScaler, raw.TileMap)
			if err != nil {
				panic(err)
			}
			if _, ok := levelNames[parsed.Name]; ok {
				panic(fmt.Sprintf("duplicated level name: act %d, %q", raw.Act, parsed.Name))
			}
			levelNames[parsed.Name] = struct{}{}
			parsed.ActNumber = raw.Act
			parsed.MissionNumber = raw.Mission
			if err := json.Unmarshal(raw.Solution, &parsed.Solution); err != nil {
				panic(err)
			}
			levelsByAct[raw.Act] = append(levelsByAct[raw.Act], parsed)
		}
		for _, levels := range levelsByAct {
			sort.SliceStable(levels, func(i, j int) bool {
				return levels[i].MissionNumber < levels[j].MissionNumber
			})
		}
		state.LevelsByAct = levelsByAct
	}

	ebiteninput.SetCursorImage(ebiteninput.CURSOR_DEFAULT, ctx.Loader.LoadImage(assets.ImagePointerNormal).Data)
	ebiteninput.SetCursorImage("hand", ctx.Loader.LoadImage(assets.ImagePointerHand).Data)

	if err := ge.RunGame(ctx, scenes.NewMainMenuController(state)); err != nil {
		panic(err)
	}
}

func getDefaultData() session.PersistentData {
	data := session.PersistentData{
		VolumeLevel: 5,
	}
	return data
}
