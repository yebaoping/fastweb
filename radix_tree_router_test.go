package fastweb

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

var regOrder = "random"

func Test_RadixTreeRouter(t *testing.T) {
	//regOrder = os.Getenv("regOrder")

	rtr := NewRadixTreeRouter()

	var routeMethodHandlerRequestPaths = []struct {
		route             string
		methodHandler     map[string]HttpHandler
		requestPathParams []struct {
			requestPath string
			Params      map[string]string
		}
	}{
		{
			route: "/",
			methodHandler: map[string]HttpHandler{
				http.MethodGet: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/" {
						return fmt.Errorf("route not match,expect:%s,result:/", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/",
				},
			},
		},
		{
			route: "/api/user",
			methodHandler: map[string]HttpHandler{
				http.MethodGet: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/user" {
						return fmt.Errorf("route not match,expect:%s,result:/api/user", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/api/user",
				},
			},
		},
		{
			route: "/api/user/add",
			methodHandler: map[string]HttpHandler{
				http.MethodPost: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/user/add" {
						return fmt.Errorf("route not match,expect:%s,result:/api/user/add", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/api/user/add",
				},
			},
		},
		{
			route: "/api/user/:id",
			methodHandler: map[string]HttpHandler{
				http.MethodGet: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/user/:id" {
						return fmt.Errorf("route not match,expect:%s,result:/api/user/:id", ctx.Params["test_match_route"])
					}
					return nil
				},
				http.MethodPut: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/user/:id" {
						return fmt.Errorf("route not match,expect:%s,result:/api/user/:id", ctx.Params["test_match_route"])
					}
					return nil
				},
				http.MethodDelete: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/user/:id" {
						return fmt.Errorf("route not match,expect:%s,result:/api/user/:id", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/api/user/2",
					Params: map[string]string{
						"id": "2",
					},
				},
				{
					requestPath: "/api/user/99",
					Params: map[string]string{
						"id": "99",
					},
				},
				{
					requestPath: "/api/user/100",
					Params: map[string]string{
						"id": "100",
					},
				},
				{
					requestPath: "/api/user/99999",
					Params: map[string]string{
						"id": "99999",
					},
				},
				{
					requestPath: "/api/user/kaxjqlhd",
					Params: map[string]string{
						"id": "kaxjqlhd",
					},
				},
			},
		},
		{
			route: "/api/user/:start/:stop",
			methodHandler: map[string]HttpHandler{
				http.MethodGet: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/user/:start/:stop" {
						return fmt.Errorf("route not match,expect:%s,result:/api/user/:start/:stop", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/api/user/0/10",
					Params: map[string]string{
						"start": "0",
						"stop":  "10",
					},
				},
				{
					requestPath: "/api/user/10/99",
					Params: map[string]string{
						"start": "10",
						"stop":  "99",
					},
				},
				{
					requestPath: "/api/user/100/999999",
					Params: map[string]string{
						"start": "100",
						"stop":  "999999",
					},
				},
			},
		},
		{
			route: "/api/user/:start/:stop/export",
			methodHandler: map[string]HttpHandler{
				http.MethodGet: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/user/:start/:stop/export" {
						return fmt.Errorf("route not match,expect:%s,result:/api/user/:start/:stop/export", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/api/user/0/9/export",
					Params: map[string]string{
						"start": "0",
						"stop":  "9",
					},
				},
				{
					requestPath: "/api/user/0/10/export",
					Params: map[string]string{
						"start": "0",
						"stop":  "10",
					},
				},
				{
					requestPath: "/api/user/10/100/export",
					Params: map[string]string{
						"start": "10",
						"stop":  "100",
					},
				},
				{
					requestPath: "/api/user/100/999999/export",
					Params: map[string]string{
						"start": "100",
						"stop":  "999999",
					},
				},
			},
		},
		{
			route: "/api/user/:id/edit",
			methodHandler: map[string]HttpHandler{
				http.MethodPost: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/user/:id/edit" {
						return fmt.Errorf("route not match,expect:%s,result:/api/user/:id/edit", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/api/user/2/edit",
					Params: map[string]string{
						"id": "2",
					},
				},
				{
					requestPath: "/api/user/99/edit",
					Params: map[string]string{
						"id": "99",
					},
				},
				{
					requestPath: "/api/user/100/edit",
					Params: map[string]string{
						"id": "100",
					},
				},
				{
					requestPath: "/api/user/99999/edit",
					Params: map[string]string{
						"id": "99999",
					},
				},
				{
					requestPath: "/api/user/kaxjqlhd/edit",
					Params: map[string]string{
						"id": "kaxjqlhd",
					},
				},
			},
		},
		{
			route: "/api/user/:id/delete",
			methodHandler: map[string]HttpHandler{
				http.MethodPost: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/user/:id/delete" {
						return fmt.Errorf("route not match,expect:%s,result:/api/user/:id/delete", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/api/user/2/delete",
					Params: map[string]string{
						"id": "2",
					},
				},
				{
					requestPath: "/api/user/99/delete",
					Params: map[string]string{
						"id": "99",
					},
				},
				{
					requestPath: "/api/user/100/delete",
					Params: map[string]string{
						"id": "100",
					},
				},
				{
					requestPath: "/api/user/99999/delete",
					Params: map[string]string{
						"id": "99999",
					},
				},
				{
					requestPath: "/api/user/kaxjqlhd/delete",
					Params: map[string]string{
						"id": "kaxjqlhd",
					},
				},
			},
		},
		{
			route: "/html/index.html",
			methodHandler: map[string]HttpHandler{
				http.MethodGet: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/html/index.html" {
						return fmt.Errorf("route not match,expect:%s,result:/html/index.html", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/html/index.html",
				},
			},
		},
		{
			route: "/html/user/index.html",
			methodHandler: map[string]HttpHandler{
				http.MethodGet: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/html/user/index.html" {
						return fmt.Errorf("route not match,expect:%s,result:/html/user/index.html", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/html/user/index.html",
				},
			},
		},
		{
			route: "/api/blog/add",
			methodHandler: map[string]HttpHandler{
				http.MethodPost: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/blog/add" {
						return fmt.Errorf("route not match,expect:%s,result:/api/blog/add", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/api/blog/add",
				},
			},
		},
		{
			route: "/api/users",
			methodHandler: map[string]HttpHandler{
				http.MethodPost: func(ctx *HttpContext) error {
					if ctx.Params["test_match_route"] != "/api/users" {
						return fmt.Errorf("route not match,expect:%s,result:/api/users", ctx.Params["test_match_route"])
					}
					return nil
				},
			},
			requestPathParams: []struct {
				requestPath string
				Params      map[string]string
			}{
				{
					requestPath: "/api/users",
				},
			},
		},
	}

	// 单个路由注册测试
	for _, routeMethodHandlerRequestPath := range routeMethodHandlerRequestPaths {
		for method, handler := range routeMethodHandlerRequestPath.methodHandler {
			rtr := NewRadixTreeRouter()
			rtr.RegisterRoute(routeMethodHandlerRequestPath.route, method, handler)

			for _, requestPathParam := range routeMethodHandlerRequestPath.requestPathParams {
				testRoute(t, rtr, routeMethodHandlerRequestPath.route, requestPathParam.requestPath, method, requestPathParam.Params)
			}
		}
	}

	if strings.EqualFold(regOrder, "desc") {
		// 逆序执行路由注册
		for i := len(routeMethodHandlerRequestPaths) - 1; i >= 0; i-- {
			routeMethodHandlerRequestPath := routeMethodHandlerRequestPaths[i]
			for method, handler := range routeMethodHandlerRequestPath.methodHandler {
				rtr.RegisterRoute(routeMethodHandlerRequestPath.route, method, handler)
			}
		}
	} else if strings.EqualFold(regOrder, "random") {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		l := len(routeMethodHandlerRequestPaths)

		// 随机打乱顺序注册
		for i := 0; i < l; i++ {
			no := r.Intn(l)
			routeMethodHandlerRequestPaths[i], routeMethodHandlerRequestPaths[no] = routeMethodHandlerRequestPaths[no], routeMethodHandlerRequestPaths[i]
		}
		for _, routeMethodHandlerRequestPath := range routeMethodHandlerRequestPaths {
			for method, handler := range routeMethodHandlerRequestPath.methodHandler {
				// t.Logf("register route:%s[%s]", routeMethodHandlerRequestPath.route, method)
				rtr.RegisterRoute(routeMethodHandlerRequestPath.route, method, handler)
			}
		}
	} else {
		// 顺序执行路由注册
		for _, routeMethodHandlerRequestPath := range routeMethodHandlerRequestPaths {
			for method, handler := range routeMethodHandlerRequestPath.methodHandler {
				rtr.RegisterRoute(routeMethodHandlerRequestPath.route, method, handler)
			}
		}
	}

	// 顺序执行路由发现
	for _, routeMethodHandlerRequestPath := range routeMethodHandlerRequestPaths {
		for method := range routeMethodHandlerRequestPath.methodHandler {
			for _, requestPathParam := range routeMethodHandlerRequestPath.requestPathParams {
				testRoute(t, rtr, routeMethodHandlerRequestPath.route, requestPathParam.requestPath, method, requestPathParam.Params)
			}
		}
	}

	// 逆序执行路由发现
	for i := len(routeMethodHandlerRequestPaths) - 1; i >= 0; i-- {
		routeMethodHandlerRequestPath := routeMethodHandlerRequestPaths[i]
		for method := range routeMethodHandlerRequestPath.methodHandler {
			for _, requestPathParam := range routeMethodHandlerRequestPath.requestPathParams {
				testRoute(t, rtr, routeMethodHandlerRequestPath.route, requestPathParam.requestPath, method, requestPathParam.Params)
			}
		}
	}

	// 错误路由测试
	ctx := &HttpContext{
		Params: map[string]string{},
	}
	route, err := rtr.Route(ctx, "/api/usr", http.MethodGet)
	if err != nil {
		t.Fatal(err)
	}
	if route != nil {
		t.Fatalf("/api/usr route exists,route:%s", route.Path)
	}

	route, err = rtr.Route(ctx, "/api/user/2", http.MethodPost)
	if !errors.Is(err, ErrMethodNotAllowed) {
		t.Fatalf("/api/user/2[POST] route Method mismatch,error:%+v", err)
	}
	if route != nil {
		t.Fatalf("/api/user/2[POST] route exists,route:%s", route.Path)
	}

	// 不推荐的路由写法测试
	rtr.RegisterRoute("/api/user/1", http.MethodGet, func(ctx *HttpContext) error {
		if ctx.Params["test_match_route"] != "/api/user/1" {
			return fmt.Errorf("route not match,expect:%s,result:/api/user/1", ctx.Params["test_match_route"])
		}
		return nil
	})

	route, err = rtr.Route(ctx, "/api/user/1/edit", http.MethodPost)
	if err != nil {
		t.Fatal(err)
	}
	if route != nil {
		t.Fatalf("/api/user/1/edit route exists,route:%s", route.Path)
	}

	route, err = rtr.Route(ctx, "/api/user/1/10", http.MethodGet)
	if err != nil {
		t.Fatal(err)
	}
	if route != nil {
		t.Fatalf("/api/user/1/10 route exists,route:%s", route.Path)
	}

	route, err = rtr.Route(ctx, "/api/user/1/10/export", http.MethodGet)
	if err != nil {
		t.Fatal(err)
	}
	if route != nil {
		t.Fatalf("/api/user/1/10/export route exists,route:%s", route.Path)
	}
}

