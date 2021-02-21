package emir

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func Test_JSONBind(t *testing.T) {
	e := New(Config{})

	e.POST("/", func(c Context) error {
		v := struct {
			Test  string `json:"test"`
			Test1 string `qs:"Test1"`
		}{}

		if err := c.Bind(&v); err != nil {
			t.Fatal(err)
		}

		if v.Test != "test" || v.Test1 != "test" {
			t.Fatalf("unexpected value. expected: %s, got: %s and %s", "test", v.Test, v.Test1)
		}

		return nil
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI("/?Test1=test")
	ctx.Request.Header.SetContentType(ContentTypeApplicationJSON)
	ctx.Request.SetBodyString(`{"test": "test"}`)
	ctx.Request.Header.SetMethod(MethodPost)

	e.Handler()(ctx)
}

func Benchmark_JSONBind(b *testing.B) {
	e := New(Config{})

	e.POST("/bench", func(c Context) error {
		v := struct {
			Test  string `json:"Test"`
			Test1 string `qs:"Test1"`
		}{}

		if err := c.Bind(&v); err != nil {
			b.Fatal(err)
		}

		return nil
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI("/bench?Test1=test")
	ctx.Request.Header.SetContentType(ContentTypeApplicationJSON)
	ctx.Request.SetBodyString(`{"test": "test"}`)
	ctx.Request.Header.SetMethod(MethodPost)

	handler := e.Handler()
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			handler(ctx)
		}
	})
}

func Benchmark_XMLBind(b *testing.B) {
	e := New(Config{})

	e.POST("/", func(c Context) error {
		v := struct {
			Test1 string `qs:"Test1"`
			Test  string `xml:"test"`
		}{}

		if err := c.Bind(&v); err != nil {
			b.Fatal(err)
		}

		return nil
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI("/?Test1=test")
	ctx.Request.Header.SetContentType(ContentTypeApplicationXML)
	ctx.Request.SetBodyString(`
		<?xml version="1.0" encoding="UTF-8"?>
		<root>
			<test>test</test>
		</root>
	`)
	ctx.Request.Header.SetMethod(MethodPost)

	handler := e.Handler()
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			handler(ctx)
		}
	})
}

func Test_XMLBind(t *testing.T) {
	e := New(Config{})

	e.POST("/", func(c Context) error {
		v := struct {
			Test1 string `qs:"Test1"`
			Test  string `xml:"test"`
		}{}

		if err := c.Bind(&v); err != nil {
			t.Fatal(err)
		}
		if v.Test != "test" || v.Test1 != "test" {
			t.Fatalf("unexpected value. expected: %s, got: %s and %s", "test", v.Test, v.Test1)
		}

		return nil
	})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI("/?Test1=test")
	ctx.Request.Header.SetContentType(ContentTypeApplicationXML)
	ctx.Request.SetBodyString(`
		<?xml version="1.0" encoding="UTF-8"?>
		<root>
			<test>test</test>
		</root>
	`)
	ctx.Request.Header.SetMethod(MethodPost)

	e.Handler()(ctx)
}
