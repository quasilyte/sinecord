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

	PlotScaler *gamedata.PlotScaler

	Input *input.Handler
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
