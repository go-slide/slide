package middleware

import (
	"compress/gzip"
	"ferry"
)

func Gzip() func (ctx *ferry.Ctx) error {
	return func(ctx *ferry.Ctx) error {
		writer, err := gzip.NewWriterLevel(ctx.Writer, gzip.BestCompression)
		if err != nil {
			return err
		}
		ctx.Writer.Header().Set("Content-Encoding", "gzip")
		defer writer.Close()
		// instead write own writer and implement common function
		// but will do that later
		ctx.GzipWriter = writer
		return ctx.Next()
	}
}