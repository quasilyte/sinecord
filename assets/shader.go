package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerShaderResources(ctx *ge.Context) {
	shaderResources := map[resource.ShaderID]resource.ShaderInfo{
		ShaderSine: {Path: "shader/sin.go"},
	}

	for id, res := range shaderResources {
		ctx.Loader.ShaderRegistry.Set(id, res)
		ctx.Loader.LoadShader(id)
	}
}

const (
	ShaderNone resource.ShaderID = iota

	ShaderSine
)
