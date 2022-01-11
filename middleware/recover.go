package middleware

import (
	"fmt"

	"github.com/yebaoping/fastweb"
)

var Recover = func(next fastweb.HttpHandler) fastweb.HttpHandler {
	return func(ctx *fastweb.HttpContext) (err error) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			if e, ok := r.(error); ok && e != nil {
				err = e
			} else {
				err = fmt.Errorf("panic:%+v", r)
			}
		}()

		err = next(ctx)
		return
	}
}
