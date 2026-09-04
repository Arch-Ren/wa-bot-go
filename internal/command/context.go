package command

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type Context struct {
	Client  *whatsmeow.Client
	Message *events.Message
	Chat    types.JID
	Sender  types.JID
	Args    []string
}
