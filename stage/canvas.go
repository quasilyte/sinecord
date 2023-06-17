package stage

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
)

type Canvas struct {
	sprites []*ge.Sprite
}

func NewCanvas() *Canvas {
	return &Canvas{
		sprites: make([]*ge.Sprite, 0, 16),
	}
}

func (c *Canvas) AddSprite(s *ge.Sprite) {
	c.sprites = append(c.sprites, s)
}

func (c *Canvas) IsDisposed() bool { return false }

func (c *Canvas) Draw(screen *ebiten.Image) {
	liveSprites := c.sprites[:0]
	for _, s := range c.sprites {
		if s.IsDisposed() {
			continue
		}
		s.Draw(screen)
		liveSprites = append(liveSprites, s)
	}
	c.sprites = liveSprites
}
