package fastweb

type HttpHandler func(ctx *HttpContext) error
