package assets

import (
	"embed"
	"fmt"
	"io"
	"strconv"
	"strings"

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
	registerShaderResources(ctx)
	registerAudioResources(ctx)

	decodeSoundfonts(ctx)
}

type RawLevelData struct {
	Act     int
	Mission int

	TileMap  []byte
	Solution []byte
}

func ReadLevelsData() []RawLevelData {
	var result []RawLevelData

	acts, err := gameAssets.ReadDir("_data/raw/level")
	if err != nil {
		panic(err)
	}
	for _, a := range acts {
		if !strings.HasPrefix(a.Name(), "act") {
			continue
		}
		actNumber, err := strconv.Atoi(a.Name()[len("act"):])
		if err != nil {
			panic(err)
		}
		missions, err := gameAssets.ReadDir("_data/raw/level/" + a.Name())
		if err != nil {
			panic(err)
		}
		for _, m := range missions {
			if !strings.HasPrefix(m.Name(), "mission") {
				continue
			}
			missionNumber, err := strconv.Atoi(m.Name()[len("mission"):])
			if err != nil {
				panic(err)
			}
			tileMap, err := gameAssets.ReadFile(fmt.Sprintf("_data/raw/level/%s/%s/map.tmj", a.Name(), m.Name()))
			if err != nil {
				panic(err)
			}
			solution, err := gameAssets.ReadFile(fmt.Sprintf("_data/raw/level/%s/%s/solution.json", a.Name(), m.Name()))
			if err != nil {
				panic(err)
			}
			result = append(result, RawLevelData{
				Act:      actNumber,
				Mission:  missionNumber,
				TileMap:  tileMap,
				Solution: solution,
			})
		}
	}

	return result
}