// 特殊路由注册测试
func Test_RadixTreeRouter_SpecialRoute(t *testing.T) {
	rtr := NewRadixTreeRouter()
	rtr.RegisterRoute("/api/user/add", http.MethodPost, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/blog", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/users", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/user", http.MethodPost, func(ctx *HttpContext) error {
		return nil
	})

	testRoute(t, rtr, "/api/user/add", "/api/user/add", http.MethodPost, map[string]string{})
	testRoute(t, rtr, "/api/blog", "/api/blog", http.MethodGet, map[string]string{})
	testRoute(t, rtr, "/api/users", "/api/users", http.MethodGet, map[string]string{})
	testRoute(t, rtr, "/api/user", "/api/user", http.MethodPost, map[string]string{})

	//

	rtr = NewRadixTreeRouter()
	rtr.RegisterRoute("/api/user", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/html/user", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})

	testRoute(t, rtr, "/api/user", "/api/user", http.MethodGet, map[string]string{})
	testRoute(t, rtr, "/html/user", "/html/user", http.MethodGet, map[string]string{})
	testRoute(t, rtr, "/", "/", http.MethodGet, map[string]string{})

	//
	rtr = NewRadixTreeRouter()
	rtr.RegisterRoute("/api/users", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/user/add", http.MethodPost, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/user/:id/delete", http.MethodPost, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/user/:id/edit", http.MethodPost, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/user/:id/", http.MethodPost, func(ctx *HttpContext) error {
		return nil
	})
	testRoute(t, rtr, "/api/users", "/api/users", http.MethodGet, map[string]string{})
	testRoute(t, rtr, "/api/user/add", "/api/user/add", http.MethodPost, map[string]string{})
	testRoute(t, rtr, "/api/user/:id/edit", "/api/user/1/edit", http.MethodPost, map[string]string{
		"id": "1",
	})
	testRoute(t, rtr, "/api/user/:id/delete", "/api/user/1/delete", http.MethodPost, map[string]string{
		"id": "1",
	})
	testRoute(t, rtr, "/api/user/:id/", "/api/user/2/", http.MethodPost, map[string]string{
		"id": "2",
	})
}

