package middleware

import (
	"github.com/yebaoping/fastweb"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int64
}

var DefaultCORSConfig = CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{fasthttp.MethodGet, fasthttp.MethodHead, fasthttp.MethodPut, fasthttp.MethodPatch, fasthttp.MethodPost, fasthttp.MethodDelete},
}

func CORS(c CORSConfig) fastweb.Middleware {
	return func(next fastweb.HttpHandler) fastweb.HttpHandler {
		return func(ctx *fastweb.HttpContext) error {
			resp := &ctx.Response

			origin := ctx.Header(fasthttp.HeaderOrigin)
			allowOrigin := ""

			for _, o := range c.AllowOrigins {
				if o == "*" && c.AllowCredentials {
					allowOrigin = origin
					break
				}
				if o == "*" || o == origin {
					allowOrigin = o
					break
				}
			}

			if ctx.Method() != fasthttp.MethodOptions {
				resp.Header.Add(fasthttp.HeaderVary, fasthttp.HeaderOrigin)
				resp.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, allowOrigin)

				if c.AllowCredentials {
					resp.Header.Set(fasthttp.HeaderAccessControlAllowCredentials, "true")
				}

				if exposeHeaders := strings.Join(c.ExposeHeaders, ","); exposeHeaders != "" {
					resp.Header.Set(fasthttp.HeaderAccessControlExposeHeaders, exposeHeaders)
				}

				return next(ctx)
			}

			resp.Header.Add(fasthttp.HeaderVary, fasthttp.HeaderOrigin)
			resp.Header.Add(fasthttp.HeaderVary, fasthttp.HeaderAccessControlRequestMethod)
			resp.Header.Add(fasthttp.HeaderVary, fasthttp.HeaderAccessControlRequestHeaders)
			resp.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, allowOrigin)
			resp.Header.Set(fasthttp.HeaderAccessControlAllowMethods, strings.Join(c.AllowMethods, ","))

			if c.AllowCredentials {
				resp.Header.Set(fasthttp.HeaderAccessControlAllowCredentials, "true")
			}

			if allowHeaders := strings.Join(c.AllowHeaders, ","); allowHeaders != "" {
				resp.Header.Set(fasthttp.HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				accessControl := ctx.Header(fasthttp.HeaderAccessControlRequestHeaders)
				if accessControl != "" {
					resp.Header.Set(fasthttp.HeaderAccessControlAllowHeaders, accessControl)
				}
			}

			if c.MaxAge > 0 {
				resp.Header.Set(fasthttp.HeaderAccessControlMaxAge, strconv.FormatInt(c.MaxAge, 10))
			}

			return ctx.NoContent()
		}
	}
}
