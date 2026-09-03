package whatsapp

import (
	"context"

	"go.mau.fi/whatsmeow/store/sqlstore"
	//waLog "go.mau.fi/whatsmeow/util/log"
	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase() (*sqlstore.Container, error) {

	container, err := sqlstore.New(
		context.Background(),
		"sqlite3",
		"file:data/whatsapp.db?_foreign_keys=on",
		nil,
	)

	if err != nil {
		return nil, err
	}

	return container, nil
}