// 特殊参数路由注册测试
func Test_RadixTreeRouter_SpecialParamsRoute(t *testing.T) {
	rtr := NewRadixTreeRouter()
	rtr.RegisterRoute("/api/user/:start/:stop", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/user/add", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})

	testRoute(t, rtr, "/api/user/:start/:stop", "/api/user/0/9", http.MethodGet, map[string]string{
		"start": "0",
		"stop":  "9",
	})
	testRoute(t, rtr, "/api/user/add", "/api/user/add", http.MethodGet, map[string]string{})

	//

	rtr = NewRadixTreeRouter()
	rtr.RegisterRoute("/api/user/:id", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/user/add", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})

	testRoute(t, rtr, "/api/user/:id", "/api/user/999", http.MethodGet, map[string]string{
		"id": "999",
	})
	testRoute(t, rtr, "/api/user/add", "/api/user/add", http.MethodGet, map[string]string{})

	//

	rtr = NewRadixTreeRouter()
	rtr.RegisterRoute("/api/user/:id", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/user", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/user/info", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})

	testRoute(t, rtr, "/api/user/:id", "/api/user/999", http.MethodGet, map[string]string{
		"id": "999",
	})
	testRoute(t, rtr, "/api/user", "/api/user", http.MethodGet, map[string]string{})
	testRoute(t, rtr, "/api/user/info", "/api/user/info", http.MethodGet, map[string]string{})

	//

	rtr = NewRadixTreeRouter()
	rtr.RegisterRoute("/api/user/:start/:stop", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})
	rtr.RegisterRoute("/api/user/:id", http.MethodGet, func(ctx *HttpContext) error {
		return nil
	})

	testRoute(t, rtr, "/api/user/:start/:stop", "/api/user/0/9", http.MethodGet, map[string]string{
		"start": "0",
		"stop":  "9",
	})
	testRoute(t, rtr, "/api/user/:id", "/api/user/99999", http.MethodGet, map[string]string{
		"id": "99999",
	})
}

