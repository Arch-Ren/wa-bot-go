package main

import (
	"fmt"

	"wa-bot-go/internal/whatsapp"
)

func main() {
	client := whatsapp.Client{}

	fmt.Println("WhatsApp Bot Starting...")
	fmt.Printf("Client: %+v\n", client)
}
