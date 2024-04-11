package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	user struct {
		ID int `json:"id"`
		Name string `json:"name"`
	}
)

var (
	// DBの代わり。users[0]とかで、*user（userへのポインタ）を得る
	users = map[int]*user{}
	// userのIDに指定する値。userが増えるたびにインクリメントする
	seq = 1
	// 排他制御をやってくれる
	mutex = sync.Mutex{}
)

func getAllUsers(c echo.Context) error {
	mutex.Lock()
	defer mutex.Unlock() // 忘れないようにdefer
	return c.JSON(http.StatusOK, users)
}

func createUser(c echo.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	// IDだけは先にセット
	u := &user{
		ID: seq,
	}
	// ID以外は、リクエストボディを構造体にバインド
	if err := c.Bind(u); err != nil {
		return err
	}
	// モックDBへ保存
	users[u.ID] = u

	seq++ // 次のCreateに備えてIDをインクリメント

	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	id, _ := strconv.Atoi(c.Param("id"))
	// 存在していなくても200を返す（という仕様でいい場合）
	return c.JSON(http.StatusOK, users[id])
}

func updateUser(c echo.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}

func deleteUser(c echo.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id) // mapから削除
	return c.NoContent(http.StatusNoContent)
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/users", getAllUsers)
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// Start server
	e.Logger.Fatal(e.Start(":5050"))
}