package command

type Command struct {
	Name        string
	Prefix      string
	Description string
	AdminOnly   bool
	Handler     Handler
}

type Handler func(ctx *Context)
