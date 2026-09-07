package command

type Router struct {
	Registry *Registry
}

func NewRouter(registry *Registry) *Router {
	return &Router{
		Registry: registry,
	}
}

func (r *Router) Route(ctx *Context, text string) bool {
	parsed := Parse(text)

	if parsed == nil {
		return false
	}

	cmd := r.Registry.Find(parsed.Prefix, parsed.Name)

	if cmd == nil {
		return false
	}

	ctx.Args = parsed.Args

	cmd.Handler(ctx)

	return true
}
