package concurrent_executor

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// ConcurrentTaskResult 并发任务执行结果（泛型版本，支持任意 Key 和 Value 类型）
type ConcurrentTaskResult[K comparable, V any] struct {
	Key   K     // 任务标识（支持任意可比较类型）
	Value V     // 任务结果（支持任意类型）
	Err   error // 任务错误
}

// ExecutorOptions 并发执行器配置选项
type ExecutorOptions struct {
	MaxConcurrency int           // 最大并发数，0表示不限制
	Timeout        time.Duration // 单个任务超时时间，0表示不设置超时
}

// DefaultExecutorOptions 返回默认配置
func DefaultExecutorOptions() *ExecutorOptions {
	return &ExecutorOptions{
		MaxConcurrency: 0, // 不限制并发数
		Timeout:        0, // 不设置超时
	}
}

// ConcurrentExecutorWithOptions 带配置选项的并发执行器（字符串Key版本）
func ConcurrentExecutorWithOptions[T any](ctx context.Context, keys []string,
	taskFunc func(context.Context, string) (T, error), opts *ExecutorOptions) ([]ConcurrentTaskResult[string, T], error) {
	return ConcurrentExecutorGeneric(ctx, keys, taskFunc, opts)
}

// ConcurrentExecutorGeneric 泛型并发执行器，支持任意可比较类型的Key
// K: Key类型约束为comparable（可比较类型，如string、int、自定义struct等）
// V: Value类型无约束，支持任意类型
func ConcurrentExecutorGeneric[K comparable, V any](ctx context.Context, keys []K,
	taskFunc func(context.Context, K) (V, error), opts *ExecutorOptions) ([]ConcurrentTaskResult[K, V], error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("keys cannot be empty")
	}

	if opts == nil {
		opts = DefaultExecutorOptions()
	}

	resultChan := make(chan ConcurrentTaskResult[K, V], len(keys))
	var wg sync.WaitGroup

	var semaphore chan struct{}
	if opts.MaxConcurrency > 0 {
		semaphore = make(chan struct{}, opts.MaxConcurrency)
	}

	for _, key := range keys {
		// 1. 优先检查 Context 是否已取消
		// 如果 Context 已取消，不再启动新任务，直接返回错误结果
		select {
		case <-ctx.Done():
			resultChan <- ConcurrentTaskResult[K, V]{
				Key: key,
				Err: ctx.Err(),
			}
			continue
		default:
		}

		// 2. 获取信号量（带 Context 取消监听）
		if semaphore != nil {
			select {
			case semaphore <- struct{}{}:
				// 成功获取信号量
			case <-ctx.Done():
				// 等待信号量期间 Context 被取消
				resultChan <- ConcurrentTaskResult[K, V]{
					Key: key,
					Err: ctx.Err(),
				}
				continue
			}
		}

		wg.Add(1)

		go func(ctx context.Context, k K) {
			// 标记任务是否已完成发送结果，防止重复发送导致死锁
			var finished bool
			defer func() {
				defer wg.Done()

				// 释放信号量
				if semaphore != nil {
					<-semaphore
				}
				// 捕获 panic，确保在发生 panic 时也能追踪到具体的 key 和错误信息
				if r := recover(); r != nil {
					// 只有在未发送结果的情况下才发送 panic 错误
					// 如果已经发送了结果（finished=true）却发生了 panic（极少见，例如在发送后的清理逻辑中），
					// 则不再发送结果，避免 resultChan 阻塞导致死锁
					if !finished {
						err := fmt.Errorf("panic recovered: %v, stack: %s", r, string(debug.Stack()))
						resultChan <- ConcurrentTaskResult[K, V]{
							Key: k,
							Err: err,
						}
					} else {
						// 记录日志或忽略，因为结果已发送
						fmt.Printf("panic recovered after result sent: %v\n", r)
					}
				}
			}()

			// 如果设置了超时，添加超时控制
			var execCtx = ctx
			if opts.Timeout > 0 {
				var cancel context.CancelFunc
				execCtx, cancel = context.WithTimeout(ctx, opts.Timeout)
				defer cancel()
			}
			value, err := taskFunc(execCtx, k)
			// 无论成功还是失败，都返回结果，确保结果集完整
			resultChan <- ConcurrentTaskResult[K, V]{
				Key:   k,
				Value: value,
				Err:   err,
			}
			finished = true

		}(ctx, key)
	}
	wg.Wait()
	close(resultChan)

	// 收集结果
	results := make([]ConcurrentTaskResult[K, V], 0, len(keys))
	for result := range resultChan {
		results = append(results, result)
	}

	return results, nil
}

// ConcurrentExecutorMap 并发执行器，返回Map结果（方便按Key查找）
func ConcurrentExecutorMap[K comparable, V any](ctx context.Context, keys []K, taskFunc func(context.Context, K) (V, error), opts *ExecutorOptions) (map[K]V, error) {
	results, err := ConcurrentExecutorGeneric(ctx, keys, taskFunc, opts)
	if err != nil {
		return nil, err
	}

	resultMap := make(map[K]V, len(results))
	for _, result := range results {
		if result.Err == nil {
			resultMap[result.Key] = result.Value
		}
	}

	return resultMap, nil
}
