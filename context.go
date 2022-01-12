package fastweb

import (
	"github.com/valyala/fasthttp"
	"github.com/yebaoping/fastweb/pkg/json"
	"github.com/yebaoping/fastweb/pkg/str"
)

type HttpContext struct {
	*fasthttp.RequestCtx
	Params requestParams
	binder Binder
}

type requestParams map[string]string

func (rp requestParams) Get(name string) string {
	return rp[name]
}

func (ctx *HttpContext) Method() string {
	return str.String(ctx.RequestCtx.Method())
}

func (ctx *HttpContext) Path() string {
	return str.String(ctx.RequestCtx.Path())
}

func (ctx *HttpContext) Header(key string) string {
	return str.String(ctx.Request.Header.Peek(key))
}

func (ctx *HttpContext) Query(key string) string {
	return str.String(ctx.QueryArgs().Peek(key))
}

func (ctx *HttpContext) Bind(val interface{}) error {
	return ctx.binder.Bind(ctx, val)
}

func (ctx *HttpContext) Json(data interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx.Success(ContentTypeJson, bs)

	return nil
}

func (ctx *HttpContext) Error(err error) error {
	ctx.RequestCtx.Error(err.Error(), fasthttp.StatusInternalServerError)
	return nil
}

func (ctx *HttpContext) NoContent() error {
	ctx.Response.Reset()
	ctx.SetStatusCode(fasthttp.StatusNoContent)
	return nil
}
