package stage

import (
	"github.com/quasilyte/sinecord/gamedata"
)

type Config struct {
	MaxInstruments int

	Targets []gamedata.Target

	Track gamedata.Track
}

type Context struct {
	runner *programRunner

	config Config

	Scaler *gamedata.PlotScaler
}

func NewContext(config Config) *Context {
	return &Context{
		config: config,
		runner: &programRunner{},
	}
}
