package command

import (
	"context"
	"testing"
)

type mockAdminChecker struct {
	allowed bool
}

func (m *mockAdminChecker) IsAdmin(ctx context.Context, jid string) (bool, error) {
	return m.allowed, nil
}

func TestRouter(t *testing.T) {
	registry := NewRegistry()

	registry.Register(Command{
		Name:        "ping",
		Prefix:      ".",
		Description: "Test Command",
		AdminOnly:   false,
		Handler:     PingHandler,
	})

	adminChecker := &mockAdminChecker{
		allowed: false,
	}

	router := NewRouter(registry, adminChecker)

	ctx := &Context{}

	handled := router.Route(ctx, ".ping")

	if !handled {
		t.Fatal("expected command to be handled")
	}
}

func TestRouterPublicCommand(t *testing.T) {
	registry := NewRegistry()

	called := false

	registry.Register(Command{
		Name:        "ping",
		Prefix:      ".",
		Description: "Test Command",
		AdminOnly:   false,
		Handler: func(ctx *Context) {
			called = true
		},
	})

	AdminChecker := &mockAdminChecker{
		allowed: false,
	}

	router := NewRouter(registry, AdminChecker)

	ctx := &Context{
		Context: context.Background(),
	}

	handled := router.Route(ctx, ".ping")

	if !handled {
		t.Fatal("Expected command to be handled")
	}

	if !called {
		t.Fatal("Expected handler to be called")
	}
}
