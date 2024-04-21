package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

func hello(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		// ws接続時のみ実行される
		defer ws.Close()

		// wsが接続している間、ずっと実行される
		for {
			// wsへ送信
			err := websocket.Message.Send(ws, "Hello!!")
			if err != nil {
				c.Logger().Error(err)
			}

			// 受診
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")

	e.GET("/ws", hello)

	// Start server
	e.Logger.Fatal(e.Start(":5050"))
}