package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func GetRoomID(ctx echo.Context) error {
	cc := ctx.(*CustomContext)
	result, err := cc.RedisClient.Incr("roomCount").Result()
	if err != nil {
		panic(err)
	}
	return ctx.JSON(http.StatusOK, result)
}
