package fastweb

type Middleware func(HttpHandler) HttpHandler

func combineMiddleware(handler HttpHandler, middlewares ...Middleware) HttpHandler {
	l := len(middlewares)
	for l > 0 {
		handler = middlewares[l-1](handler)

		l--
	}

	return handler
}
