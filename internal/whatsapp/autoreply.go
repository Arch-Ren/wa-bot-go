package whatsapp

import (
	"context"
	"log"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type AutoReplyHandler struct {
	Client   *whatsmeow.Client
	Settings *SettingsStore
}

func NewAutoReplyHandler(client *whatsmeow.Client, settings *SettingsStore) *AutoReplyHandler {
	return &AutoReplyHandler{
		Client:   client,
		Settings: settings,
	}
}

func (h *AutoReplyHandler) Register() {
	h.Client.AddEventHandler(h.handleEvent)
	log.Println("Auto-reply handler registered")
}

func (h *AutoReplyHandler) handleEvent(evt any) {
	msg, ok := evt.(*events.Message)
	if !ok {
		return
	}

	// Skip pesan dari diri sendiri
	if msg.Info.IsFromMe {
		return
	}

	// Skip group chat
	if msg.Info.IsGroup {
		return
	}

	// Extract text dari pesan
	text := extractText(msg)
	if text == "" {
		return
	}

	chatJID := msg.Info.Chat.String()

	// Handle command
	if strings.HasPrefix(text, "/autoreply ") {
		h.handleCommand(msg, strings.TrimPrefix(text, "/autoreply "))
		return
	}

	// Auto-reply logic
	if !h.Settings.IsAutoReplyEnabled() {
		return
	}

	//cooldown state
	if h.Settings.HasRepliedToday(chatJID) {
		return
	}

	// Kirim auto-reply
	replyText := h.Settings.GetAutoReplyText()
	h.sendReply(msg, replyText)
	h.Settings.MarkReplied(chatJID)
}

func (h *AutoReplyHandler) handleCommand(msg *events.Message, command string) {
	chatJID := msg.Info.Chat
	command = strings.TrimSpace(strings.ToLower(command))

	var reply string

	switch command {
	case "on":
		h.Settings.SetAutoReplyEnabled(true)
		reply = "✅ Auto-reply diaktifkan"
	case "off":
		h.Settings.SetAutoReplyEnabled(false)
		reply = "❌ Auto-reply dimatikan"
	case "status":
		if h.Settings.IsAutoReplyEnabled() {
			reply = "📊 Status: Auto-reply AKTIF"
		} else {
			reply = "📊 Status: Auto-reply MATI"
		}
	case "help":
		reply = "📖 Perintah auto-reply:\n" +
			"/autoreply on - Aktifkan auto-reply\n" +
			"/autoreply off - Matikan auto-reply\n" +
			"/autoreply status - Lihat status\n" +
			"/autoreply help - Bantuan ini"
	default:
		reply = "❓ Perintah tidak dikenal. Ketik /autoreply help"
	}

	h.sendReplyDirect(chatJID, reply)
}

func (h *AutoReplyHandler) sendReply(msg *events.Message, text string) {
	h.sendReplyDirect(msg.Info.Chat, text)
}

func (h *AutoReplyHandler) sendReplyDirect(chatJID types.JID, text string) {
	_, err := h.Client.SendMessage(context.Background(), chatJID, &waE2E.Message{
		Conversation: &text,
	})
	if err != nil {
		log.Printf("Gagal kirim auto-reply ke %s: %v", chatJID.String(), err)
	}
}

func extractText(msg *events.Message) string {
	if msg.Message == nil {
		return ""
	}

	// Pesan teks biasa
	if text := msg.Message.GetConversation(); text != "" {
		return text
	}

	// Pesan extended text
	if extText := msg.Message.GetExtendedTextMessage(); extText != nil {
		if text := extText.GetText(); text != "" {
			return text
		}
	}

	return ""
}
