package command

import (
	"context"
	"log"
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
		senderJID := ctx.Sender.ToNonAD().String()
		log.Printf("Checking admin for: %s", senderJID)
		isAdmin, err := r.AdminChecker.IsAdmin(ctx.Context, senderJID)

		if err != nil {
			log.Printf("Admin check error: %v", err)
			return true
		}

		if !isAdmin {
			log.Printf("Access denied for: %s", senderJID)
			return false
		}
	}

	cmd.Handler(ctx)

	return true
}