// 测试相似路由
func Test_RadixTreeRouter_SimilarPath(t *testing.T) {
	rtr := NewRadixTreeRouter()

	paths := []string{
		"/api/user/:id",
		"/api/user/:id/delete",
		"/api/user/:start/:stop",
		"/api/user/:start/:stop/export",
	}
	for _, path := range paths {
		rtr.RegisterRoute(path, http.MethodGet, nil)
	}

	similarPaths := []string{
		"/api/user/:userid",
		"/api/user/:userid/delete",
		"/api/user/:begin/:end",
		"/api/user/:begin/:end/export",
	}
	for _, path := range similarPaths {
		t.Run(path, func(t *testing.T) {
			defer func() {
				err := recover()
				if err == nil {
					t.Error("similar path register success")
				} else {
					t.Log(err)
				}
			}()

			rtr.RegisterRoute(path, http.MethodGet, nil)
		})
	}
}

func Test_RadixTreeRouter_parsePathParams(t *testing.T) {
	pathParams := map[string]struct {
		parsedPath string
		isParam    bool
		params     []string
	}{
		":id":                                {parsedPath: ":", isParam: true, params: []string{":id"}},
		"/:id":                               {parsedPath: "/:", isParam: true, params: []string{":id"}},
		"/api/user/:id":                      {parsedPath: "/api/user/:", isParam: true, params: []string{":id"}},
		"/api/user/:id/edit":                 {parsedPath: "/api/user/:/edit", isParam: true, params: []string{":id"}},
		"/api/user/:start/:stop":             {parsedPath: "/api/user/:/:", isParam: true, params: []string{":start", ":stop"}},
		"/api/user/:page/:size/:sort":        {parsedPath: "/api/user/:/:/:", isParam: true, params: []string{":page", ":size", ":sort"}},
		"/api/user/:start/:stop/export":      {parsedPath: "/api/user/:/:/export", isParam: true, params: []string{":start", ":stop"}},
		"/api/user/:page/:size/:sort/export": {parsedPath: "/api/user/:/:/:/export", isParam: true, params: []string{":page", ":size", ":sort"}},
		"/api/user/:uid/role/:rid":           {parsedPath: "/api/user/:/role/:", isParam: true, params: []string{":uid", ":rid"}},
		"/api/user/:uid/role/:rid/delete":    {parsedPath: "/api/user/:/role/:/delete", isParam: true, params: []string{":uid", ":rid"}},
		"/api/user":                          {parsedPath: "/api/user", isParam: false, params: make([]string, 0)},
	}

	for path, param := range pathParams {
		parsedPath, isParam, params, err := parsePathParams(path)
		if err != nil {
			t.Fatalf("path:'%s' parse error:%+v", path, err)
		}
		if parsedPath != param.parsedPath || isParam != param.isParam || !reflect.DeepEqual(params, param.params) {
			t.Fatalf("path:%s parse fail,parsedPath:%s,isParam:%v,params:%+v", path, parsedPath, isParam, params)
		}
	}

	// 异常路径转换
	invalidPaths := []string{
		"u:id",
		"id:",
		"/u:id",
		"/id:",
		"/api/user/u:id",
		"/api/user/id:",
	}
	for _, path := range invalidPaths {
		_, _, _, err := parsePathParams(path)
		if err == nil {
			t.Fatalf("path:'%s' parse success", path)
		}
	}
}

