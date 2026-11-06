package whatsapp

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"wa-bot-go/internal/command"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

const (
	stickerSize  = 512
	cornerRadius = 40
)

type StickerHandler struct {
	Client *whatsmeow.Client
}

func NewStickerHandler(client *whatsmeow.Client) *StickerHandler {
	return &StickerHandler{
		Client: client,
	}
}

func (s *StickerHandler) RegisterCommand(registry *command.Registry) {
	registry.Register(command.Command{
		Name:        "sticker",
		Prefix:      ".",
		Description: "Ubah foto/gif/video jadi stiker",
		AdminOnly:   false,
		Handler:     s.handleSticker,
	})

	registry.Register(command.Command{
		Name:        "s",
		Prefix:      ".",
		Description: "Shortcut untuk .sticker",
		AdminOnly:   false,
		Handler:     s.handleSticker,
	})
}

func (s *StickerHandler) handleSticker(ctx *command.Context) {
	msg := ctx.Message
	log.Printf("Menerima permintaan stiker dari %s", ctx.Sender.String())

	// Cek apakah caption mengandung .sticker atau .s
	text := extractText(msg)
	isStickerCmd := strings.Contains(text, ".sticker") || strings.Contains(text, ".s")

	mediaMsg := s.extractMedia(msg)

	// Jika bukan reply, kita harus pastikan pesan ini mengandung media DAN perintah
	if mediaMsg == nil {
		log.Println("Media tidak ditemukan")
		s.sendText(ctx, "Kirim atau reply foto/gif/video dengan caption *.sticker* atau *.s*")
		return
	}

	// Jika media ditemukan, tapi tidak ada perintah di caption,
	// dan bukan reply, jangan proses.
	if !isStickerCmd && !s.isReply(msg) {
		return
	}

	log.Printf("Media ditemukan, mulai download...")
	data, err := s.downloadMedia(mediaMsg)
	if err != nil {
		log.Printf("Gagal download media: %v", err)
		s.sendText(ctx, "Gagal download media")
		return
	}
	log.Printf("Download selesai (%d bytes)", len(data))

	var webpData []byte
	var isAnimated bool

	doCrop := false
	if len(ctx.Args) > 0 && (ctx.Args[0] == "sq" || ctx.Args[0] == "square") {
		doCrop = true
		log.Println("Crop square diaktifkan")
	}

	switch m := mediaMsg.(type) {
	case *waE2E.ImageMessage:
		log.Printf("Memproses Image (mimetype: %s)", m.GetMimetype())
		webpData, err = s.convertToWebP(data, m.GetMimetype(), false, doCrop)
		isAnimated = false
	case *waE2E.VideoMessage:
		log.Printf("Memproses Video/GIF (mimetype: %s)", m.GetMimetype())
		webpData, err = s.convertToWebP(data, m.GetMimetype(), true, doCrop)
		isAnimated = true
	}
// ... rest of the code

	if err != nil {
		log.Printf("gagal konversi stiker: %v", err)
		s.sendText(ctx, "Gagal membuat stiker: "+err.Error())
		return
	}

	err = s.sendSticker(ctx, webpData, isAnimated)
	if err != nil {
		log.Printf("gagal kirim stiker: %v", err)
		s.sendText(ctx, "Gagal mengirim stiker")
	}
}

// extractMedia mengembalikan media message dari pesan langsung atau quoted message
func (s *StickerHandler) extractMedia(msg *events.Message) interface{} {
	if msg.Message == nil {
		return nil
	}

	// Cek media langsung di pesan
	if img := msg.Message.GetImageMessage(); img != nil {
		return img
	}
	if vid := msg.Message.GetVideoMessage(); vid != nil {
		return vid
	}

	// Cek di quoted message (reply)
	extMsg := msg.Message.GetExtendedTextMessage()
	if extMsg == nil {
		return nil
	}

	quoted := extMsg.GetContextInfo().GetQuotedMessage()
	if quoted == nil {
		return nil
	}

	if img := quoted.GetImageMessage(); img != nil {
		return img
	}
	if vid := quoted.GetVideoMessage(); vid != nil {
		return vid
	}

	return nil
}

func (s *StickerHandler) downloadMedia(mediaMsg interface{}) ([]byte, error) {
	switch m := mediaMsg.(type) {
	case *waE2E.ImageMessage:
		return s.Client.Download(context.Background(), m)
	case *waE2E.VideoMessage:
		return s.Client.Download(context.Background(), m)
	default:
		return nil, fmt.Errorf("tipe media tidak didukung")
	}
}

