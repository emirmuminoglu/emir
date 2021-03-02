# emir

[![Build Status](https://travis-ci.com/emirmuminoglu/emir.svg?branch=v0.0.3)](https://travis-ci.com/emirmuminoglu/emir)
[![Go Report Card](https://goreportcard.com/badge/github.com/emirmuminoglu/emir)](https://goreportcard.com/report/github.com/emirmuminoglu/emir)
[![codecov](https://codecov.io/gh/emirmuminoglu/emir/branch/master/graph/badge.svg?token=M0IH7CMZNS)](https://codecov.io/gh/emirmuminoglu/emir)
[![Go Reference](https://pkg.go.dev/badge/github.com/emirmuminoglu/emir.svg)](https://pkg.go.dev/github.com/emirmuminoglu/emir)

A lightweight, high performance micro web framework.

It's based on [FastHTTP](https://github.com/valyala/fasthttp) and inspired by [atreugo](https://github.com/savsgio/atreugo) and [echo](https://github.com/labstack).

# Feature Overview

- Based on [fasthttp/router](https://github.com/fasthttp/router).
- It's based on [FastHTTP](https://github.com/valyala/fasthttp). It's faster up to 10 times faster than net/http.
- Uses [uber/zap](https://github.com/uber-go/zap).
- Path parameters.
- Multiple handlers to single route.
- Before and after middlewares to a router or to a specific route.
- Define an error handler to a router or to a specific route.
- Data binding for JSON, XML, form and query payload.
- Customizable Request Context.
- Common HTTP responses like JSON, HTML, plain text.
