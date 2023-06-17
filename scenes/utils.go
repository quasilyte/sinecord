package scenes

import (
	"bytes"
	"encoding/binary"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/sinecord/eui"
)

func initUI(scene *ge.Scene, root *widget.Container) {
	uiObject := eui.NewSceneObject(root)
	scene.AddGraphics(uiObject)
	scene.AddObject(uiObject)
}

func generatePCM(left, right []float32) []byte {
	length := len(left)

	a := float32(32768.0)

	data := make([]int16, 2*length)

	for i := 0; i < length; i++ {
		data[2*i] = int16(a * left[i])
		data[2*i+1] = int16(a * right[i])
	}

	var buf bytes.Buffer
	buf.Grow(len(data) * 2)

	binary.Write(&buf, binary.LittleEndian, data)
	return buf.Bytes()
}
