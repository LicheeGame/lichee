package main

import (
	"fmt"

	"web/config"
	"web/dao"
	"web/logger"
	"web/router"
)

func main() {
	logger.Init(&config.Conf.Log)
	dao.InitDB()
	logger.Info("start router")

	r := router.Router()
	if err := r.Run(fmt.Sprintf(":%d", config.Conf.Port)); err != nil {
		panic(err)
	}
	//r.RunTLS(":8080", "./testdata/server.pem", "./testdata/server.key")
}
