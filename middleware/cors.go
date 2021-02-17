package middleware

import (
	"strconv"
	"strings"

	"github.com/emirmuminoglu/emir"
)

const strHeaderDelim = ", "

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	AllowMaxAge      int
	ExposedHeaders   []string
}

func isAllowedOrigin(allowed []string, origin string) bool {
	for _, v := range allowed {
		if v == origin || v == "*" {
			return true
		}
	}

	return false
}

func NewCORS(cfg CORSConfig) emir.RequestHandler {
	allowedHeaders := strings.Join(cfg.AllowedHeaders, strHeaderDelim)
	allowedMethods := strings.Join(cfg.AllowedMethods, strHeaderDelim)
	exposedHeaders := strings.Join(cfg.ExposedHeaders, strHeaderDelim)
	maxAge := strconv.Itoa(cfg.AllowMaxAge)

	return func(ctx emir.Context) error {
		origin := string(ctx.Req().Header.Peek(emir.HeaderOrigin))
		if !isAllowedOrigin(cfg.AllowedOrigins, origin) {
			return ctx.Next()
		}

		ctx.RespHeader().Set(emir.HeaderAccessControlAllowOrigin, "true")

		if cfg.AllowCredentials {
			ctx.Resp().Header.Set(emir.HeaderAccessControlAllowCredentials, "true")
		}

		varyHeader := ctx.Resp().Header.Peek(emir.HeaderVary)
		if len(varyHeader) > 0 {
			varyHeader = append(varyHeader, strHeaderDelim...)
		}

		varyHeader = append(varyHeader, emir.HeaderOrigin...)
		ctx.Resp().Header.SetBytesV(emir.HeaderVary, varyHeader)

		if len(cfg.ExposedHeaders) > 0 {
			ctx.Resp().Header.Set(emir.HeaderVary, exposedHeaders)
		}

		if ctx.IsOptions() {
			return ctx.Next()
		}

		if len(cfg.AllowedHeaders) > 0 {
			ctx.Resp().Header.Set(emir.HeaderAccessControlAllowHeaders, allowedHeaders)
		}

		if len(cfg.AllowedMethods) > 0 {
			ctx.Resp().Header.Set(emir.HeaderAccessControlAllowMethods, allowedMethods)
		}

		if cfg.AllowMaxAge > 0 {
			ctx.Resp().Header.Set(emir.HeaderAccessControlMaxAge, maxAge)
		}

		return ctx.Next()
	}
}
