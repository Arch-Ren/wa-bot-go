package command

import (
	"context"
)

type Router struct {
	Registry     *Registry
	AdminChecker AdminChecker
}

type AdminChecker interface {
	IsAdmin(ctx context.Context, jid string) (bool, error)
}

func NewRouter(registry *Registry, adminChecker AdminChecker) *Router {
	return &Router{
		Registry:     registry,
		AdminChecker: adminChecker,
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

	if cmd.AdminOnly {
		isAdmin, err := r.AdminChecker.IsAdmin(ctx.Context, ctx.Sender.String())

		if err != nil {
			return true
		}

		if !isAdmin {
			return false
		}
	}

	cmd.Handler(ctx)

	return true
}
