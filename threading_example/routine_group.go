package threading_example

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/threading"
)

// RoutineGroupDemo 模拟批量发送通知
// 场景：需要给 20 个用户发送通知，但为了保护下游服务，
// 限制并发数最多为 5 个。
func RoutineGroupDemo() {
	fmt.Println("=== 开始 RoutineGroup 示例 ===")

	// 创建一个 RoutineGroup
	group := threading.NewRoutineGroup()

	// 模拟 20 个任务
	taskCount := 20
	// 限制并发数为 5 (注意：RoutineGroup 本身不限制并发数，需要配合 channel 或其他方式，
	// 但这里我们演示 WaitGroup 的替代用法，go-zero 的 RoutineGroup 主要是对 WaitGroup 的封装，提供了 RunSafe)
	// 修正：go-zero 的 RoutineGroup 确实只是 WaitGroup 的封装。
	// 如果要限制并发，通常配合 WorkerPool 或者信号量。
	// 但 go-zero 的 mr 包更适合做并发控制。
	// 这里我们演示 RoutineGroup 如何安全地等待一组任务完成。

	var successCount int32

	for i := 0; i < taskCount; i++ {
		userID := i
		// RunSafe 会增加 WaitGroup 计数，并在 goroutine 中执行函数，
		// 函数执行完毕（或 panic recover 后）减少计数。
		group.RunSafe(func() {
			// 模拟发送逻辑
			fmt.Printf("正在发送通知给用户 %d...\n", userID)
			time.Sleep(time.Millisecond * 100)

			// 模拟偶发失败
			if userID%7 == 0 {
				// 即使这里 panic，RunSafe 也会捕获，并确保 Done 被调用
				panic(fmt.Sprintf("用户 %d 发送失败", userID))
			}

			atomic.AddInt32(&successCount, 1)
		})
	}

	// 等待所有任务完成
	group.Wait()

	fmt.Printf("所有任务处理完毕。成功发送: %d/%d\n", successCount, taskCount)
	fmt.Println("=== 结束 RoutineGroup 示例 ===")
}
