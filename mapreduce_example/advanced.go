package mapreduce_example

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/zeromicro/go-zero/core/mr"
)

// DemoMultiStage 演示多级扇入扇出
// 场景：
// 1. 第一阶段：从大量用户ID中筛选出活跃用户（模拟过滤）
// 2. 第二阶段：并发获取活跃用户的订单金额并计算总额（模拟聚合）
func DemoMultiStage() {
	fmt.Println("=== DemoMultiStage Start ===")

	// 模拟 100 个用户ID
	totalUsers := 100
	userIDs := make([]int, totalUsers)
	for i := 0; i < totalUsers; i++ {
		userIDs[i] = i
	}

	// 第一阶段：筛选活跃用户 (Fan-out -> Fan-in)
	// 返回类型为 []int
	activeUsersResult, err := mr.MapReduce(func(source chan<- int) {
		for _, uid := range userIDs {
			source <- uid
		}
	}, func(item int, writer mr.Writer[int], cancel func(error)) {
		// 模拟检查用户是否活跃（随机）
		if isUserActive(item) {
			writer.Write(item)
		}
	}, func(pipe <-chan int, writer mr.Writer[[]int], cancel func(error)) {
		var active []int
		for item := range pipe {
			active = append(active, item)
		}
		writer.Write(active)
	})

	if err != nil {
		log.Printf("Stage 1 error: %v", err)
		return
	}

	activeUsers := activeUsersResult
	fmt.Printf("Stage 1 complete. Active users count: %d\n", len(activeUsers))

	// 第二阶段：计算活跃用户的订单总额 (Fan-out -> Fan-in)
	// 输入是第一阶段的结果
	totalAmountResult, err := mr.MapReduce(func(source chan<- int) {
		for _, uid := range activeUsers {
			source <- uid
		}
	}, func(item int, writer mr.Writer[float64], cancel func(error)) {
		// 模拟并发获取订单金额
		amount := fetchUserOrderAmount(item)
		writer.Write(amount)
	}, func(pipe <-chan float64, writer mr.Writer[float64], cancel func(error)) {
		var total float64
		for item := range pipe {
			total += item
		}
		writer.Write(total)
	})

	if err != nil {
		log.Printf("Stage 2 error: %v", err)
		return
	}

	fmt.Printf("Stage 2 complete. Total Order Amount: %.2f\n", totalAmountResult)
	fmt.Println("=== DemoMultiStage End ===")
}

func isUserActive(uid int) bool {
	// 模拟耗时
	time.Sleep(time.Millisecond * 2)
	// 偶数ID为活跃用户
	return uid%2 == 0
}

func fetchUserOrderAmount(uid int) float64 {
	// 模拟网络请求耗时
	time.Sleep(time.Millisecond * 5)
	// 随机金额 10-100
	return 10.0 + rand.Float64()*90.0
}