// convertToWebP mengkonversi image/video/gif ke WebP menggunakan ffmpeg
// Semua konversi menggunakan ffmpeg agar konsisten dan mendukung rounded corners
func (s *StickerHandler) convertToWebP(data []byte, mimetype string, isVideo bool, doCrop bool) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "wa-sticker-*")
	if err != nil {
		return nil, fmt.Errorf("gagal buat temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Tentukan ekstensi berdasarkan mimetype
	ext := detectExtension(mimetype)

	inputPath := filepath.Join(tmpDir, "input"+ext)
	outputPath := filepath.Join(tmpDir, "output.webp")

	if err := os.WriteFile(inputPath, data, 0644); err != nil {
		return nil, fmt.Errorf("gagal tulis file input: %w", err)
	}

	// Build filter: Resize ke 512x512
	var vf string
	if doCrop {
		// Scale to cover 512x512, then crop center
		vf = fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=increase,crop=%d:%d", stickerSize, stickerSize, stickerSize, stickerSize)
	} else {
		// Scale to fit within 512x512
		vf = fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", stickerSize, stickerSize)
	}

	// Rounded Corners via geq: definisikan r, g, b, a agar FFmpeg tidak error
	// r, g, b memproses nilai warna asal, a memproses transparansi
	roundedVf := "format=rgba,geq=" +
		"r='r(X,Y)':" +
		"g='g(X,Y)':" +
		"b='b(X,Y)':" +
		"a='if(gt(abs(W/2-X),W/2-40)*gt(abs(H/2-Y),H/2-40),if(gt(hypot(X-(if(lt(X,W/2),40,W-40)),Y-(if(lt(Y,H/2),40,H-40))),40),0,255),255)'"
	vf = vf + "," + roundedVf

	var args []string
	if isVideo {
		args = []string{
			"-i", inputPath,
			"-t", "8",
			"-vf", vf,
			"-loop", "0",
			"-compression_level", "4",
			"-quality", "60",
			"-an",
			"-y",
			outputPath,
		}
	} else {
		// Image statis — single frame webp
		args = []string{
			"-i", inputPath,
			"-vf", vf,
			"-frames:v", "1",
			"-compression_level", "4",
			"-quality", "80",
			"-y",
			outputPath,
		}
	}

	cmd := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg gagal: %w\nstderr: %s", err, stderr.String())
	}

	result, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("gagal baca output webp: %w", err)
	}

	// WhatsApp stiker max ~1MB
	if len(result) > 1024*1024 {
		return nil, fmt.Errorf("stiker terlalu besar (%.1f KB), coba media dengan resolusi lebih kecil atau video lebih pendek", float64(len(result))/1024)
	}

	return result, nil
}

func (s *StickerHandler) sendSticker(ctx *command.Context, webpData []byte, isAnimated bool) error {
	uploaded, err := s.Client.Upload(context.Background(), webpData, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("gagal upload stiker: %w", err)
	}

	stickerMsg := &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           &uploaded.URL,
			DirectPath:    &uploaded.DirectPath,
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String("image/webp"),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(webpData))),
			IsAnimated:    proto.Bool(isAnimated),
		},
	}

	_, err = s.Client.SendMessage(context.Background(), ctx.Chat, stickerMsg)
	return err
}

func (s *StickerHandler) sendText(ctx *command.Context, text string) {
	_, err := s.Client.SendMessage(context.Background(), ctx.Chat, &waE2E.Message{
		Conversation: &text,
	})
	if err != nil {
		log.Printf("gagal kirim pesan: %v", err)
	}
}

func (s *StickerHandler) isReply(msg *events.Message) bool {
	if msg.Message == nil {
		return false
	}
	extMsg := msg.Message.GetExtendedTextMessage()
	return extMsg != nil && extMsg.ContextInfo != nil && extMsg.ContextInfo.QuotedMessage != nil
}


// detectExtension mengembalikan ekstensi file berdasarkan mimetype
func detectExtension(mimetype string) string {
	switch {
	case strings.Contains(mimetype, "jpeg"), strings.Contains(mimetype, "jpg"):
		return ".jpg"
	case strings.Contains(mimetype, "png"):
		return ".png"
	case strings.Contains(mimetype, "gif"):
		return ".gif"
	case strings.Contains(mimetype, "webp"):
		return ".webp"
	case strings.Contains(mimetype, "mp4"):
		return ".mp4"
	case strings.Contains(mimetype, "3gpp"):
		return ".3gp"
	default:
		return ".bin"
	}
}
