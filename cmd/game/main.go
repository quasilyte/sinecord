package main

import (
	"time"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/eui"
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
	}

	state.UIResources = eui.PrepareResources(ctx.Loader)

	session.ReloadLanguage(ctx, "en")

	if err := ge.RunGame(ctx, scenes.NewMainMenuController(state)); err != nil {
		panic(err)
	}
}
