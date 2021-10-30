package main

import (
	"context"
	"log"
	"math/rand"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	key1 := "rank1"
	for i := 1; i <= 1000; i++ {
		score := rand.Intn(1000) * 100

		if err := rdb.ZAdd(ctx, key1, &redis.Z{
			Score:  float64(score),
			Member: i,
		}).Err(); err != nil {
			log.Fatalf("ZAdd error: %v", err)
		}
	}

	count, err := rdb.ZCard(ctx, key1).Result()
	if err != nil {
		log.Fatalf("ZScore error: %v", err)
	}
	log.Printf("all members count %d\n", int(count))
	mem500score, err := rdb.ZScore(ctx, key1, "500").Result()
	if err != nil {
		log.Fatalf("ZScore error: %v", err)
	}
	mem500rank, err := rdb.ZRevRank(ctx, key1, "500").Result()
	if err != nil {
		log.Fatalf("ZScore error: %v", err)
	}
	log.Printf("rank %d, user 500, score %d\n", mem500rank+1, int(mem500score))
	if err := rdb.ZIncrBy(ctx, key1, 100, "500").Err(); err != nil {
		log.Fatalf("ZIncrBy error: %v", err)
	}
	mem500score2, err := rdb.ZScore(ctx, key1, "500").Result()
	if err != nil {
		log.Fatalf("ZScore error: %v", err)
	}
	log.Printf("rank -, user 500, score %d\n", int(mem500score2))

	top10, err := rdb.ZRevRangeWithScores(ctx, key1, int64(0), int64(9)).Result()
	if err != nil {
		log.Fatalf("Redis ZRevRangeWithScores error: %v", err)
	}
	for i, m := range top10 {
		log.Printf("rank %d, user %s, score %d\n", i+1, m.Member, int(m.Score))
	}
	around50, err := rdb.ZRevRangeWithScores(ctx, key1, int64(45), int64(54)).Result()
	if err != nil {
		log.Fatalf("Redis ZRevRangeWithScores error: %v", err)
	}
	for i, m := range around50 {
		log.Printf("rank %d, user %s, score %d\n", i+45, m.Member, int(m.Score))
	}
}
