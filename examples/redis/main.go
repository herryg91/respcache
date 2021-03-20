package main

import (
	"fmt"
	"log"

	rdsInit "github.com/herryg91/homegym-be/libs/db/redis"
	"github.com/herryg91/respcache"
)

func main() {
	rdsPool, err := rdsInit.Connect("localhost", 6379, "")
	if err != nil {
		panic(fmt.Sprintf("Failed to Initialized redis: %v", err))
	}

	rc := respcache.NewRedisCache(rdsPool)
	var result []Author
	iscached, err := rc.Run("test_redis_simple_cache", 10, &result, func() (interface{}, error) {
		return GetAuthors(), nil
	})

	log.Println(fmt.Sprintf("[cached:%v] %v", iscached, result))
}

type Author struct {
	Id   int
	Name string
}

func GetAuthors() []Author {
	return []Author{
		Author{Id: 1, Name: "Name 1"},
		Author{Id: 2, Name: "Name 2"},
		Author{Id: 3, Name: "Name 3"},
	}
}
