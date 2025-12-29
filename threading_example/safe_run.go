package threading_example

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/threading"
)

// SafeRunDemo 模拟后台日志处理任务
// 场景：后台启动一个协程处理任务，任务中可能发生 Panic，
// 使用 GoSafe 确保 Panic 被捕获，不会导致整个进程崩溃。
func SafeRunDemo() {
	fmt.Println("=== 开始 SafeRun 示例 ===")

	// 模拟一个不稳定的任务
	riskyTask := func() {
		fmt.Println("后台任务正在运行...")
		time.Sleep(time.Millisecond * 500)

		// 模拟发生 Panic
		panic("遇到意外错误：日志格式解析失败")
	}

	// 使用 GoSafe 启动协程
	threading.GoSafe(riskyTask)

	// 主协程等待一会，确保看到 Panic 被捕获的日志（go-zero 默认会打印 recover 日志）
	time.Sleep(time.Second)
	fmt.Println("主进程依然存活，未受 Panic 影响")
	fmt.Println("=== 结束 SafeRun 示例 ===")
}
