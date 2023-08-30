package fastweb

import (
	"net/http"
	"strings"
)

type Group struct {
	app         *App
	prefix      string
	middlewares []Middleware
}

func (a *App) Group(prefix string, middlewares ...Middleware) *Group {
	prefix = strings.TrimSuffix(prefix, "/")

	return &Group{
		app:         a,
		prefix:      prefix,
		middlewares: middlewares,
	}
}

func (g *Group) combinePath(path string) string {
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return g.prefix + path
}

func (g *Group) combineMiddleware(handler HttpHandler, middlewares ...Middleware) HttpHandler {
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

func (g *Group) combine(path string, handler HttpHandler, middlewares ...Middleware) (string, HttpHandler) {
	path = g.combinePath(path)
	handler = g.combineMiddleware(handler, middlewares...)

	return path, handler
}

func (g *Group) Handle(path, method string, handler HttpHandler, middlewares ...Middleware) {
	path, handler = g.combine(path, handler, middlewares...)
	g.app.Handle(path, method, handler)
}

func (g *Group) Static(path, file, contentType string, data interface{}, middlewares ...Middleware) {
	path = g.combinePath(path)

	s := make([]Middleware, 0, len(middlewares)+len(g.middlewares))
	s = append(s, g.middlewares...)
	s = append(s, middlewares...)

	g.app.Static(path, file, contentType, data, s...)
}

func (g *Group) Html(path, file string, data interface{}, middlewares ...Middleware) {
	g.Static(path, file, ContentTypeTextHtml, data, middlewares...)
}

func (g *Group) Get(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodGet, handler, middlewares...)
}

func (g *Group) Post(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodPost, handler, middlewares...)
}

func (g *Group) Put(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodPut, handler, middlewares...)
}

func (g *Group) Delete(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodDelete, handler, middlewares...)
}

func (g *Group) Head(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodHead, handler, middlewares...)
}

func (g *Group) Patch(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodPatch, handler, middlewares...)
}

func (g *Group) Options(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodOptions, handler, middlewares...)
}

func (g *Group) Connect(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodConnect, handler, middlewares...)
}

func (g *Group) Trace(path string, handler HttpHandler, middlewares ...Middleware) {
	g.Handle(path, http.MethodTrace, handler, middlewares...)
}
