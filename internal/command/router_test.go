package command

import "testing"

func TestRouter(t *testing.T) {
	registry := NewRegistry()

	registry.Register(Command{
		Name:        "ping",
		Prefix:      ".",
		Description: "Test Command",
		AdminOnly:   false,
		Handler:     PingHandler,
	})

	router := NewRouter(registry)

	ctx := &Context{}

	handled := router.Route(ctx, ".ping")

	if !handled {
		t.Fatal("expected command to be handled")
	}
}