func testRoute(t *testing.T, rtr *RadixTreeRouter, route, requestPath, method string, params map[string]string) {
	ctx := &HttpContext{
		Params: requestParams{
			"test_match_route": route,
		},
	}

	r, err := rtr.Route(ctx, requestPath, method)
	if err != nil {
		t.Fatalf("%s route error:%+v", requestPath, err)
	}
	if r == nil {
		t.Fatalf("%s Handler is null", requestPath)
	}

	if r.Path != route {
		t.Fatalf("route not match,request Path:%s,expect:%s,result:%s", requestPath, route, r.Path)
	}

	_ = rtr.Params(ctx, r.Path, requestPath)
	for key, val := range params {
		if ctx.Params[key] != val {
			t.Fatalf("param not equal.route:%s,request Path:%s,param name:%s,expect:%s,result:%s",
				route, requestPath, key, val, ctx.Params[key])
		}
	}

	err = r.Handler(ctx)
	if err != nil {
		t.Fatalf("Handler:%s execute error:%+v", r.Path, err)
	}
}

var routePathParams = []struct {
	route string
	paths []struct {
		path        string
		replaced    string
		paramValues []string
	}
}{
	{
		route: ":",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "1",
				replaced:    "1",
				paramValues: []string{"1"},
			},
			{
				path:        "999",
				replaced:    "999",
				paramValues: []string{"999"},
			},
			{
				path:        "1/edit",
				replaced:    "1",
				paramValues: []string{"1"},
			},
			{
				path:        "999999/edit",
				replaced:    "999999",
				paramValues: []string{"999999"},
			},
		},
	},
	{
		route: ":/edit",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "1/edit",
				replaced:    "1/edit",
				paramValues: []string{"1"},
			},
			{
				path:        "999/edit",
				replaced:    "999/edit",
				paramValues: []string{"999"},
			},
			{
				path:        "1",
				replaced:    "1/edit",
				paramValues: []string{"1"},
			},
			{
				path:        "99999",
				replaced:    "99999/edit",
				paramValues: []string{"99999"},
			},
		},
	},
	{
		route: "/:",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "/1",
				replaced:    "/1",
				paramValues: []string{"1"},
			},
			{
				path:        "/99999",
				replaced:    "/99999",
				paramValues: []string{"99999"},
			},
			{
				path:        "/1/edit",
				replaced:    "/1",
				paramValues: []string{"1"},
			},
			{
				path:        "/999999/edit",
				replaced:    "/999999",
				paramValues: []string{"999999"},
			},
		},
	},
	{
		route: "/:/edit",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "/1/edit",
				replaced:    "/1/edit",
				paramValues: []string{"1"},
			},
			{
				path:        "/99999/edit",
				replaced:    "/99999/edit",
				paramValues: []string{"99999"},
			},
			{
				path:        "/1",
				replaced:    "/1/edit",
				paramValues: []string{"1"},
			},
			{
				path:        "/99999",
				replaced:    "/99999/edit",
				paramValues: []string{"99999"},
			},
		},
	},
	{
		route: "/user/:",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "/user/1",
				replaced:    "/user/1",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/99999",
				replaced:    "/user/99999",
				paramValues: []string{"99999"},
			},
			{
				path:        "/user/1/edit",
				replaced:    "/user/1",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/99999/edit",
				replaced:    "/user/99999",
				paramValues: []string{"99999"},
			},
		},
	},
	{
		route: "/user/:/edit",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "/user/1/edit",
				replaced:    "/user/1/edit",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/999999/edit",
				replaced:    "/user/999999/edit",
				paramValues: []string{"999999"},
			},
			{
				path:        "/user/1",
				replaced:    "/user/1/edit",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/999999",
				replaced:    "/user/999999/edit",
				paramValues: []string{"999999"},
			},
		},
	},
	{
		route: "/user/:/:",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "/user/1/10",
				replaced:    "/user/1/10",
				paramValues: []string{"1", "10"},
			},
			{
				path:        "/user/1000/99999",
				replaced:    "/user/1000/99999",
				paramValues: []string{"1000", "99999"},
			},
			{
				path:        "/user/1",
				replaced:    "/user/1/:",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/1000",
				replaced:    "/user/1000/:",
				paramValues: []string{"1000"},
			},
			{
				path:        "/user/1/10/export",
				replaced:    "/user/1/10",
				paramValues: []string{"1", "10"},
			},
			{
				path:        "/user/1000/99999/export",
				replaced:    "/user/1000/99999",
				paramValues: []string{"1000", "99999"},
			},
		},
	},
	{
		route: "/user/:/:/:",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "/user/1/10/100",
				replaced:    "/user/1/10/100",
				paramValues: []string{"1", "10", "100"},
			},
			{
				path:        "/user/1000/99999/9",
				replaced:    "/user/1000/99999/9",
				paramValues: []string{"1000", "99999", "9"},
			},
			{
				path:        "/user/1",
				replaced:    "/user/1/:/:",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/1000",
				replaced:    "/user/1000/:/:",
				paramValues: []string{"1000"},
			},
			{
				path:        "/user/1/10",
				replaced:    "/user/1/10/:",
				paramValues: []string{"1", "10"},
			},
			{
				path:        "/user/1000/99999",
				replaced:    "/user/1000/99999/:",
				paramValues: []string{"1000", "99999"},
			},
			{
				path:        "/user/1/10/99/export",
				replaced:    "/user/1/10/99",
				paramValues: []string{"1", "10", "99"},
			},
			{
				path:        "/user/1000/99999/10/export",
				replaced:    "/user/1000/99999/10",
				paramValues: []string{"1000", "99999", "10"},
			},
		},
	},
	{
		route: "/user/:/role/:",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "/user/1/role/10",
				replaced:    "/user/1/role/10",
				paramValues: []string{"1", "10"},
			},
			{
				path:        "/user/1000/role/99999",
				replaced:    "/user/1000/role/99999",
				paramValues: []string{"1000", "99999"},
			},
			{
				path:        "/user/1/role/10/edit",
				replaced:    "/user/1/role/10",
				paramValues: []string{"1", "10"},
			},
			{
				path:        "/user/1000/role/99999/edit",
				replaced:    "/user/1000/role/99999",
				paramValues: []string{"1000", "99999"},
			},
			{
				path:        "/user/1",
				replaced:    "/user/1/role/:",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/1000",
				replaced:    "/user/1000/role/:",
				paramValues: []string{"1000"},
			},
			{
				path:        "/user/1/role",
				replaced:    "/user/1/role/:",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/1000/role",
				replaced:    "/user/1000/role/:",
				paramValues: []string{"1000"},
			},
		},
	},
	{
		route: "/user/:/role/:/:",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "/user/1/role/1/10",
				replaced:    "/user/1/role/1/10",
				paramValues: []string{"1", "1", "10"},
			},
			{
				path:        "/user/1000/role/10000/99999",
				replaced:    "/user/1000/role/10000/99999",
				paramValues: []string{"1000", "10000", "99999"},
			},
			{
				path:        "/user/1/role/1/10/edit",
				replaced:    "/user/1/role/1/10",
				paramValues: []string{"1", "1", "10"},
			},
			{
				path:        "/user/1000/role/10000/99999/edit",
				replaced:    "/user/1000/role/10000/99999",
				paramValues: []string{"1000", "10000", "99999"},
			},
			{
				path:        "/user/1",
				replaced:    "/user/1/role/:/:",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/999",
				replaced:    "/user/999/role/:/:",
				paramValues: []string{"999"},
			},
			{
				path:        "/user/1/role",
				replaced:    "/user/1/role/:/:",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/1000/role",
				replaced:    "/user/1000/role/:/:",
				paramValues: []string{"1000"},
			},
			{
				path:        "/user/1/role/999",
				replaced:    "/user/1/role/999/:",
				paramValues: []string{"1", "999"},
			},
			{
				path:        "/user/1000/role/999",
				replaced:    "/user/1000/role/999/:",
				paramValues: []string{"1000", "999"},
			},
		},
	},
	{
		route: "/user/:/:/export",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "/user/1/10/export",
				replaced:    "/user/1/10/export",
				paramValues: []string{"1", "10"},
			},
			{
				path:        "/user/1000/99999/export",
				replaced:    "/user/1000/99999/export",
				paramValues: []string{"1000", "99999"},
			},
			{
				path:        "/user/1",
				replaced:    "/user/1/:/export",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/1000",
				replaced:    "/user/1000/:/export",
				paramValues: []string{"1000"},
			},
			{
				path:        "/user/1/999",
				replaced:    "/user/1/999/export",
				paramValues: []string{"1", "999"},
			},
			{
				path:        "/user/1000/99999",
				replaced:    "/user/1000/99999/export",
				paramValues: []string{"1000", "99999"},
			},
		},
	},
	{
		route: "/user/:/:/:/export",
		paths: []struct {
			path        string
			replaced    string
			paramValues []string
		}{
			{
				path:        "/user/1/10/99/export",
				replaced:    "/user/1/10/99/export",
				paramValues: []string{"1", "10", "99"},
			},
			{
				path:        "/user/1000/99999/9/export",
				replaced:    "/user/1000/99999/9/export",
				paramValues: []string{"1000", "99999", "9"},
			},

			{
				path:        "/user/1",
				replaced:    "/user/1/:/:/export",
				paramValues: []string{"1"},
			},
			{
				path:        "/user/1000",
				replaced:    "/user/1000/:/:/export",
				paramValues: []string{"1000"},
			},
			{
				path:        "/user/1/99999",
				replaced:    "/user/1/99999/:/export",
				paramValues: []string{"1", "99999"},
			},
			{
				path:        "/user/1000/99999",
				replaced:    "/user/1000/99999/:/export",
				paramValues: []string{"1000", "99999"},
			},
		},
	},
}

