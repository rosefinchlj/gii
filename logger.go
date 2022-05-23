package gii

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		ctx.Next()
		log.Printf("[%d] %s in %s", ctx.StatusCode, ctx.Req.RequestURI, time.Since(start))
	}
}
