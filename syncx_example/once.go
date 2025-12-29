package syncx_example

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
)

// OnceDemo 模拟配置懒加载（支持失败重试）
// 场景：系统启动后，首次访问时加载配置。
// 使用 syncx.ImmutableResource 来实现资源的懒加载。
// 它确保资源只被加载一次，并且支持加载失败后的重试策略。
func OnceDemo() {
	fmt.Println("=== 开始 ImmutableResource (Once) 示例 ===")

	// 模拟前 3 次尝试都失败，第 4 次成功
	var tryCount int
	var lock sync.Mutex // 保护 tryCount

	loadConfig := func() (any, error) {
		lock.Lock()
		tryCount++
		currentCount := tryCount
		lock.Unlock()

		fmt.Printf("尝试加载配置 (第 %d 次)...\n", currentCount)

		if currentCount < 4 {
			return nil, errors.New("网络超时")
		}

		return "ServerPort=8080", nil
	}

	// 创建 ImmutableResource
	// WithRefreshIntervalOnFailure(0) 表示失败后下次立即重试，不等待
	resource := syncx.NewImmutableResource(loadConfig, syncx.WithRefreshIntervalOnFailure(0))

	// 模拟并发访问
	var wg sync.WaitGroup

	// 启动 5 个协程并发尝试加载
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 稍微错开一点时间，模拟持续请求
			time.Sleep(time.Millisecond * time.Duration(id*50))

			val, err := resource.Get()
			if err != nil {
				fmt.Printf("协程 %d 加载失败: %v\n", id, err)
			} else {
				fmt.Printf("协程 %d 获取配置: %s\n", id, val)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("=== 结束 ImmutableResource (Once) 示例 ===")
}
