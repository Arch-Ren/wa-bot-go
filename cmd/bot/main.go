package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdp/qrterminal/v3"

	"wa-bot-go/internal/whatsapp"
)

func main() {
	database, err := whatsapp.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	client, err := whatsapp.NewClient(database)
	if err != nil {
		log.Fatal(err)
	}

	settings, err := whatsapp.NewSettingsStore("file:data/whatsapp.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	defer settings.Close()

	//Register auto-reply handler
	autoreply := whatsapp.NewAutoReplyHandler(client.WhatsApp, settings)
	autoreply.Register()

	//Belum pernah login
	if client.WhatsApp.Store.ID == nil {
		qrChan, err := client.WhatsApp.GetQRChannel(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		if err := client.WhatsApp.Connect(); err != nil {
			log.Fatal(err)
		}

		for evt := range qrChan {
			switch evt.Event {
			case "code":
				fmt.Println("Scan this QR Code:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)

			case "succes":
				fmt.Println("WhatsApp Connected!")

			default:
				fmt.Printf("Login event: %s\n", evt.Event)
			}
		}
	} else {
		//sudah pernah login
		if err := client.WhatsApp.Connect(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("WhatsApp Connected")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	client.WhatsApp.Disconnect()
	fmt.Println("WhatsApp Disconnected!")
}
