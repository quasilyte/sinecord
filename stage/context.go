package stage

import "github.com/quasilyte/gmath"

type Config struct {
	MaxInstruments int

	Targets []Target
}

type Context struct {
	runner *programRunner

	config Config

	PlotScale  float64
	PlotOffset gmath.Vec
}

func NewContext(config Config) *Context {
	return &Context{
		config: config,
		runner: &programRunner{},
	}
}
