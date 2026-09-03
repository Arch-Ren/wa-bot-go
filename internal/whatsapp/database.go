package whatsapp

import (
	"context"

	"go.mau.fi/whatsmeow/store/sqlstore"
	//waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite"
)

func NewDatabase() (*sqlstore.Container, error) {

	container, err := sqlstore.New(
		context.Background(),
		"sqlite",
		"file:data/whatsapp.db?_foreign_keys=on",
		nil,
	)

	if err != nil {
		return nil, err
	}

	return container, nil
}
