package main

import (
	"os"

	"github.com/eric-orenge/message-server/api"
	"github.com/eric-orenge/message-server/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type App struct {
	Server *echo.Echo
}

func (app *App) Start() {
	app.Server = echo.New()

	pool := websocket.NewPool()
	go pool.Start()

	app.Server.Use(middleware.CORS())
	app.Server.Use(middleware.Logger())
	app.Server.Use(middleware.Recover())
	app.Server.GET("/ws/:roomid", func(ctx echo.Context) error {
		return api.WsEndpoint(pool, ctx)
	})
	app.Server.GET("/roomid", api.GetRoomID)
	port := ":" + os.Getenv("PORT")
	app.Server.Logger.Fatal(app.Server.Start(port))
}
