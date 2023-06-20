package main

import (
	"time"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/sinecord/assets"
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

	keymap := input.Keymap{}
	state := &session.State{
		Input: ctx.Input.NewHandler(0, keymap),
		PlotScaler: &gamedata.PlotScaler{
			Factor: 46,
			Offset: gmath.Vec{
				X: 4,
				Y: 46 * 3,
			},
		},
		UIResources: eui.PrepareResources(ctx.Loader),
	}

	session.ReloadLanguage(ctx, "en")

	tileset, err := gamedata.ParseTileset(ctx.Loader.LoadRaw(assets.RawLevelTilesetJSON).Data)
	if err != nil {
		panic(err)
	}
	state.LevelTileset = tileset

	if err := ge.RunGame(ctx, scenes.NewMainMenuController(state)); err != nil {
		panic(err)
	}
}
