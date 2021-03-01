package middleware

import (
	"github.com/emirmuminoglu/emir"
	"github.com/google/uuid"
)

func NewRequestID() emir.RequestHandler {
	return func(c *emir.Context) error {
		c.ReqHeader().Set(emir.HeaderXRequestID, uuid.New().String())

		return c.Next()
	}
}
