package fastweb

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type RadixTreeRouter struct {
	rootRoute       *radixTreeRoute
	paramValuesPool sync.Pool
}

type radixTreeRoute struct {
	Prefix     string
	Method     map[string]HttpHandler
	Params     []string
	Children   childRoutes
	IsRoute    bool
	IsParam    bool
	OriginPath string
}

type childRoutes []*radixTreeRoute

func (routes childRoutes) Len() int {
	return len(routes)
}

func (routes childRoutes) Less(i, j int) bool {
	//return routes[j].Prefix == ":"
	return strings.HasPrefix(routes[j].Prefix, ":")
}

func (routes childRoutes) Swap(i, j int) {
	routes[i], routes[j] = routes[j], routes[i]
}

func NewRadixTreeRouter() *RadixTreeRouter {
	return &RadixTreeRouter{
		rootRoute: &radixTreeRoute{
			Children: make([]*radixTreeRoute, 0),
		},
		paramValuesPool: sync.Pool{
			New: func() any {
				return make([]string, 0, 10)
			},
		},
	}
}

func (rtr *RadixTreeRouter) RegisterRoute(path, method string, handler HttpHandler) {
	originPath := path

	if len(path) == 0 {
		panic("path cannot be empty")
	}

	path, isParam, params, err := parsePathParams(path)
	if err != nil {
		panic(err)
	}

	if len(rtr.rootRoute.Children) == 0 {
		rtr.rootRoute.Children = append(rtr.rootRoute.Children, &radixTreeRoute{
			Prefix:     path,
			Method:     map[string]HttpHandler{method: handler},
			IsParam:    isParam,
			Params:     params,
			IsRoute:    true,
			Children:   make([]*radixTreeRoute, 0),
			OriginPath: originPath,
		})
		return
	}

	parent := rtr.rootRoute

forChild:
	for j := 0; j < len(parent.Children); j++ {
		child := parent.Children[j]

		lenPath := len(path)
		lenPrefix := len(child.Prefix)
		max := lenPath
		if lenPrefix < max {
			max = lenPrefix
		}

		i := 0
		for i < max && child.Prefix[i] == path[i] {
			i++
		}

		// 完全相等，已经存在
		if i == lenPrefix && i == lenPath {
			if _, exists := child.Method[method]; exists {
				panic("exist similar path:" + child.OriginPath)
			}

			child.Method[method] = handler
			child.IsRoute = true
			child.IsParam = isParam
			child.Params = params
			child.OriginPath = originPath

			break
		}

		// 完全不相等
		if i <= 0 {
			// 同级最后一个节点不相等，附加到当前父节点的子节点集合
			if j == len(parent.Children)-1 {
				parent.Children = append(parent.Children, &radixTreeRoute{
					Prefix:     path,
					Method:     map[string]HttpHandler{method: handler},
					IsParam:    isParam,
					Params:     params,
					IsRoute:    true,
					Children:   make([]*radixTreeRoute, 0),
					OriginPath: originPath,
				})

				sort.Sort(parent.Children)
				break
			}

			// 继续下一个同级的节点对比
			continue
		}

		// /user
		// /user/:id
		if i == lenPrefix && i < lenPath {
			if len(child.Children) > 0 {
				parent = child
				path = path[i:]

				goto forChild
			} else {
				child.Children = append(child.Children, &radixTreeRoute{
					Prefix:     path[i:],
					Method:     map[string]HttpHandler{method: handler},
					IsParam:    isParam,
					Params:     params,
					IsRoute:    true,
					Children:   make([]*radixTreeRoute, 0),
					OriginPath: originPath,
				})

				sort.Sort(child.Children)
				break
			}
		}

		// /user/:id
		// /user
		if i < lenPrefix && i == lenPath {
			newChild := &radixTreeRoute{
				Prefix:     child.Prefix[i:],
				Method:     child.Method,
				IsRoute:    child.IsRoute,
				IsParam:    child.IsParam,
				Params:     child.Params,
				Children:   child.Children,
				OriginPath: child.OriginPath,
			}
			child.Children = []*radixTreeRoute{
				newChild,
			}

			child.Prefix = path
			child.Method = map[string]HttpHandler{method: handler}
			child.IsRoute = true
			child.IsParam = isParam
			child.Params = params
			child.OriginPath = originPath

			sort.Sort(parent.Children)
			break
		}

		// /user/:id/edit
		// /user/:id/delete
		if i < lenPrefix && i < lenPath {
			newChild := child.Children

			child.Children = []*radixTreeRoute{
				{
					Prefix:     child.Prefix[i:],
					Method:     child.Method,
					IsRoute:    child.IsRoute,
					IsParam:    child.IsParam,
					Params:     child.Params,
					Children:   newChild,
					OriginPath: child.OriginPath,
				}, {
					Prefix:     path[i:],
					Method:     map[string]HttpHandler{method: handler},
					IsRoute:    true,
					IsParam:    isParam,
					Params:     params,
					Children:   make([]*radixTreeRoute, 0),
					OriginPath: originPath,
				},
			}

			sort.Sort(child.Children)

			child.Prefix = path[:i]
			child.Method = make(map[string]HttpHandler)
			child.IsRoute = false
			child.IsParam = false
			child.Params = make([]string, 0)
			child.OriginPath = ""

			sort.Sort(parent.Children)
			break
		}
	}
}

