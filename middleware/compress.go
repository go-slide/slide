package middleware

import (
	"ferry"

	"github.com/valyala/fasthttp"
)

// Compress brotli compression
func Compress() func(ctx *ferry.Ctx) error {
	// give user config to use level // TODO
	compressHandler := fasthttp.CompressHandlerBrotliLevel(func(c *fasthttp.RequestCtx) {}, fasthttp.CompressBrotliBestSpeed, fasthttp.CompressBestSpeed)
	return func(ctx *ferry.Ctx) error {
		compressHandler(ctx.RequestCtx)
		return ctx.Next()
	}
}
