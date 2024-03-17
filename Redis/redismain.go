package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "thaer", // no password set
		Password: "thaer", // no password set
		DB:       0,
	})

	val, err := rdb.Get(ctx, "user:2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("user:2", val)

	err2 := rdb.Set(ctx, "user:3", "Thaer", 0).Err()
	if err != nil {
		panic(err2)
	}

	val, err3 := rdb.Get(ctx, "user:3").Result()
	if err != nil {
		panic(err3)
	}
	fmt.Println("user:3", val)

}
