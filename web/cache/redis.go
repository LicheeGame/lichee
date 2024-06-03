package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ctx context.Context
	rdb *redis.Client
)

func InitRedis() {
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 100, Member: "player1"}).Err()

	// 获取排行榜上的成员
	leaderboard, err := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	for _, scoreMember := range leaderboard {
		fmt.Printf("Member: %s, Score: %f\n", scoreMember.Member, scoreMember.Score)
	}

}
