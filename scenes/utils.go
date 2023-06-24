package scenes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/sinecord/assets"
	"github.com/quasilyte/sinecord/eui"
)

func initUI(scene *ge.Scene, root *widget.Container) {
	bg := scene.NewSprite(assets.ImageMenuBackground)
	bg.Centered = false
	scene.AddGraphics(bg)

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

func formatDateISO8601(d time.Time, withTime bool) string {
	s := fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
	if withTime {
		s += fmt.Sprintf(" %02d:%02d", d.Hour(), d.Minute())
	}
	return s
}
