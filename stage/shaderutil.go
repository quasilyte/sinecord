package stage

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

var plotColors = []string{
	`vec4(0.349, 0.529, 0.940, 1.0)`,
	`vec4(0.682, 0.576, 0.843, 1.0)`,
	`vec4(0.847, 0.784, 0.506, 1.0)`,
	`vec4(0.925, 0.090, 0.514, 1.0)`,
}

func CompilePlotShader(id int, fx string) (*ebiten.Shader, error) {
	shaderText := []byte(`package main
	func tex2pixCoord(texCoord vec2) vec2 {
		pixSize := imageSrcTextureSize()
		originTexCoord, _ := imageSrcRegionOnTexture()
		actualTexCoord := texCoord - originTexCoord
		actualPixCoord := actualTexCoord * pixSize
		return actualPixCoord
	}
	func Fragment(_ vec4, _texCoord vec2, _ vec4) vec4 {
		const _scale float = 46.0
	
		_pixPos := tex2pixCoord(_texCoord)
		_realX := floor(_pixPos.x)
		_realY := floor(_pixPos.y)
	
		const _offsetX = 4.0
	
		x := (_realX - _offsetX) / _scale
		_y := -(_realY - (_scale * 3.0)) / _scale

		_ = x

		_clr := imageSrc0UnsafeAt(_texCoord)
		_delta := abs(_y - (` + fx + `))
		if _realX > _offsetX && _delta < 0.08 {
			return ` + plotColors[id] + `+(_clr*0.25)
		}
	
		return _clr
	}`)

	return compileShaderSafe(shaderText)
}

func compileShaderSafe(text []byte) (shader *ebiten.Shader, err error) {
	defer func() {
		rv := recover()
		if rv != nil {
			shader = nil
			err = errors.New("shader compilation panic")
		}
	}()
	return ebiten.NewShader(text)
}