func Test_RadixTreeRouter_replaceRouteParams(t *testing.T) {
	rtr := NewRadixTreeRouter()
	for _, routePathParam := range routePathParams {
		for _, pathParam := range routePathParam.paths {
			replaced, paramValues := rtr.replaceRouteParams(routePathParam.route, pathParam.path)
			if replaced != pathParam.replaced || !reflect.DeepEqual(pathParam.paramValues, paramValues) {
				t.Fatalf("route:'%s',Path:'%s',expectReplaced:'%s',replaced:'%s',expectParamValues:%+v,paramValues:%+v",
					routePathParam.route, pathParam.path, pathParam.replaced, replaced, pathParam.paramValues, paramValues)
			}
		}
	}
}

func Test_replaceRouteParams(t *testing.T) {
	for _, routePathParam := range routePathParams {
		for _, pathParam := range routePathParam.paths {
			replaced, paramValues := replaceRouteParams(routePathParam.route, pathParam.path)
			if replaced != pathParam.replaced || !reflect.DeepEqual(pathParam.paramValues, paramValues) {
				t.Fatalf("route:'%s',Path:'%s',expectReplaced:'%s',replaced:'%s',expectParamValues:%+v,paramValues:%+v",
					routePathParam.route, pathParam.path, pathParam.replaced, replaced, pathParam.paramValues, paramValues)
			}
		}
	}
}

func Benchmark_RadixTreeRouter_replaceRouteParams(b *testing.B) {
	rtr := NewRadixTreeRouter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, routePathParam := range routePathParams {
			for _, pathParam := range routePathParam.paths {
				_, _ = rtr.replaceRouteParams(routePathParam.route, pathParam.path)
			}
		}
	}
}

func Benchmark_replaceRouteParams(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, routePathParam := range routePathParams {
			for _, pathParam := range routePathParam.paths {
				_, _ = replaceRouteParams(routePathParam.route, pathParam.path)
			}
		}
	}
}
