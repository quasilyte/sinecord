package stage

type Config struct {
	MaxInstruments int
}

type Context struct {
	runner *programRunner

	config Config
}

func NewContext(config Config) *Context {
	return &Context{
		config: config,
		runner: &programRunner{},
	}
}
