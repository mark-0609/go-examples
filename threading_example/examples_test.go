package threading_example

import (
	"fmt"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/threading"
)

// TestGoSafeSimple 简单的 GoSafe 测试
func TestGoSafeSimple(t *testing.T) {
	t.Run("AsyncTask", func(t *testing.T) {
		done := make(chan bool)

		threading.GoSafe(func() {
			time.Sleep(100 * time.Millisecond)
			done <- true
		})

		select {
		case <-done:
			t.Log("异步任务执行成功")
		case <-time.After(1 * time.Second):
			t.Error("异步任务超时")
		}
	})

	t.Run("PanicRecovery", func(t *testing.T) {
		threading.GoSafe(func() {
			panic("测试 panic")
		})

		time.Sleep(100 * time.Millisecond)
		t.Log("panic 已被捕获，程序继续运行")
	})
}

// TestRunSafeSimple 简单的 RunSafe 测试
func TestRunSafeSimple(t *testing.T) {
	t.Run("NormalExecution", func(t *testing.T) {
		executed := false

		threading.RunSafe(func() {
			executed = true
		})

		if !executed {
			t.Error("函数未执行")
		}
	})

	t.Run("PanicRecovery", func(t *testing.T) {
		threading.RunSafe(func() {
			panic("测试 panic")
		})

		t.Log("panic 已被捕获，测试继续")
	})
}

// TestRoutineGroupSimple 简单的 RoutineGroup 测试
func TestRoutineGroupSimple(t *testing.T) {
	t.Run("WaitForAll", func(t *testing.T) {
		group := threading.NewRoutineGroup()
		counter := 0

		for i := 0; i < 5; i++ {
			group.RunSafe(func() {
				time.Sleep(100 * time.Millisecond)
				counter++
			})
		}

		group.Wait()

		if counter != 5 {
			t.Errorf("期望执行 5 次，实际执行 %d 次", counter)
		}
	})
}

// TestAllThreadingFeatures 测试所有 threading 功能
func TestAllThreadingFeatures(t *testing.T) {
	fmt.Println("========================================")
	fmt.Println("go-zero threading 包完整测试")
	fmt.Println("========================================")

	// 1. GoSafe 测试
	t.Run("GoSafe", func(t *testing.T) {
		fmt.Println("\n【1. GoSafe - 安全启动协程】")

		t.Run("Basic", func(t *testing.T) {
			fmt.Println("测试基本功能...")
			done := make(chan bool)

			threading.GoSafe(func() {
				fmt.Println("  后台任务执行中...")
				time.Sleep(100 * time.Millisecond)
				done <- true
			})

			<-done
			fmt.Println("  ✅ 基本功能测试通过")
		})

		t.Run("WithPanic", func(t *testing.T) {
			fmt.Println("测试 panic 捕获...")
			threading.GoSafe(func() {
				panic("模拟 panic")
			})

			time.Sleep(100 * time.Millisecond)
			fmt.Println("  ✅ Panic 捕获测试通过")
		})
	})

	// 2. RunSafe 测试
	t.Run("RunSafe", func(t *testing.T) {
		fmt.Println("\n【2. RunSafe - 当前协程安全执行】")

		t.Run("Basic", func(t *testing.T) {
			fmt.Println("测试基本功能...")
			executed := false

			threading.RunSafe(func() {
				executed = true
			})

			if !executed {
				t.Error("函数未执行")
			}
			fmt.Println("  ✅ 基本功能测试通过")
		})

		t.Run("WithPanic", func(t *testing.T) {
			fmt.Println("测试 panic 捕获...")
			threading.RunSafe(func() {
				panic("模拟 panic")
			})
			fmt.Println("  ✅ Panic 捕获测试通过")
		})
	})

	// 3. RoutineGroup 测试
	t.Run("RoutineGroup", func(t *testing.T) {
		fmt.Println("\n【3. RoutineGroup - 协程组管理】")

		t.Run("WaitForAll", func(t *testing.T) {
			fmt.Println("测试等待所有协程完成...")
			group := threading.NewRoutineGroup()
			counter := 0

			for i := 0; i < 5; i++ {
				group.RunSafe(func() {
					time.Sleep(100 * time.Millisecond)
					counter++
				})
			}

			group.Wait()

			if counter == 5 {
				fmt.Println("  ✅ 所有协程执行完成")
			} else {
				t.Errorf("期望执行 5 次，实际执行 %d 次", counter)
			}
		})
	})

	fmt.Println("\n========================================")
	fmt.Println("所有测试执行完成！")
	fmt.Println("========================================")
}
