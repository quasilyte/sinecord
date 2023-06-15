package eui

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
)

type SceneObject struct {
	ui *ebitenui.UI
}

func NewSceneObject(root *widget.Container) *SceneObject {
	return &SceneObject{
		ui: &ebitenui.UI{
			Container: root,
		},
	}
}

func (o *SceneObject) IsDisposed() bool { return false }

func (o *SceneObject) Init(scene *ge.Scene) {}

func (o *SceneObject) Update(delta float64) {
	o.ui.Update()
}

func (o *SceneObject) Draw(screen *ebiten.Image) {
	o.ui.Draw(screen)
}
