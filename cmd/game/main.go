package main

import (
	"time"

	"github.com/quasilyte/ebitengine-gamejam2023/assets"
	"github.com/quasilyte/ebitengine-gamejam2023/eui"
	"github.com/quasilyte/ebitengine-gamejam2023/scenes"
	"github.com/quasilyte/ebitengine-gamejam2023/session"
	"github.com/quasilyte/ge"
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

	state := &session.State{}
	state.UIResources = eui.PrepareResources(ctx.Loader)

	session.ReloadLanguage(ctx, "en")

	if err := ge.RunGame(ctx, scenes.NewMainMenuController(state)); err != nil {
		panic(err)
	}
}
