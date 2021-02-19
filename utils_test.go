package emir

import (
	"bytes"
	"testing"

	"github.com/valyala/fasthttp"
)

func Test_B2S(t *testing.T) {
	originalString := "test"
	originalBytes := []byte(originalString)

	if value := B2S(originalBytes); value != originalString {
		t.Errorf("result and expected is not equal, result: %v", value)
	}
}

func Test_S2B(t *testing.T) {
	str := "test"
	originalBytes := []byte(str)

	if value := S2B(str); !bytes.Equal(value, originalBytes) {
		t.Errorf("result and expected is not equal, result: %v", value)
	}

}

func Test_ConvertArgsToValues(t *testing.T) {
	args := new(fasthttp.Args)
	args.Set("test", "test1")
	args.Add("test", "test2")

	values := ConvertArgsToValues(args)

	if length := len(values["test"]); length != 2 {
		t.Errorf("unexpected map length, expected: 2, result: %v", length)
	}
}
