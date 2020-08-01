package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-slide/slide"
)

// CorsConfig configuration for Corsfeat
type CorsConfig struct {
	AllowMethods     []string
	AllowOrigins     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

var (
	// DefaultCORSConfig defeault config for cors
	DefaultCORSConfig = CorsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
)

// Cors Middleware with default config
// Reference https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
func Cors() func(ctx *slide.Ctx) error {
	return CorsWithConfig(DefaultCORSConfig)
}

// CorsWithConfig Cors with a config
func CorsWithConfig(config CorsConfig) func(ctx *slide.Ctx) error {
	return func(ctx *slide.Ctx) error {
		if len(config.AllowOrigins) == 0 {
			config.AllowOrigins = DefaultCORSConfig.AllowOrigins
		}
		if len(config.AllowMethods) == 0 {
			config.AllowMethods = DefaultCORSConfig.AllowMethods
		}
		allowMethods := strings.Join(config.AllowMethods, ",")
		allowHeaders := strings.Join(config.AllowHeaders, ",")
		exposeHeaders := strings.Join(config.ExposeHeaders, ",")
		maxAge := strconv.Itoa(config.MaxAge)
		origin := string(ctx.RequestCtx.Request.Header.Peek(slide.HeaderOrigin))
		allowOrigin := ""
		for _, o := range config.AllowOrigins {
			if o == "*" && config.AllowCredentials {
				allowOrigin = origin
				break
			}
			if o == "*" || o == origin {
				allowOrigin = o
				break
			}
		}
		if string(ctx.RequestCtx.Method()) != http.MethodOptions {
			ctx.RequestCtx.Response.Header.Set(slide.HeaderAccessControlAllowOrigin, allowOrigin)
			ctx.RequestCtx.Response.Header.Set(slide.HeaderVary, slide.HeaderOrigin)
			if config.AllowCredentials {
				ctx.RequestCtx.Response.Header.Set(slide.HeaderAccessControlAllowCredentials, "true")
			}
			if exposeHeaders != "" {
				ctx.RequestCtx.Response.Header.Set(slide.HeaderAccessControlExposeHeaders, exposeHeaders)
			}
			return ctx.Next()
		}
		// Options request
		ctx.RequestCtx.Response.Header.Set(slide.HeaderVary, slide.HeaderOrigin)
		ctx.RequestCtx.Response.Header.Set(slide.HeaderVary, slide.HeaderAccessControlRequestMethod)
		ctx.RequestCtx.Response.Header.Set(slide.HeaderVary, slide.HeaderAccessControlRequestHeaders)
		ctx.RequestCtx.Response.Header.Set(slide.HeaderAccessControlAllowOrigin, allowOrigin)
		ctx.RequestCtx.Response.Header.Set(slide.HeaderAccessControlAllowMethods, allowMethods)

		if config.AllowCredentials {
			ctx.RequestCtx.Response.Header.Set(slide.HeaderAccessControlAllowCredentials, "true")
		}
		if allowHeaders != "" {
			ctx.RequestCtx.Response.Header.Set(slide.HeaderAccessControlAllowHeaders, allowHeaders)
		} else {
			h := string(ctx.RequestCtx.Request.Header.Peek(slide.HeaderAccessControlRequestHeaders))
			if h != "" {
				ctx.RequestCtx.Response.Header.Set(slide.HeaderAccessControlAllowHeaders, h)
			}
		}
		if config.MaxAge > 0 {
			ctx.RequestCtx.Response.Header.Set(slide.HeaderAccessControlMaxAge, maxAge)
		}
		return ctx.SendStatusCode(http.StatusNoContent)
	}

}
