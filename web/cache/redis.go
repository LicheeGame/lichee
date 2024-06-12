package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ctx context.Context
	Rdb *redis.Client
)

func InitRedis() {
	ctx = context.Background()
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	/*
		err := rdb.Set(ctx, "key", "value", 0).Err()
		if err != nil {
			panic(err)
		}

		val, err := rdb.Get(ctx, "key").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("key", val)
	*/
}

// 设置分数
func SetUserScore(appid string, uid string, score int) {
	Rdb.ZAdd(ctx, fmt.Sprintf("%s_rank", appid), redis.Z{Score: float64(score), Member: uid}).Err()
}

func GetUserScoreRank(appid string) ([]redis.Z, error) {
	// 获取排行榜上的成员前20
	scoreRank, err := Rdb.ZRevRangeWithScores(ctx, fmt.Sprintf("%s_rank", appid), 0, 19).Result()
	return scoreRank, err
}
