package cliexec

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

func Address() string {
	return fmt.Sprintf("http://%s:%d", "localhost", 2310)
}

func Get(header []byte, route string) ([]byte, int, error) {
	code, d, err := fasthttp.Get(header, fmt.Sprintf("%s/%s", Address(), route))
	return d, code, err
}

func Post(header []byte, route string, args *fasthttp.Args) ([]byte, int, error) {
	code, d, err := fasthttp.Post(header, fmt.Sprintf("%s/%s", Address(), route), args)
	return d, code, err
}
