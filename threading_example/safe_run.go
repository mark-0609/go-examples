package threading_example

import (
	"fmt"
	"sync"
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

// downloadFiles 并发下载文件
func downloadFiles(urls []string) {
	group := threading.NewRoutineGroup()

	for _, url := range urls {
		url := url // 创建局部副本
		group.RunSafe(func() {
			fmt.Printf("开始下载: %s\n", url)
			time.Sleep(time.Second) // 模拟下载
			fmt.Printf("完成下载: %s\n", url)
		})
	}

	fmt.Println("等待所有下载完成...")
	group.Wait()
	fmt.Println("所有文件下载完成！")
}

func SafeRunDemo2() {
	urls := []string{
		"http://example.com/file1.zip",
		"http://example.com/file2.zip",
		"http://example.com/file3.zip",
	}
	downloadFiles(urls)
}

// queryMultipleTables 并发数据库查询
func queryMultipleTables(tableNames []string) map[string]struct{} {
	group := threading.NewRoutineGroup()
	results := make(map[string]struct{})
	var mu sync.Mutex

	for _, table := range tableNames {
		table := table
		group.RunSafe(func() {
			records := struct{}{}

			mu.Lock()
			results[table] = records
			mu.Unlock()
		})
	}

	group.Wait()
	return results
}
