package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ebitengine-gamejam2023/eui"
	"github.com/quasilyte/ge"
)

func initUI(scene *ge.Scene, root *widget.Container) {
	uiObject := eui.NewSceneObject(root)
	scene.AddGraphics(uiObject)
	scene.AddObject(uiObject)
}
