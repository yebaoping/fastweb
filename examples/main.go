package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/yebaoping/fastweb"

	"github.com/yebaoping/fastweb/middleware"
)

func main() {
	logger := fastweb.DefaultLogger(os.Stdout, true)
	app := fastweb.NewApp(logger, fastweb.NewRadixTreeRouter(), nil)

	app.ErrorHandler = func(ctx *fastweb.HttpContext, err error) {
		if err != nil {
			logger.Errorf("ErrorHandler:%+v", err)

			_ = ctx.Error(err)
		}
	}

	app.Use(middleware.Recover, func(next fastweb.HttpHandler) fastweb.HttpHandler {
		return func(ctx *fastweb.HttpContext) error {
			logger.Debug("app use middleware exec 1")
			return next(ctx)
		}
	}, middleware.CORS(middleware.DefaultCORSConfig))

	app.Static("/", "./static/index.html", fastweb.ContentTypeTextHtml, nil)

	app.Head("/", func(ctx *fastweb.HttpContext) error {
		return nil
	})

	userGroup := app.Group("/user",
		func(next fastweb.HttpHandler) fastweb.HttpHandler {
			return func(ctx *fastweb.HttpContext) error {
				logger.Debug("app useTime middleware exec")

				start := time.Now()

				err := next(ctx)
				if err != nil {
					return err
				}

				logger.Infof("url:%s,method:%s,user-agent:%s,receiveTime:%s,useTime:%fms",
					ctx.RequestURI(), ctx.Method(), ctx.UserAgent(), start.Format(time.RFC3339), time.Since(start).Seconds()*1000)

				return nil
			}
		},
	)

	userGroup.Get("", func(ctx *fastweb.HttpContext) error {
		return ctx.Json(users)
	})

	userGroup.Get("/", func(ctx *fastweb.HttpContext) error {
		return ctx.Json(users)
	})

	userGroup.Get("/:start/:stop", func(ctx *fastweb.HttpContext) error {
		pstart := ctx.Params.Get("start")
		pstop := ctx.Params.Get("stop")

		start, err := strconv.ParseInt(pstart, 10, 64)
		if err != nil {
			return err
		}

		stop, err := strconv.ParseInt(pstop, 10, 64)
		if err != nil {
			return err
		}

		s := make([]*User, 0)
		for _, u := range users {
			s = append(s, u)
		}

		return ctx.Json(s[start : stop+1])
	})

	userGroup.Get("/:id", func(ctx *fastweb.HttpContext) error {
		pId := ctx.Params.Get("id")
		id, err := strconv.ParseInt(pId, 10, 64)
		if err != nil {
			return err
		}

		return ctx.Json(users[id])
	}, func(next fastweb.HttpHandler) fastweb.HttpHandler {
		return func(ctx *fastweb.HttpContext) error {
			logger.Debug("app use middleware exec 2")
			return next(ctx)
		}
	})

	userGroup.Post("/:id", func(ctx *fastweb.HttpContext) error {
		pId := ctx.Params.Get("id")
		id, err := strconv.ParseInt(pId, 10, 64)
		if err != nil {
			return err
		}

		u := &User{}
		if err := ctx.Bind(u); err != nil {
			return err
		}

		users[id] = &User{
			ID:    id,
			Name:  u.Name,
			Addrs: u.Addrs,
			Roles: u.Roles,
		}

		return ctx.Json(users[id])
	})

	userGroup.Get("/:start/:stop/export", func(ctx *fastweb.HttpContext) error {
		pstart := ctx.Params.Get("start")
		pstop := ctx.Params.Get("stop")

		start, err := strconv.ParseInt(pstart, 10, 64)
		if err != nil {
			return err
		}

		stop, err := strconv.ParseInt(pstop, 10, 64)
		if err != nil {
			return err
		}

		s := make([]*User, 0)
		for _, u := range users {
			s = append(s, u)
		}

		fileName := "export_" + time.Now().Format("20060102150405") + ".csv"
		filePath := path.Join(os.TempDir(), fileName)
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return ctx.Error(err)
		}
		defer func() {
			file.Close()
			os.Remove(filePath)
		}()

		file.Write([]byte{0xEF, 0xBB, 0xBF})
		file.WriteString("id, name, address\n")
		for i := start; i <= stop; i++ {
			u := s[i]
			file.WriteString(fmt.Sprintf("%d, %s, %v\n", u.ID, u.Name, u.Addrs))
		}
		file.Sync()

		// ctx.SendFile(filePath)
		return ctx.Attachment(filePath, fileName)
	})

	userGroup.Get("/:id/delete", func(ctx *fastweb.HttpContext) error {
		pId := ctx.Params.Get("id")
		id, err := strconv.ParseInt(pId, 10, 64)
		if err != nil {
			return err
		}

		delete(users, id)

		return ctx.Json(users)
	})

	userGroup.Get("/panic", func(ctx *fastweb.HttpContext) error {
		panic("throw panic.")
	})

	app.Get("/panic", func(ctx *fastweb.HttpContext) error {
		panic("throw app panic")
	})

	app.Use(func(next fastweb.HttpHandler) fastweb.HttpHandler {
		return func(ctx *fastweb.HttpContext) error {
			logger.Debug("app use middleware exec 3")

			return next(ctx)
		}
	})

	roleGroup := app.Group("/role", func(next fastweb.HttpHandler) fastweb.HttpHandler {
		return func(ctx *fastweb.HttpContext) error {
			logger.Debug("app use middleware exec 4")

			return next(ctx)
		}
	})

	roleGroup.Get("", func(ctx *fastweb.HttpContext) error {
		logger.Debug("exec get role list.")
		return ctx.Json(roles)
	}, func(next fastweb.HttpHandler) fastweb.HttpHandler {
		return func(ctx *fastweb.HttpContext) error {
			logger.Debug("app use middleware exec 5")
			return next(ctx)
		}
	})

	roleGroup.Get("/:id", func(ctx *fastweb.HttpContext) error {
		pId := ctx.Params.Get("id")
		id, err := strconv.ParseInt(pId, 10, 64)
		if err != nil {
			return err
		}

		return ctx.Json(roles[id])
	})

	roleGroup.Html("/about", "./static/about.html", roles, func(next fastweb.HttpHandler) fastweb.HttpHandler {
		return func(ctx *fastweb.HttpContext) error {
			logger.Debug("app use middleware exec 6")
			return next(ctx)
		}
	}, func(next fastweb.HttpHandler) fastweb.HttpHandler {
		return func(ctx *fastweb.HttpContext) error {
			logger.Debug("app use middleware exec 7")
			return next(ctx)
		}
	})

	logger.Fatal(app.Start(":8000"))
}

type User struct {
	ID    int64    `json:"id"`
	Name  string   `json:"name"`
	Addrs []string `json:"address"`
	Roles []*Role  `json:"roles"`
}

var users = map[int64]*User{
	1: {
		ID:    1,
		Name:  "tom",
		Addrs: []string{"hangzhou"},
		Roles: []*Role{
			roles[1],
			roles[2],
		},
	},
	2: {
		ID:    2,
		Name:  "jack",
		Addrs: []string{"beijing"},
	},
	3: {
		ID:    3,
		Name:  "jim",
		Addrs: []string{"shanghai"},
	},
}

type Role struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Level       int64        `json:"level"`
	Score       float64      `json:"score"`
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	Resource string   `json:"resource"`
	Verb     []string `json:"verb"`
}

var roles = map[int64]*Role{
	1: {
		ID:    1,
		Name:  "admin",
		Level: 1,
		Score: 100,
	},
	2: {
		ID:    2,
		Name:  "owner",
		Level: 2,
		Score: 90,
	},
	3: {
		ID:    3,
		Name:  "guest",
		Level: 3,
		Score: -1,
	},
}
