package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Host struct {
	Echo *echo.Echo
}

func main() {
	hosts := map[string]*Host{}

	// ----------
	// サブドメイン api
	// ----------
	api := echo.New()
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())

	hosts["api.localhost:5050"] = &Host{api}

	api.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "API")
	})

	// ----------
	// ウェブサイト
	// ----------
	site := echo.New()
	site.Use(middleware.Logger())
	site.Use(middleware.Recover())

	hosts["localhost:5050"] = &Host{site}

	site.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Website")
	})

	// ----------
	// 本体
	// ----------
	e := echo.New()
	// Anyに渡したfuncはリクエストのたびに実行される
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})

	e.Logger.Fatal(e.Start(":5050"))
}