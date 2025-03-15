package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

func main() {
	// redis connection pool을 생성
	redisAddrs := []string{
		"localhost:6379",
		"localhost:16379",
		"localhost:26379",
		"localhost:36379",
		"localhost:46379",
	}
	var pools []redsyncredis.Pool
	for _, addr := range redisAddrs {
		client := redis.NewClient(&redis.Options{
			Addr: addr,
		})

		// ping
		if err := client.Ping(context.Background()).Err(); err != nil {
			log.Fatal(err)
		}

		pools = append(pools, goredis.NewPool(client))
	}

	// redis connection pool을 이용하여 redsync 인스턴스를 생성
	rs := redsync.New(pools...)

	mutexname := "my-global-mutex"

	// 주어진 mutexname을 이용하여 Mutex 인스턴스를 생성
	mutex := rs.NewMutex(mutexname)

	// Lock을 획득하여 다른 프로세스나 스레드가 Lock을 획득할 수 없도록 함
	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	// 작업 수행
	{
		// do something
		time.Sleep(1 * time.Second)
	}

	// Lock을 해제하여 다른 프로세스나 스레드가 Lock을 획득할 수 있도록 함
	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
}