func (rtr *RadixTreeRouter) Route(ctx *HttpContext, path, method string) (*Route, error) {
	paramValues := make([]string, 0)

	parent := rtr.rootRoute

forChild:
	for j := 0; j < len(parent.Children); j++ {
		child := parent.Children[j]
		prefix := child.Prefix

		// 如果是参数路由，先把路由中的参数替换成值
		if child.IsParam {
			p, pv := rtr.replaceRouteParams(prefix, path)
			prefix = p
			paramValues = append(paramValues, pv...)
		}

		lenPath := len(path)
		lenPrefix := len(prefix)
		max := lenPath
		if lenPrefix < max {
			max = lenPrefix
		}

		i := 0
		for i < max && prefix[i] == path[i] {
			i++
		}

		// 完全相等
		if i == lenPrefix && i == lenPath && child.IsRoute {
			if handler, exists := child.Method[method]; exists {
				for k, param := range child.Params {
					ctx.Params[param[1:]] = paramValues[k]
				}

				return &Route{
					Path:    child.OriginPath,
					Method:  method,
					Handler: handler,
				}, nil
			}

			return nil, ErrMethodNotAllowed
		}

		// 完全不相等
		if i == 0 {
			// 继续同级的下一个节点匹配
			continue
		}

		// /user
		// /user/add
		if i == lenPrefix && i < lenPath {
			if len(child.Children) > 0 {
				parent = child
				path = path[i:]

				goto forChild
			} else {
				return nil, nil
			}
		}

		// /user/add
		// /user
		if i < lenPrefix && i == lenPath {
			return nil, nil
		}

		// /user/add
		// /user/edit
		if i < lenPrefix && i < lenPath {
			return nil, nil
		}
	}
	return nil, nil
}

func (rtr *RadixTreeRouter) Params(ctx *HttpContext, regPath, path string) error {
	// 已经在路由过程中绑定了参数，直接返回
	return nil
}

func parsePathParams(path string) (string, bool, []string, error) {
	r := regexp.MustCompile(`:\w*`)

	parsedPath := path

	indexes := r.FindAllStringIndex(path, -1)
	params := make([]string, len(indexes))
	isParam := len(indexes) > 0

	parsedLen := 0
	for i, index := range indexes {
		start := index[0]
		stop := index[1]
		if start > 0 && path[start-1] != '/' {
			return "", false, nil, fmt.Errorf("path contains invalid parameter string:'%s'", params[i])
		}

		params[i] = path[start:stop]
		parsedPath = parsedPath[:start-parsedLen+1] + parsedPath[stop-parsedLen:]

		parsedLen = parsedLen + stop - start - 1
	}

	return parsedPath, isParam, params, nil
}

func (rtr *RadixTreeRouter) replaceRouteParams(prefix, path string) (string, []string) {
	paramValues := rtr.paramValuesPool.Get().([]string)
	defer rtr.paramValuesPool.Put(paramValues[0:0])

	var start int
	for {
		start = strings.IndexByte(prefix, ':')
		// prefix不包含参数则退出
		if start < 0 || start >= len(path) {
			break
		}

		stop := start
		if start < len(path) {
			stop = strings.IndexByte(path[start:], '/')
			if stop < 0 {
				stop = len(path)
			} else {
				stop += start
			}
		}

		paramValue := path[start:stop]

		prefix = prefix[:start] + paramValue + prefix[start+1:]

		paramValues = append(paramValues, paramValue)
	}

	return prefix, paramValues
}

var paramValuesPool = sync.Pool{
	New: func() any {
		return make([]string, 0, 10)
	},
}

var bytesPool = sync.Pool{
	New: func() any {
		return make([]byte, 0, 100)
	},
}

func replaceRouteParams(prefix, path string) (string, []string) {
	paramValues := paramValuesPool.Get().([]string)
	defer paramValuesPool.Put(paramValues[0:0])

	newBytes := bytesPool.Get().([]byte)
	defer bytesPool.Put(newBytes[0:0])

	var i, prefixStart, pathStart, pathStop, copied int
	for {
		if len(prefix)+copied-i <= pathStop {
			break
		}

		if i > 0 {
			prefixStart++
		}
		index := strings.IndexByte(prefix[prefixStart:], ':')
		if index < 0 {
			newBytes = append(newBytes, prefix[pathStop-copied+i:]...)
			break
		}

		prefixStart += index
		pathStart = prefixStart + copied - i

		if pathStart < len(path) {
			if i > 0 {
				newBytes = append(newBytes, prefix[pathStop-copied+i:prefixStart]...)
			}

			pathStop = strings.IndexByte(path[pathStart:], '/')
			if pathStop < 0 {
				pathStop = len(path)
			} else {
				pathStop += pathStart
			}
		} else {
			newBytes = append(newBytes, prefix[pathStop-copied+i:]...)
			break
		}

		if i == 0 && i < pathStart {
			newBytes = append(newBytes, prefix[0:pathStart]...)
		}

		newBytes = append(newBytes, path[pathStart:pathStop]...)
		copied = copied + pathStop - pathStart

		paramValues = append(paramValues, path[pathStart:pathStop])
		i++
	}

	return string(newBytes), paramValues
}
