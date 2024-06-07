package dao

import (
	"context"
	"time"
	"web/config"
	"web/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB *mongo.Client
)

func InitDB() {
	moongoConf := config.Conf.Mongodb

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	DB, err := mongo.Connect(ctx, options.Client().ApplyURI(moongoConf.Dns))
	if err != nil {
		logger.Error("mongodb connect error :%s", err)
		return
	}

	err = DB.Ping(context.TODO(), nil)
	if err != nil {
		return
	}
}
