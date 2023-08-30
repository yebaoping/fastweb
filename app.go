package fastweb

import (
	"net"
	"net/http"
	"os"

	"html/template"

	"github.com/valyala/fasthttp"
	"github.com/yebaoping/fastweb/pkg/str"
)

type App struct {
	server       *fasthttp.Server
	router       Router
	log          Logger
	binder       Binder
	middlewares  []Middleware
	ErrorHandler func(ctx *HttpContext, err error)
}

func NewApp(l Logger, r Router, b Binder) *App {
	a := &App{
		server:      &fasthttp.Server{},
		router:      r,
		log:         l,
		binder:      b,
		middlewares: make([]Middleware, 0),
	}

	if a.log == nil {
		a.log = DefaultLogger(os.Stdout, false)
	}

	if a.router == nil {
		a.router = DefaultRouter()
	}

	if a.binder == nil {
		a.binder = DefaultBinder()
	}

	a.server.Handler = a.handleRequest()

	return a
}

var fmtStr = "fasthttp server started on [%s]"

func (a *App) Serve(ln net.Listener) error {
	a.log.Infof(fmtStr, ln.Addr().String())
	return a.server.Serve(ln)
}

func (a *App) ServeTLS(ln net.Listener, certFile, keyFile string) error {
	a.log.Infof(fmtStr, ln.Addr().String())
	return a.server.ServeTLS(ln, certFile, keyFile)
}

func (a *App) ServeTLSEmbed(ln net.Listener, certData, keyData []byte) error {
	a.log.Infof(fmtStr, ln.Addr().String())
	return a.server.ServeTLSEmbed(ln, certData, keyData)
}

func (a *App) Start(addr string) error {
	a.log.Infof(fmtStr, addr)
	return a.server.ListenAndServe(addr)
}

func (a *App) StartTLS(addr, certFile, keyFile string) error {
	a.log.Infof(fmtStr, addr)
	return a.server.ListenAndServeTLS(addr, certFile, keyFile)
}

func (a *App) StartTLSEmbed(addr string, certData, keyData []byte) error {
	a.log.Infof(fmtStr, addr)
	return a.server.ListenAndServeTLSEmbed(addr, certData, keyData)
}

func (a *App) Shutdown() error {
	a.log.Info("fasthttp server shutdown")
	return a.server.Shutdown()
}

func (a *App) Use(middlewares ...Middleware) {
	a.middlewares = append(a.middlewares, middlewares...)
}

// 包装层次：ErrorHandler -> App.middlewares -> Group.middlewares -> HttpHandler middlewares -> HttpHandler
func (a *App) combineMiddleware(handler HttpHandler, middlewares ...Middleware) HttpHandler {
	/*
		l := len(a.middlewares)
		for l > 0 {
			handler = a.middlewares[l-1](handler)

			l--
		}

		return handler
	*/
	s := make([]Middleware, 0, len(middlewares)+len(a.middlewares))
	s = append(s, a.middlewares...)
	s = append(s, middlewares...)

	return combineMiddleware(handler, s...)
}

func (a *App) Handle(path, method string, handler HttpHandler, middlewares ...Middleware) {
	handler = a.combineMiddleware(handler, middlewares...)

	a.router.RegisterRoute(path, method, func(ctx *HttpContext) error {
		a.log.Debugf("%s request,path:%s", method, path)

		return handler(ctx)
	})
}

func (a *App) Static(path, file, contentType string, data interface{}, middlewares ...Middleware) {
	handler := a.combineMiddleware(func(ctx *HttpContext) error {
		t, err := template.ParseFiles(file)
		if err != nil {
			return err
		}

		ctx.SetContentType(contentType)

		return t.Execute(ctx, data)
	}, middlewares...)

	a.router.RegisterRoute(path, http.MethodGet, func(ctx *HttpContext) error {
		a.log.Debugf("static request,path:%s", path)

		return handler(ctx)
	})
}

func (a *App) Html(path, file string, data interface{}, middlewares ...Middleware) {
	a.Static(path, file, ContentTypeTextHtml, data, middlewares...)
}

func (a *App) Get(path string, handler HttpHandler, middlewares ...Middleware) {
	a.Handle(path, http.MethodGet, handler, middlewares...)
}

func (a *App) Post(path string, handler HttpHandler, middlewares ...Middleware) {
	a.Handle(path, http.MethodPost, handler, middlewares...)
}

func (a *App) Put(path string, handler HttpHandler, middlewares ...Middleware) {
	a.Handle(path, http.MethodPut, handler, middlewares...)
}

func (a *App) Delete(path string, handler HttpHandler, middlewares ...Middleware) {
	a.Handle(path, http.MethodDelete, handler, middlewares...)
}

func (a *App) Head(path string, handler HttpHandler, middlewares ...Middleware) {
	a.Handle(path, http.MethodHead, handler, middlewares...)
}

func (a *App) Patch(path string, handler HttpHandler, middlewares ...Middleware) {
	a.Handle(path, http.MethodPatch, handler, middlewares...)
}

func (a *App) Options(path string, handler HttpHandler, middlewares ...Middleware) {
	a.Handle(path, http.MethodOptions, handler, middlewares...)
}

func (a *App) Connect(path string, handler HttpHandler, middlewares ...Middleware) {
	a.Handle(path, http.MethodConnect, handler, middlewares...)
}

func (a *App) Trace(path string, handler HttpHandler, middlewares ...Middleware) {
	a.Handle(path, http.MethodTrace, handler, middlewares...)
}

// 执行顺序：HttpHandler -> HttpHandler middlewares -> Group.middlewares -> App.middlewares -> ErrorHandler
func (a *App) handleRequest() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := str.String(ctx.Path())

		a.log.Debugf("request path:%s", path)

		httpCtx := &HttpContext{
			RequestCtx: ctx,
			Params:     make(map[string]string),
			binder:     a.binder,
		}

		// 根据path和http method获取对应的handler
		route, err := a.router.Route(httpCtx, path, httpCtx.Method())
		if err != nil || route == nil {
			switch err {
			case ErrMethodNotAllowed:
				a.log.Debugf("path:%s method:%s %s", path, httpCtx.Method(), err.Error())
				ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			case ErrNotFound:
				a.log.Debugf("path:%s %s", path, err.Error())
				fallthrough
			default:
				ctx.SetStatusCode(fasthttp.StatusNotFound)
			}

			return
		}

		//绑定参数
		err = a.router.Params(httpCtx, route.Path, path)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		// 执行handler
		err = route.Handler(httpCtx)
		if err != nil {
			if a.ErrorHandler != nil {
				a.ErrorHandler(httpCtx, err)
				return
			}

			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
	}
}
