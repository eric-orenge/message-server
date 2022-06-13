package api

import (
	"github.com/labstack/echo"
	"gopkg.in/redis.v3"
)

type CustomContext struct {
	echo.Context
	RedisClient *redis.Client
}

func (c *CustomContext) SetRedis(client *redis.Client) {
	c.RedisClient = client
}
