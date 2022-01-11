package fastweb

import (
	"net/http"
	"reflect"
	"testing"
)

var r *defaultRouter

func TestMain(m *testing.M) {
	r = DefaultRouter()

	m.Run()
}

func Test_parseRoute(t *testing.T) {
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

	ctx := &HttpContext{
		Params: make(requestParams),
	}

	for regRoute, pathMethodHandler := range routePathMethodHandler {
		for path, methodHandler := range pathMethodHandler {
			for method := range methodHandler {
				route, err := r.Route(ctx, path, method)
				if err != nil {
					t.Errorf("path:%s route error:%+v", path, err)
					continue
				}

				if regRoute != route.path {
					t.Errorf("routePath:%s,path:%s,route:%s", regRoute, path, route.path)
					continue
				}

				//t.Logf("routePath:%s,path:%s,route:%s", routePath, path, route.path)
			}
		}
	}
}

func Test_params(t *testing.T) {
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
