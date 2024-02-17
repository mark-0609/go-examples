package idempotent

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 通过token 机制实现接口的幂等性 https://www.cnblogs.com/taojietaoge/p/15845948.html
func Idempotent(token string) string {
	// fmt.Println("generate token...")
	// token := GenerateToken()

	err := RedisSetnXToken(token)
	if err != nil {
		fmt.Println("redis err:", err)
		return ""
	}
	// 向第三方执行http请求逻辑，若超时，自动重试/记录日志
	channel := make(chan struct{}, 1)
	defer close(channel)
	go func() {
		// 模拟超时
		// time.Sleep(time.Second * 3)
		// 正常执行
		// time.Sleep(time.Second)
		channel <- struct{}{}
	}()
	select {
	case <-channel:
		fmt.Println("正常执行结束...")
		return token
	case <-time.After(time.Second * 2):
		fmt.Println("超时...记录日志...")
		return "TimeOutToken"
	}

}

var redisClient *redis.Client

const (
	Addr     = "114.132.210.241:36379"
	Password = "G62m50oigInC30sf"
	DB       = 1
)

func RedisInit() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       DB,
	})
	// fmt.Println("redisClient:%v", redisClient)
	return redisClient
}

func RedisSetnXToken(token string) error {
	redisClient = RedisInit()
	defer redisClient.Close()
	redisClient.SetNX(token, 1, time.Minute)
	return nil
}

func GenerateToken() string {
	return time.Now().String() + strconv.Itoa(rand.Int())
}
