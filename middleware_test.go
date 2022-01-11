package fastweb

import (
	"github.com/yebaoping/fastweb/pkg/str"
	"testing"
)

func Test_combineMiddleware(t *testing.T) {
	key := "values"

	h := func(ctx *HttpContext) error {
		ctx.Params[key], _ = str.Connect(ctx.Params[key], "1")
		return nil
	}

	m := []Middleware{
		func(next HttpHandler) HttpHandler {
			return func(ctx *HttpContext) error {
				ctx.Params[key], _ = str.Connect(ctx.Params[key], "3")
				return next(ctx)
			}
		},

		func(next HttpHandler) HttpHandler {
			return func(ctx *HttpContext) error {
				ctx.Params[key], _ = str.Connect(ctx.Params[key], "2")
				return next(ctx)
			}
		},
	}

	h = combineMiddleware(h, m...)

	ctx := HttpContext{
		Params: map[string]string{
			key: "",
		},
	}

	_ = h(&ctx)

	if ctx.Params[key] != "321" {
		t.Fail()
	}
}
