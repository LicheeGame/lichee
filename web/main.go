package main

import (
	"fmt"

	"web/auth"
	"web/cache"
	"web/config"
	"web/dao"
	"web/logger"
	"web/router"
)

func main() {
	logger.Init(&config.Conf.Log)
	auth.InitJwt()
	cache.InitRedis()
	dao.InitDB()
	logger.Info("start router")

	r := router.Router()
	if err := r.Run(fmt.Sprintf(":%d", config.Conf.Port)); err != nil {
		panic(err)
	}
}
