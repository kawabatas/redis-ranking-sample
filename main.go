package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if err := rdb.FlushAll(ctx).Err(); err != nil {
		log.Fatalf("FlushAll error: %v", err)
	}

	// ソート済みセット
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
		log.Fatalf("ZRevRangeWithScores error: %v", err)
	}
	for i, m := range top10 {
		log.Printf("rank %d, user %s, score %d\n", i+1, m.Member, int(m.Score))
	}
	around50, err := rdb.ZRevRangeWithScores(ctx, key1, int64(45), int64(54)).Result()
	if err != nil {
		log.Fatalf("ZRevRangeWithScores error: %v", err)
	}
	for i, m := range around50 {
		log.Printf("rank %d, user %s, score %d\n", i+45, m.Member, int(m.Score))
	}

	// 文字列
	if err := rdb.Set(ctx, "total", 0, 0).Err(); err != nil {
		log.Fatalf("Set error: %v", err)
	}
	total, err := rdb.Get(ctx, "total").Result()
	if err != nil {
		log.Fatalf("Get error: %v", err)
	}
	log.Printf("total %s\n", total)
	if err := rdb.IncrBy(ctx, "total", 100).Err(); err != nil {
		log.Fatalf("IncrBy error: %v", err)
	}
	total2, err := rdb.Get(ctx, "total").Result()
	if err != nil {
		log.Fatalf("Get error: %v", err)
	}
	log.Printf("total %s\n", total2)

	// 文字列 Expire
	if err := rdb.SetEX(ctx, "limit", 0, time.Minute).Err(); err != nil {
		log.Fatalf("SetEX error: %v", err)
	}
	limit, err := rdb.Get(ctx, "limit").Result()
	if err != nil {
		log.Fatalf("Get error: %v", err)
	}
	log.Printf("limit %s\n", limit)
	d, err := rdb.TTL(ctx, "limit").Result()
	if err != nil {
		log.Fatalf("TTL error: %v", err)
	}
	log.Printf("limit duration %v\n", d)
	time.Sleep(3 * time.Second)
	if err := rdb.IncrBy(ctx, "limit", 1).Err(); err != nil {
		log.Fatalf("IncrBy error: %v", err)
	}
	limit2, err := rdb.Get(ctx, "limit").Result()
	if err != nil {
		log.Fatalf("Get error: %v", err)
	}
	log.Printf("limit %s\n", limit2)
	d2, err := rdb.TTL(ctx, "limit").Result()
	if err != nil {
		log.Fatalf("TTL error: %v", err)
	}
	log.Printf("limit duration %v\n", d2)

	// リスト
	lkey := "list1"
	for i := 1; i <= 500; i++ {
		r := time.Now()
		if err := rdb.LPush(ctx, lkey, r).Err(); err != nil {
			log.Fatalf("LPush error: %v", err)
		}
	}
	llen, err := rdb.LLen(ctx, lkey).Result()
	if err != nil {
		log.Fatalf("LLen error: %v", err)
	}
	log.Printf("list count %d\n", int(llen))
	ltop10, err := rdb.LRange(ctx, lkey, int64(0), int64(9)).Result()
	if err != nil {
		log.Fatalf("LRange error: %v", err)
	}
	for i, v := range ltop10 {
		log.Printf("index %d, value %v\n", i, v)
	}
	laround50, err := rdb.LRange(ctx, lkey, int64(45), int64(54)).Result()
	if err != nil {
		log.Fatalf("LRange error: %v", err)
	}
	for i, v := range laround50 {
		log.Printf("index %d, value %v\n", i+44, v)
	}
}
