package syncx_example

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
)

// SingleFlightDemo 模拟热点商品详情查询
// 场景：高并发下，多个请求同时查询同一个 Key（如商品ID），
// 使用 SingleFlight 确保同一时刻只有一个请求穿透到数据库，
// 其他请求共享该结果，从而保护数据库。
func SingleFlightDemo() {
	fmt.Println("=== 开始 SingleFlight 示例 ===")

	// 创建 SingleFlight 组
	g := syncx.NewSingleFlight()

	// 模拟并发请求数
	concurrentReqs := 10
	// 模拟数据库调用次数
	var dbCallCount int32

	var wg sync.WaitGroup
	wg.Add(concurrentReqs)

	key := "product_12345"

	for i := 0; i < concurrentReqs; i++ {
		go func(id int) {
			defer wg.Done()

			// Do 方法：针对同一个 key，函数只会被执行一次
			val, err := g.Do(key, func() (interface{}, error) {
				fmt.Printf("[DB] 正在从数据库加载数据 (请求ID: %d)...\n", id)
				atomic.AddInt32(&dbCallCount, 1)
				time.Sleep(time.Millisecond * 200) // 模拟耗时
				return "Product: iPhone 15 Pro", nil
			})

			if err != nil {
				fmt.Printf("请求 %d 失败: %v\n", id, err)
			} else {
				fmt.Printf("请求 %d 获取结果: %v\n", id, val)
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("总请求数: %d, 实际数据库调用次数: %d\n", concurrentReqs, dbCallCount)
	fmt.Println("=== 结束 SingleFlight 示例 ===")
}
