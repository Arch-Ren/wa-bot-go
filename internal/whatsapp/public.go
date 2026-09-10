package whatsapp

import (
	"context"
	"log"
	"strings"

	"wa-bot-go/internal/command"

	"go.mau.fi/whatsmeow/proto/waE2E"
)

type PublicComands struct {
	Registry *command.Registry
}

func NewPublicCommand(registry *command.Registry) *PublicComands {
	return &PublicComands{
		Registry: registry,
	}
}

func (p *PublicComands) RegiserCommand(registry *command.Registry) {
	registry.Register(command.Command{
		Name:        "about",
		Prefix:      ".",
		Description: "Tentang bot ini",
		AdminOnly:   false,
		Handler:     p.handleAbout,
	})

	registry.Register(command.Command{
		Name:        "help",
		Prefix:      ".",
		Description: "Menampilkan daftar command",
		AdminOnly:   false,
		Handler:     p.handleHelp,
	})
}

func (p *PublicComands) handleAbout(ctx *command.Context) {
	text := "Ini wa-bot"

	_, err := ctx.Client.SendMessage(
		context.Background(),
		ctx.Chat,
		&waE2E.Message{
			Conversation: &text,
		},
	)

	if err != nil {
		log.Printf("gagal kirim .about: %v", err)
	}
}

func (p *PublicComands) handleHelp(ctx *command.Context) {
	command := p.Registry.PublicComands()

	var builder strings.Builder

	builder.WriteString("Daftar command yang tersedia: \n\n")

	for _, cmd := range command {
		builder.WriteString(cmd.Prefix)
		builder.WriteString(cmd.Name)
		builder.WriteString(" - ")
		builder.WriteString(cmd.Description)
		builder.WriteString("\n")
	}

	text := builder.String()

	_, err := ctx.Client.SendMessage(
		context.Background(),
		ctx.Chat,
		&waE2E.Message{
			Conversation: &text,
		},
	)

	if err != nil {
		log.Printf("gagal kirim .help: %v", err)
	}
}
