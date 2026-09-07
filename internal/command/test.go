package command

import "log"

func PingHandler(ctx *Context) {
	log.Printf("ping command from %s", ctx.Sender.String())
}
