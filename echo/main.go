package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// DBを用意していないので、簡易で
	if username != "jon" || password != "shhh!" {
		return echo.ErrUnauthorized
	}

	claims := jwtCustomClaims{
		"jon Snow",
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// トークン生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 署名済みトークンへ
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println("ほげ")
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login
	e.POST("/login", login)

	// 認証していなくても利用できるリソース
	e.GET("/", accessible)

	// 認証していないと使えないリソースは/restricted/*
	r := e.Group("/restricted")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("", restricted)

	// Start server
	e.Logger.Fatal(e.Start(":5050"))
}
