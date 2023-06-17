package stage

type Context struct {
	runner *programRunner
}

func NewContext() *Context {
	return &Context{
		runner: &programRunner{},
	}
}
