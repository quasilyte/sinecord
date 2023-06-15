package assets

import (
	"embed"
	"io"

	"github.com/quasilyte/ge"
)

//go:embed all:_data
var gameAssets embed.FS

func MakeOpenAssetFunc(ctx *ge.Context) func(path string) io.ReadCloser {
	return func(path string) io.ReadCloser {
		f, err := gameAssets.Open("_data/" + path)
		if err != nil {
			ctx.OnCriticalError(err)
		}
		return f
	}
}
func RegisterResources(ctx *ge.Context) {
	registerFontResources(ctx)
	registerRawResources(ctx)
	registerImageResources(ctx)
}
