package whatsapp

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type Client struct {
	Container *sqlstore.Container
	WhatsApp  *whatsmeow.Client
}

func NewClient(container *sqlstore.Container) (*Client, error) {
	device, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(device, nil)

	return &Client{
		Container: container,
		WhatsApp:  client,
	}, nil
}
