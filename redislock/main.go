package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
}

func newRedisClient(addr, password string, db int) *redisClient {
	return &redisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,     // Redis 주소
			Password: password, // 패스워드 (없으면 빈 문자열)
			DB:       db,       // 기본 DB
		}),
	}
}

func (r *redisClient) acquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	ok, err := r.client.SetNX(ctx, key, "locked", ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *redisClient) releaseLock(ctx context.Context, key string) {
	r.client.Del(ctx, key)
}

// 락 관련 상수
const lockKey = "cron:job:my-task"
const lockTTL = 30 * time.Second // 락 유지 시간

func main() {
	ctx := context.Background()

	// Redis 클라이언트 생성
	client := newRedisClient("localhost:6379", "", 0)

	// 분산 락 획득 시도
	acquired, err := client.acquireLock(ctx, lockKey, lockTTL)
	if err != nil {
		log.Fatalf("Redis 연결 오류: %v", err)
	}

	if !acquired {
		fmt.Println("다른 서버에서 이미 크론 작업 실행 중! 🚫")
		return
	}

	fmt.Println("크론 작업 실행 시작... ✅")

	// 크론 작업 실행 (예: 데이터 백업)
	runCronTask()

	// 락 해제
	client.releaseLock(ctx, lockKey)
	fmt.Println("크론 작업 완료, 락 해제 🔓")
}

func runCronTask() {
	// 실제 실행할 크론 작업 (여기선 10초 대기)
	fmt.Println("작업 중... ⏳")
	time.Sleep(10 * time.Second)
	fmt.Println("작업 완료! 🎉")
}
