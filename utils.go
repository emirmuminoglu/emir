package emir

import (
	"net/url"
	"reflect"
	"unsafe"

	"github.com/valyala/fasthttp"
)

// B2S converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2B converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func S2B(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len

	return
}

// ConvertArgsToValues converts given fasthttp.Args to url.Values
func ConvertArgsToValues(args *fasthttp.Args) url.Values {
	var values url.Values
	args.VisitAll(func(key, value []byte) {
		keyStr := B2S(key)
		values[keyStr] = append(values[keyStr], B2S(value))
	})

	return values
}
