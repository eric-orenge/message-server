package main

import (
	"os"
	"strconv"

	"github.com/eric-orenge/message-server/api"
	"github.com/eric-orenge/message-server/controller"
	"github.com/eric-orenge/message-server/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type App struct {
	RedisCtrl     controller.RedisCtrl
	WebsocketCtrl controller.WebsocketCtrl
	Server        *echo.Echo
}

func (app *App) Start() {
	app.Server = echo.New()

	app.RedisCtrl.Address = os.Getenv("REDIS_URL")
	app.RedisCtrl.Password = os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)
	app.RedisCtrl.DB = db

	err := app.RedisCtrl.Connect()
	if err != nil {
		panic(err)
	}

	pool := websocket.NewPool()
	go pool.Start()

	app.Server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cc := &api.CustomContext{ctx, nil}
			cc.SetRedis(app.RedisCtrl.Client)
			return next(cc)
		}
	})
	app.Server.Use(middleware.CORS())
	app.Server.Use(middleware.Logger())
	app.Server.Use(middleware.Recover())
	app.Server.GET("/ws/:roomid", func(ctx echo.Context) error {
		return app.WebsocketCtrl.WsEndpoint(pool, ctx)
	})
	app.Server.GET("/roomid", api.GetRoomID)
	port := ":" + os.Getenv("PORT")
	app.Server.Logger.Fatal(app.Server.Start(port))
}
