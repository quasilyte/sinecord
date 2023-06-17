//go:build ignore
// +build ignore

package main

func sqr(x float) float { return x * x }

func Fragment(_ vec4, texCoord vec2, _ vec4) vec4 {
	const scale float = 46.0

	pixPos := tex2pixCoord(texCoord)
	realX := floor(pixPos.x)
	realY := floor(pixPos.y)

	const offsetX = 4.0

	x := (realX - offsetX) / scale
	y := -(realY - (scale * 3.0)) / scale

	epsilon := 0.05

	delta1 := abs(y - log2(x))
	if realX > offsetX && delta1 < epsilon {
		return vec4(0.616, 0.843, 0.576, 1.0)
	}
	// delta2 := abs(y - sqrt(x)/2.0)
	// if realX > offsetX && delta2 < epsilon {
	// 	return vec4(0.682, 0.576, 0.843, 1.0)
	// }
	// delta3 := abs(y - sin(x))
	// if realX > offsetX && delta3 < epsilon {
	// 	return vec4(0.847, 0.784, 0.506, 1.0)
	// }
	// delta4 := abs(y - 2*cos(x/2))
	// if realX > offsetX && delta4 < epsilon {
	// 	return vec4(0.925, 0.09, 0.514, 1.0)
	// }

	return imageSrc0At(texCoord)
}

/*
func Fragment(_ vec4, texCoord vec2, _ vec4) vec4 {
	const scale float = 46.0

	pixPos := tex2pixCoord(texCoord)
	realX := floor(pixPos.x)
	realY := floor(pixPos.y)

	const offsetX = 4.0

	// x := (realX - (64.0))
	// y := -(realY - (scale * 3) - 5.0)

	x := realX - offsetX
	y := -(realY - (scale * 3.0))

	if realX > offsetX && abs(y-scale*sin(x/scale)) < 1.0 {
		return vec4(1.0, 0.0, 0.0, 1.0)
	}

	// if abs(y-sqrt(x*scale)) < 1.0 {
	// 	return vec4(1.0, 0.0, 0.0, 1.0)
	// }

	return imageSrc0At(texCoord)
}
*/

func tex2pixCoord(texCoord vec2) vec2 {
	pixSize := imageSrcTextureSize()
	originTexCoord, _ := imageSrcRegionOnTexture()
	actualTexCoord := texCoord - originTexCoord
	actualPixCoord := actualTexCoord * pixSize
	return actualPixCoord
}
