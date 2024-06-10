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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(moongoConf.Dns).SetMaxPoolSize(100))
	if err != nil {
		logger.Error("mongodb connect error :%s", err)
		return
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Error("mongodb ping error :%s", err)
		return
	}
	DB = client
}

// 根据appid获取db
func GetDB(appid string) *mongo.Database {
	return DB.Database(config.GetMongoDBByAppID(appid))
}
