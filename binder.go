package emir

import (
	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/pasztorpisti/qs"
)

// Binder is the interface that wraps the Bind method.
type Binder interface {
	Bind(c Context, v interface{}) error
}

// DefaultBinder is the default implementation of the Binder interface.
type DefaultBinder struct {
}

// Bind implements the Binder#Bind function. Binding is done in following order:
// Binder will bind the body first then binds the query params
// If the request method is not POST, PUT or PATCH then the binder will skip the body
// Struct tag for query params will be "qs".
func (*DefaultBinder) Bind(c Context, v interface{}) error {
	req := c.Req()
	contentType := B2S(req.Header.ContentType())

	if req.Header.IsPost() || req.Header.IsPut() || req.Header.IsPatch() {

		switch {
		case strings.HasPrefix(contentType, ContentTypeApplicationJSON):
			if err := json.Unmarshal(c.PostBody(), v); err != nil {
				return err
			}
		case strings.HasPrefix(contentType, ContentTypeApplicationXML) ||
			strings.HasPrefix(contentType, ContentTypeTextXML):

			if err := xml.Unmarshal(c.PostBody(), v); err != nil {
				return err
			}
		case strings.HasPrefix(contentType, ContentTypeApplicationForm):
			if err := qs.UnmarshalValues(v, ConvertArgsToValues(c.QueryArgs())); err != nil {
				return err
			}
		}
	}

	return qs.UnmarshalValues(v, ConvertArgsToValues(c.QueryArgs()))
}
