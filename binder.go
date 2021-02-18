package emir

import (
	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/pasztorpisti/qs"
)

type Binder interface {
	Bind(c Context, v interface{}) error
}

type DefaultBinder struct {
}

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
			if err := qs.Unmarshal(v, B2S(c.PostArgs().QueryString())); err != nil {
				return err
			}
		}
	}

	return qs.Unmarshal(v, B2S(c.QueryArgs().QueryString()))
}
