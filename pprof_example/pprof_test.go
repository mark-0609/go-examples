package pprof_example

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // 必须引入以开启 HTTP pprof
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
	"testing"
	"time"
)

// TestPprof 演示了如何在异常（Panic）或手动触发（信号）时获取 pprof 数据
// 运行命令: go test -v -run TestPprof
func TestPprof(t *testing.T) {
	// 1. 开启 HTTP pprof 服务
	// 这是最推荐的方式，生产环境通常都会开启，配合防火墙或鉴权使用
	// 访问 http://localhost:6060/debug/pprof/ 可以查看各种指标
	go func() {
		log.Println("Starting HTTP pprof server at :6060")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Printf("HTTP pprof server failed: %v", err)
		}
	}()

	// 2. 监听信号，支持手动触发 Dump
	// 在发现程序假死（死锁）但未 Crash 时，可以发送信号触发 Dump
	setupSignalHandler()

	log.Println("Test started. Waiting 2 seconds before simulating a panic...")

	// 3. 模拟触发 Panic (在子协程中进行，以免中断主测试流程)
	go func() {
		// 设置 Panic 捕获，自动 Dump 堆栈
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
				dumpGoroutineStack("panic")
			}
		}()

		time.Sleep(2 * time.Second)
		simulatePanic()
	}()

	// 4. 模拟死锁 (如果不 panic)
	// 放在子协程中，模拟部分死锁场景
	go func() {
		time.Sleep(5 * time.Second)
		// simulateDeadlock() // 取消注释以模拟死锁
	}()

	log.Println("Test is running. You can visit http://localhost:6060/debug/pprof/")
	log.Println("Test will finish in 30 seconds automatically.")

	// 保持测试运行一段时间，以便观察
	time.Sleep(30 * time.Second)
}

func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	// 监听中断信号 (Ctrl+C) 和终止信号
	// 在 Linux 下通常使用 syscall.SIGUSR1 来专门触发 dump 而不退出程序
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for sig := range c {
			log.Printf("Received signal: %v", sig)

			// 收到信号时 dump pprof
			dumpGoroutineStack("signal")

			// 如果是退出信号，则退出
			if sig == os.Interrupt || sig == syscall.SIGTERM {
				log.Println("Exiting...")
				os.Exit(0)
			}
		}
	}()
}

// dumpGoroutineStack 将当前所有 goroutine 的堆栈信息写入文件
func dumpGoroutineStack(reason string) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("goroutine_dump_%s_%s.pprof", reason, timestamp)

	f, err := os.Create(filename)
	if err != nil {
		log.Printf("Could not create pprof file: %v", err)
		return
	}
	defer f.Close()

	// 获取所有 goroutine 的堆栈信息
	// debug=0: 二进制格式 (go tool pprof 使用)
	// debug=1: 文本格式 (人类可读)
	// debug=2: 文本格式，包含所有 goroutine 的完整堆栈 (类似 panic 时的输出)

	// 这里我们使用 pprof.Lookup("goroutine").WriteTo 输出二进制格式，方便工具分析
	// 如果想要文本格式，可以使用 pprof.Lookup("goroutine").WriteTo(f, 2)
	// 或者手动 runtime.Stack()

	if err := pprof.Lookup("goroutine").WriteTo(f, 0); err != nil {
		log.Printf("Could not write pprof: %v", err)
	} else {
		log.Printf("Successfully dumped goroutine stack to %s", filename)
	}
}

func simulatePanic() {
	log.Println("Simulating panic now!")
	panic("something went wrong!")
}

func simulateDeadlock() {
	log.Println("Simulating partial deadlock...")
	var mu sync.Mutex
	mu.Lock()
	// 再次加锁导致死锁
	mu.Lock()
}
