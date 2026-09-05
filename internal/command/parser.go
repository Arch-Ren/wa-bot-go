package command

import "strings"

type ParsedCommand struct {
	Prefix string
	Name   string
	Args   []string
}

func Parse(text string) *ParsedCommand {
	text = strings.TrimSpace(text)

	if text == "" {
		return nil
	}

	prefix := string(text[0])

	if prefix != "." && prefix != "/" {
		return nil
	}

	parts := strings.Fields(text[1:])

	if len(parts) == 0 {
		return nil
	}

	return &ParsedCommand{
		Prefix: prefix,
		Name:   parts[0],
		Args:   parts[1:],
	}
}
