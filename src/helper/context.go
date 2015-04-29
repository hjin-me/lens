package helper

import (
	"net/http"

	"golang.org/x/net/context"
)

func Req(ctx context.Context) (*http.Request, bool) {
	v, ok := ctx.Value("req").(*http.Request)
	return v, ok
}
func Res(ctx context.Context) (http.ResponseWriter, bool) {
	v, ok := ctx.Value("res").(http.ResponseWriter)
	return v, ok
}
