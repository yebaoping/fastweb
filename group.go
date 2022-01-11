package fastweb

import (
	"net/http"
	"strings"
)

type group struct {
	app         *App
	prefix      string
	middlewares []Middleware
}

func (a *App) Group(prefix string, middlewares ...Middleware) *group {
	prefix = strings.TrimSuffix(prefix, "/")

	return &group{
		app:         a,
		prefix:      prefix,
		middlewares: middlewares,
	}
}

func (g *group) combinePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return g.prefix + path
}

func (g *group) combineMiddleware(handler HttpHandler, middlewares ...Middleware) HttpHandler {
	/*
		l := len(g.middlewares)
		for l > 0 {
			handler = g.middlewares[l-1](handler)

			l--
		}

		return handler
	*/
	s := make([]Middleware, 0, len(middlewares)+len(g.middlewares))
	s = append(s, g.middlewares...)
	s = append(s, middlewares...)

	return combineMiddleware(handler, s...)
}

func (g *group) combine(path string, handler HttpHandler, middlewares ...Middleware) (string, HttpHandler) {
	path = g.combinePath(path)
	handler = g.combineMiddleware(handler, middlewares...)

	return path, handler
}

func (g *group) Handle(path, method string, handler HttpHandler, middlewares ...Middleware) {
	path, handler = g.combine(path, handler, middlewares...)
	g.app.Handle(path, method, handler)
}

func (g *group) Static(path, file, contentType string, data interface{}, middlewares ...Middleware) {
	path = g.combinePath(path)

	s := make([]Middleware, 0, len(middlewares)+len(g.middlewares))
	s = append(s, g.middlewares...)
	s = append(s, middlewares...)

	g.app.Static(path, file, contentType, data, s...)
}

func (g *group) Html(path, file string, data interface{}, middlewares ...Middleware) {
	g.Static(path, file, "text/html", data, middlewares...)
}

func (g *group) Get(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodGet, handler, middlewares...)
}

func (g *group) Post(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodPost, handler, middlewares...)
}

func (g *group) Put(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodPut, handler, middlewares...)
}

func (g *group) Delete(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodDelete, handler, middlewares...)
}

func (g *group) Head(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodHead, handler, middlewares...)
}

func (g *group) Patch(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodPatch, handler, middlewares...)
}

func (g *group) Options(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodOptions, handler, middlewares...)
}

func (g *group) Connect(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodConnect, handler, middlewares...)
}

func (g *group) Trace(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodTrace, handler, middlewares...)
}
