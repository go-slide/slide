package middleware

import (
	"github.com/slide-go/slide"
	"github.com/valyala/fasthttp"
)

// Compress brotli compression
func Compress() func(ctx *slide.Ctx) error {
	// give user config to use level // TODO
	compressHandler := fasthttp.CompressHandlerBrotliLevel(func(c *fasthttp.RequestCtx) {}, fasthttp.CompressBrotliBestSpeed, fasthttp.CompressBestSpeed)
	return func(ctx *slide.Ctx) error {
		compressHandler(ctx.RequestCtx)
		return ctx.Next()
	}
}
