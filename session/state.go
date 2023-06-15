package session

import (
	"fmt"

	"github.com/quasilyte/ebitengine-gamejam2023/assets"
	"github.com/quasilyte/ebitengine-gamejam2023/eui"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/langs"
)

type State struct {
	UIResources *eui.Resources
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
