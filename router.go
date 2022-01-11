package fastweb

import (
	"errors"
	"regexp"
	"strings"
)

type Router interface {
	RegisterRoute(path, method string, h HttpHandler)
	Route(path, method string) (*route, error)
	Params(ctx *HttpContext, regPath, path string) error
}

type route struct {
	path    string // 注册的原始路径
	method  string
	handler HttpHandler
}

type defaultRouter struct {
	routes map[int]*route
	index  []string
}

func DefaultRouter() *defaultRouter {
	return &defaultRouter{
		routes: make(map[int]*route),
		index:  make([]string, 0),
	}
}

func (r *defaultRouter) parseRoute(path string) string {
	path = strings.TrimSuffix(path, "/")

	paths := strings.Split(path, "/")
	for i, p := range paths {
		if strings.Contains(p, ":") {
			paths[i] = `\w+`
		}
	}

	//return fmt.Sprintf(`^%s\/*$`, strings.Join(paths, "/"))
	return `^` + strings.Join(paths, "/") + `\/*$`
}

func (r *defaultRouter) RegisterRoute(path, method string, h HttpHandler) {
	pr := r.parseRoute(path)

	// map无序，先把路由保存在有序的slice中，保证顺序
	r.index = append(r.index, pr)

	r.routes[len(r.index)-1] = &route{
		path:    path,
		method:  method,
		handler: h,
	}
}

/*
func (r *router) route(ctx *HttpContext, path, method string) *route {
	// 逆序，后注册的路由会覆盖之前注册的同名路由
	l := len(r.index)
	for l > 0 {
		routePath := r.index[l-1]

		// 正则表达式匹配
		if matched, err := regexp.MatchString(routePath, path); err == nil && matched {
			return r.routes[routePath]
		}

		l--
	}

	return nil
}
*/

func (r *defaultRouter) Route(path, method string) (*route, error) {
	s := make([]*route, 0, len(r.index))

	// 逆序，后注册的路由会覆盖之前注册的同名路由
	l := len(r.index)
	for l > 0 {
		i := l - 1
		routePath := r.index[i]

		// 正则表达式匹配
		if matched, err := regexp.MatchString(routePath, path); err == nil && matched {
			route := r.routes[i]
			if route == nil {
				continue
			}

			if route.method == method {
				return route, nil
			}

			s = append(s, route)
		}

		l--
	}

	// 路由匹配，但是http method不匹配
	if len(s) > 0 {
		return nil, ErrMethodNotAllowed
	}

	return nil, ErrNotFound
}

func (r *defaultRouter) Params(ctx *HttpContext, regPath, path string) error {
	path = strings.TrimSuffix(path, "/")
	regPath = strings.TrimSuffix(regPath, "/")

	regPaths := strings.Split(regPath, "/")
	paths := strings.Split(path, "/")

	if len(regPaths) != len(paths) {
		return errors.New("route bind error")
	}

	for i, rp := range regPaths {
		if strings.HasPrefix(rp, ":") {
			k := strings.TrimPrefix(rp, ":")
			v := paths[i]

			ctx.Params[k] = v
		}
	}

	return nil
}
