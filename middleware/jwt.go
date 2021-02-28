package middleware

import (
	"bytes"
	"sync"

	"github.com/emirmuminoglu/emir"
	"github.com/emirmuminoglu/jwt"
)

var spaceByte = []byte(" ")

type JWTConfig struct {
	TokenLookupIn   string
	TokenLookupName string
	Key             []byte
	Algo            string
	AuthScheme      []byte
	ClaimsKey       string
}

func NewJWT(cfg JWTConfig) emir.RequestHandler {
	var extractor func(ctx emir.Context) []byte
	var parser func(key, token []byte, claims *jwt.Claims) error
	var pool sync.Pool
	var zeroClaims = &jwt.Claims{}

	acqClaims := func() *jwt.Claims {
		v := pool.Get()
		if v == nil {
			return new(jwt.Claims)
		}

		return v.(*jwt.Claims)
	}

	relClaims := func(claims *jwt.Claims) {
		*claims = *zeroClaims

		pool.Put(claims)
	}
	switch cfg.Algo {
	case "hs256":
		parser = jwt.ParseHS256
	case "hs384":
		parser = jwt.ParseHS384
	case "hs512":
		parser = jwt.ParseHS512
	}

	switch cfg.TokenLookupIn {
	case "header":
		extractor = func(ctx emir.Context) []byte {
			rawHeader := ctx.ReqHeader().Peek(cfg.TokenLookupName)
			splitted := bytes.Split(rawHeader, spaceByte)
			if len(splitted) != 2 {
				return nil
			}

			if bytes.Equal(splitted[0], cfg.AuthScheme) {
				return nil
			}

			return splitted[1]
		}
	case "cookie":
		extractor = func(ctx emir.Context) []byte {
			return ctx.ReqHeader().Cookie(cfg.TokenLookupName)
		}
	}

	return func(c emir.Context) error {
		token := extractor(c)
		if token == nil {
			return emir.NewBasicError(401, "missing token")
		}

		claims := acqClaims()
		defer relClaims(claims)
		err := parser(cfg.Key, token, claims)
		if err != nil {
			return emir.NewBasicError(401, "malformed token")
		}

		c.SetUserValue(cfg.ClaimsKey, claims)
		return c.Next()
	}
}

type JWTWithCustomConfig struct {
	TokenLookupIn   string
	TokenLookupName string
	Key             []byte
	Algo            string
	AuthScheme      []byte
	ClaimsKey       string
	ClaimFactory    func() interface{}
	ClaimReleaser   func(interface{})
	Validator       jwt.ValidatorFunction
}

func NewJWTWithCustomClaims(cfg JWTWithCustomConfig) emir.RequestHandler {
	var extractor func(ctx emir.Context) []byte
	var parser func(key, token []byte, claims interface{}, validator jwt.ValidatorFunction) error

	switch cfg.Algo {
	case "hs256":
		parser = jwt.ParseHS256Custom
	case "hs384":
		parser = jwt.ParseHS384Custom
	case "hs512":
		parser = jwt.ParseHS512Custom
	}

	switch cfg.TokenLookupIn {
	case "header":
		extractor = func(ctx emir.Context) []byte {
			rawHeader := ctx.ReqHeader().Peek(cfg.TokenLookupName)
			splitted := bytes.Split(rawHeader, spaceByte)
			if len(splitted) != 2 {
				return nil
			}

			if bytes.Equal(splitted[0], cfg.AuthScheme) {
				return nil
			}

			return splitted[1]
		}
	case "cookie":
		extractor = func(ctx emir.Context) []byte {
			return ctx.ReqHeader().Cookie(cfg.TokenLookupName)
		}
	}

	return func(c emir.Context) error {
		token := extractor(c)
		if token == nil {
			return emir.NewBasicError(401, "missing token")
		}

		claims := cfg.ClaimFactory()
		defer cfg.ClaimReleaser(claims)
		err := parser(cfg.Key, token, claims, cfg.Validator)
		if err != nil {
			return emir.NewBasicError(401, "malformed token")
		}

		c.SetUserValue(cfg.ClaimsKey, claims)
		return c.Next()
	}
}
