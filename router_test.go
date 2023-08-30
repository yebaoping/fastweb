package fastweb

import (
	"net/http"
	"reflect"
	"testing"
)

var r *defaultRouter

func TestMain(m *testing.M) {
	m.Run()
}

func Test_parseRoute(t *testing.T) {
	r = DefaultRouter()

	pathRoute := map[string]string{
		"/":                       `^\/*$`,
		"/api/user":               `^/api/user\/*$`,
		"/api/user/":              `^/api/user\/*$`,
		"/api/user/:start/:stop":  `^/api/user/\w+/\w+\/*$`,
		"/api/user/:start/:stop/": `^/api/user/\w+/\w+\/*$`,
		"/api/user/:id":           `^/api/user/\w+\/*$`,
		"/api/user/:id/":          `^/api/user/\w+\/*$`,
		"/api/user/:id/edit":      `^/api/user/\w+/edit\/*$`,
		"/api/user/:id/edit/":     `^/api/user/\w+/edit\/*$`,
	}

	for path, route := range pathRoute {
		if route != r.parseRoute(path) {
			t.Errorf("parse route failed,path:%s,route:%s,parseRoute:%s", path, route, r.parseRoute(path))
		}
	}
}

func Test_route(t *testing.T) {
	r = DefaultRouter()

	// 路由/路径/HTTP Method
	routePathMethodHandler := map[string]map[string]map[string]HttpHandler{
		"/": {
			"/": {
				http.MethodGet: func(ctx *HttpContext) error {
					return nil
				},
				http.MethodHead: func(ctx *HttpContext) error {
					return nil
				},
			},
		},
		"/user": {
			"/user": {
				http.MethodGet: func(ctx *HttpContext) error {
					return nil
				},
			},
			"/user/": {
				http.MethodGet: func(ctx *HttpContext) error {
					return nil
				}},
		},
		"/user/:id": {
			"/user/1": {
				http.MethodGet: func(ctx *HttpContext) error {
					return nil
				},
				http.MethodPost: func(ctx *HttpContext) error {
					return nil
				},
				http.MethodPut: func(ctx *HttpContext) error {
					return nil
				},
			},
			"/user/1/": {
				http.MethodGet: func(ctx *HttpContext) error {
					return nil
				},
				http.MethodPost: func(ctx *HttpContext) error {
					return nil
				},
				http.MethodPut: func(ctx *HttpContext) error {
					return nil
				},
			},
		},
		"/user/:start/:stop": {
			"/user/1/200": {
				http.MethodGet: func(ctx *HttpContext) error {
					return nil
				},
			},
			"/user/1/200/": {
				http.MethodGet: func(ctx *HttpContext) error {
					return nil
				},
			},
		},
	}

	for route, pathMethodHandler := range routePathMethodHandler {
		for _, methodHandler := range pathMethodHandler {
			for method, handler := range methodHandler {
				r.RegisterRoute(route, method, handler)
			}
		}
	}

	for regRoute, pathMethodHandler := range routePathMethodHandler {
		for path, methodHandler := range pathMethodHandler {
			for method := range methodHandler {
				ctx := &HttpContext{
					Params: make(requestParams),
				}
				route, err := r.Route(ctx, path, method)
				if err != nil {
					t.Errorf("path:%s route error:%+v", path, err)
					continue
				}

				if regRoute != route.Path {
					t.Errorf("routePath:%s,path:%s,route:%s", regRoute, path, route.Path)
					continue
				}

				//t.Logf("routePath:%s,path:%s,route:%s", routePath, path, route.Path)
			}
		}
	}
}

func Test_route_repeat(t *testing.T) {
	r = DefaultRouter()

	path := "/user"
	key := "value"
	val := ""
	ctx := &HttpContext{Params: map[string]string{
		key: "",
	}}

	val = "http get 1"
	r.RegisterRoute(path, http.MethodGet, func(ctx *HttpContext) error {
		ctx.Params[key] = val
		return nil
	})
	route, err := r.Route(ctx, path, http.MethodGet)
	if err != nil {
		t.Error(err)
	}
	err = route.Handler(ctx)
	if err != nil {
		t.Error(err)
	}
	if ctx.Params[key] != val {
		t.Fail()
	}

	val = "http get 2"
	r.RegisterRoute(path, http.MethodGet, func(ctx *HttpContext) error {
		ctx.Params[key] = val
		return nil
	})
	route, err = r.Route(ctx, path, http.MethodGet)
	if err != nil {
		t.Error(err)
	}
	err = route.Handler(ctx)
	if err != nil {
		t.Error(err)
	}
	if ctx.Params[key] != val {
		t.Fail()
	}
}

func Test_route_error(t *testing.T) {
	r = DefaultRouter()

	path := "/user"
	r.RegisterRoute(path, http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})

	ctx := &HttpContext{
		RequestCtx: nil,
		Params:     make(requestParams),
		binder:     nil,
	}

	// 路由不存在的情况
	route, err := r.Route(ctx, "/user/1", http.MethodGet)
	if err == nil || route != nil {
		t.Fail()
	}
	if err != ErrNotFound {
		t.Fail()
	}

	// http method 不匹配的情况
	route, err = r.Route(ctx, "/user", http.MethodPost)
	if err == nil || route != nil {
		t.Fail()
	}
	if err != ErrMethodNotAllowed {
		t.Fail()
	}
}

func Test_params(t *testing.T) {
	r = DefaultRouter()

	regPathPathParams := map[string]map[string]requestParams{
		"/": {
			"/": requestParams{},
		},
		"/api/user": {
			"/api/user": requestParams{},
		},
		"/api/user/:id": {
			"/api/user/1": {
				"id": "1",
			},
		},
		"/api/user/:start/:stop": {
			"/api/user/1/100": {
				"start": "1",
				"stop":  "100",
			},
		},
	}

	for regPath, pathParams := range regPathPathParams {
		for path, params := range pathParams {
			ctx := &HttpContext{
				Params: make(requestParams),
			}

			if err := r.Params(ctx, regPath, path); err != nil {
				t.Error(err)
			}

			if ctx.Params != nil && params != nil && !reflect.DeepEqual(ctx.Params, params) {
				t.Errorf("regPath:%s,path:%s,expectation:%+v,actual:%+v", regPath, path, params, ctx.Params)
			}
		}
	}
}
