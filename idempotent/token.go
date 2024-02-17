package idempotent

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// 通过token 机制实现接口的幂等性 https://www.cnblogs.com/taojietaoge/p/15845948.html
var redisClient *redis.Client

const (
	Addr           = ""
	Password       = ""
	DB             = 1
	ExpirationTime = time.Minute
)

func RedisInit() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       DB,
	})
	return redisClient
}

// GenerateToken 1.生成token
func GenerateToken() string {
	return uuid.New().String()
}

// ProcessRequest 2.模拟请求
func ProcessRequest(requestID string) (string, error) {
	redisClient := RedisInit()
	// 检查Redis中是否存在请求标识符
	existingResult, err := redisClient.Get(requestID).Result()
	if err == nil {
		// 如果存在，直接返回之前的处理结果
		return existingResult, nil
	} else if err != redis.Nil {
		// 处理其他Redis错误
		return "", err
	}
	result := ""
	// 处理请求逻辑...
	result, err = Work(requestID)
	if err != nil {
		return "", err
	}
	// 将请求标识符存储到Redis，并设置过期时间
	err = redisClient.SetNX(requestID, result, ExpirationTime).Err()
	if err != nil {
		// 处理存储状态错误
		return "", err
	}
	return result, nil
}

// Work 3.模拟处理逻辑
func Work(requestID string) (string, error) {
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
		return "success result", nil
	case <-time.After(time.Second * 2):
		fmt.Println("超时...记录日志...")
		return "", errors.New("timeout")
	}
}
