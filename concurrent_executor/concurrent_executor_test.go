package concurrent_executor

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestConcurrentExecutorGeneric_Timeout(t *testing.T) {
	keys := []string{"task1", "task2"}
	opts := &ExecutorOptions{
		Timeout: 100 * time.Millisecond,
	}

	taskFunc := func(ctx context.Context, key string) (string, error) {
		select {
		case <-time.After(200 * time.Millisecond):
			return "done", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	results, err := ConcurrentExecutorGeneric(context.Background(), keys, taskFunc, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, res := range results {
		if res.Err != context.DeadlineExceeded {
			t.Errorf("expected DeadlineExceeded error for key %s, got %v", res.Key, res.Err)
		}
	}
}

func TestConcurrentExecutorGeneric_Panic(t *testing.T) {
	keys := []string{"panic_task1", "panic_task2"}
	opts := DefaultExecutorOptions()

	taskFunc := func(ctx context.Context, key string) (string, error) {
		panic("something went wrong")
	}

	results, err := ConcurrentExecutorGeneric(context.Background(), keys, taskFunc, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range results {
		if r.Err == nil {
			t.Error("expected error from panic, got nil")
		} else {
			t.Logf("Got expected key: %v ,panic error:%v", r.Key, r.Err)
		}
	}
}

func TestConcurrentExecutorGeneric_ContextPropagation(t *testing.T) {
	type contextKey string
	const traceKey contextKey = "trace_id"

	keys := []string{"task1"}
	opts := DefaultExecutorOptions()

	ctx := context.WithValue(context.Background(), traceKey, "12345")

	taskFunc := func(ctx context.Context, key string) (string, error) {
		val := ctx.Value(traceKey)
		if val != "12345" {
			return "", fmt.Errorf("trace_id mismatch: got %v", val)
		}
		return "ok", nil
	}

	results, err := ConcurrentExecutorGeneric(ctx, keys, taskFunc, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if results[0].Err != nil {
		t.Errorf("task failed: %v", results[0].Err)
	}
}

func demo(ctx context.Context) (int, error) {
	return rand.Int(), nil
}

func TestConcurrentExecutorWithOptions(t *testing.T) {
	taskFunc := func(ctx context.Context, playId string) (int, error) {
		return demo(ctx)
	}

	keys := []string{"1", "2", "3"}

	// 使用通用并发执行器
	results, err := ConcurrentExecutorWithOptions(context.Background(), keys, taskFunc, nil)
	if err != nil {
		fmt.Printf("concurrent execution failed: %v\n", err)
		return
	}

	// 3. 收集查询结果
	for _, result := range results {
		if result.Err != nil {
			fmt.Printf("result failed: %s\n", result.Err)
		} else {
			fmt.Printf("result success,val: %v\n", result.Value)
		}
	}

}

func TestConcurrentExecutorGeneric_Cancel(t *testing.T) {
	keys := []string{"task1", "task2", "task3"}
	// 限制并发为1，确保任务串行启动，方便控制取消时机
	opts := &ExecutorOptions{
		MaxConcurrency: 1,
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 用于记录实际执行的任务数量
	var executedCount int

	taskFunc := func(ctx context.Context, key string) (string, error) {
		executedCount++
		if key == "task1" {
			// 在第一个任务中取消 Context
			cancel()
			// 模拟一点耗时，确保主循环有机会检测到 Done
			time.Sleep(50 * time.Millisecond)
			return "done", nil
		}
		// 其他任务不应该被执行
		return "should_not_run", nil
	}

	results, err := ConcurrentExecutorGeneric(ctx, keys, taskFunc, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 验证结果
	for _, res := range results {
		if res.Key == "task1" {
			if res.Err != nil {
				// task1 可能会成功，也可能因为 sleep 期间 ctx 被取消而报错，这取决于具体实现细节
				// 但我们的重点是后续任务
			}
		} else {
			// task2 和 task3 应该是因为 Context 取消而失败
			if res.Err != context.Canceled {
				t.Errorf("expected context.Canceled for key %s, got %v", res.Key, res.Err)
			}
		}
	}

	// 验证实际执行的任务数
	// 理论上只有 task1 执行了。task2 和 task3 应该在主循环中被拦截。
	// 注意：由于并发和调度的不确定性，如果 task1 取消得太慢，task2 可能会刚好进入 semaphore 等待。
	// 但在我们的优化逻辑中，semaphore 等待也会监听 ctx.Done()。
	// 所以 executedCount 应该是 1。
	if executedCount != 1 {
		t.Errorf("expected executedCount to be 1, got %d", executedCount)
	}
}
