package mapreduce_example

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

// DemoContextTrace 演示如何在 MapReduce 中正确传递 context
// 问题：既要保证超时控制，又要让子协程能够继承 traceID
func DemoContextTrace() {
	fmt.Println("=== DemoContextTrace Start ===")

	// 模拟带有 traceID 的 context（通常来自 HTTP 请求或 RPC 调用）
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-12345678")
	ctx = context.WithValue(ctx, "user_id", "user-999")

	// 设置超时时间
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 方案1：使用 MapReduce + WithContext 选项（推荐）
	result, err := demoWithContext(ctx)
	if err != nil {
		log.Printf("MapReduce with context error: %v", err)
	} else {
		fmt.Printf("Result: %v\n", result)
	}

	fmt.Println("=== DemoContextTrace End ===")
}

// demoWithContext 使用 MapReduce + WithContext 正确传递 context
func demoWithContext(ctx context.Context) ([]int, error) {
	nums := []int{1, 2, 3, 4, 5}

	// 使用 MapReduce + WithContext 选项，它会正确处理 context 传递
	return mr.MapReduce(
		func(source chan<- int) {
			// Generator: 可以访问 ctx
			traceID := ctx.Value("trace_id")
			fmt.Printf("[Generator] TraceID: %v\n", traceID)

			for _, v := range nums {
				select {
				case <-ctx.Done():
					// 检查是否超时或取消
					fmt.Println("[Generator] Context cancelled")
					return
				case source <- v:
				}
			}
		},
		func(item int, writer mr.Writer[int], cancel func(error)) {
			// Mapper: 这里的 ctx 已经被 go-zero 处理过
			// 会继承 traceID 等值，但不会因为父 context 超时而立即取消
			traceID := ctx.Value("trace_id")
			userID := ctx.Value("user_id")

			// 使用 logx 记录日志，traceID 会自动传递
			logx.WithContext(ctx).Infof("Processing item: %d, TraceID: %v, UserID: %v",
				item, traceID, userID)

			// 模拟耗时操作
			time.Sleep(100 * time.Millisecond)

			result := item * item
			writer.Write(result)
		},
		func(pipe <-chan int, writer mr.Writer[[]int], cancel func(error)) {
			// Reducer: 聚合结果
			traceID := ctx.Value("trace_id")
			fmt.Printf("[Reducer] TraceID: %v\n", traceID)

			var results []int
			for v := range pipe {
				results = append(results, v)
			}
			writer.Write(results)
		},
		mr.WithContext(ctx),
	)
}

// DemoContextSeparation 演示手动分离 context 的方法（适用于 Go 1.21+）
func DemoContextSeparation() {
	fmt.Println("=== DemoContextSeparation Start ===")

	// 原始 context 带超时
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-87654321")
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	nums := []int{1, 2, 3, 4, 5}

	// 方案2：手动分离 context（Go 1.21+）
	// 使用 context.WithoutCancel 创建一个新的 context
	// 它会继承所有的值（如 traceID），但不会继承取消信号
	detachedCtx := context.WithoutCancel(ctx)

	result, err := mr.MapReduce(
		func(source chan<- int) {
			for _, v := range nums {
				source <- v
			}
		},
		func(item int, writer mr.Writer[int], cancel func(error)) {
			// 使用分离后的 context
			traceID := detachedCtx.Value("trace_id")
			logx.WithContext(detachedCtx).Infof("Processing with detached context: %d, TraceID: %v",
				item, traceID)

			time.Sleep(100 * time.Millisecond)
			writer.Write(item * 2)
		},
		func(pipe <-chan int, writer mr.Writer[[]int], cancel func(error)) {
			var results []int
			for v := range pipe {
				results = append(results, v)
			}
			writer.Write(results)
		},
	)

	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Result: %v\n", result)
	}

	fmt.Println("=== DemoContextSeparation End ===")
}

// DemoManualContextCopy 演示手动复制 context 值的方法（适用于 Go 1.20 及以下）
func DemoManualContextCopy() {
	fmt.Println("=== DemoManualContextCopy Start ===")

	// 原始 context 带超时
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-11111111")
	ctx = context.WithValue(ctx, "user_id", "user-888")
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// 方案3：手动复制 context 值到新的 context（Go 1.20 及以下）
	newCtx := copyContextValues(ctx, "trace_id", "user_id")

	nums := []int{1, 2, 3, 4, 5}

	result, err := mr.MapReduce(
		func(source chan<- int) {
			for _, v := range nums {
				source <- v
			}
		},
		func(item int, writer mr.Writer[int], cancel func(error)) {
			// 使用新的 context
			traceID := newCtx.Value("trace_id")
			userID := newCtx.Value("user_id")

			logx.WithContext(newCtx).Infof("Processing: %d, TraceID: %v, UserID: %v",
				item, traceID, userID)

			time.Sleep(100 * time.Millisecond)
			writer.Write(item * 3)
		},
		func(pipe <-chan int, writer mr.Writer[[]int], cancel func(error)) {
			var results []int
			for v := range pipe {
				results = append(results, v)
			}
			writer.Write(results)
		},
	)

	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Result: %v\n", result)
	}

	fmt.Println("=== DemoManualContextCopy End ===")
}

// copyContextValues 手动复制 context 中的值到新的 context
// 适用于 Go 1.20 及以下版本
func copyContextValues(ctx context.Context, keys ...string) context.Context {
	newCtx := context.Background()
	for _, key := range keys {
		if val := ctx.Value(key); val != nil {
			newCtx = context.WithValue(newCtx, key, val)
		}
	}
	return newCtx
}

// DemoRealWorldExample 真实场景示例：批量查询用户信息
func DemoRealWorldExample() {
	fmt.Println("=== DemoRealWorldExample Start ===")

	// 模拟 HTTP 请求的 context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "http-trace-99999")
	ctx = context.WithValue(ctx, "request_id", "req-12345")
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	userIDs := []int64{1001, 1002, 1003, 1004, 1005}

	// 使用 MapReduce + WithContext 批量查询
	users, err := mr.MapReduce(
		func(source chan<- int64) {
			for _, uid := range userIDs {
				select {
				case <-ctx.Done():
					logx.WithContext(ctx).Error("Generator cancelled due to timeout")
					return
				case source <- uid:
				}
			}
		},
		func(uid int64, writer mr.Writer[*UserInfo], cancel func(error)) {
			// 查询单个用户信息
			user, err := queryUserInfo(ctx, uid)
			if err != nil {
				// 记录错误但不取消整个任务
				logx.WithContext(ctx).Errorf("Failed to query user %d: %v", uid, err)
				return
			}
			writer.Write(user)
		},
		func(pipe <-chan *UserInfo, writer mr.Writer[[]*UserInfo], cancel func(error)) {
			var users []*UserInfo
			for user := range pipe {
				users = append(users, user)
			}
			writer.Write(users)
		},
		mr.WithContext(ctx),
	)

	if err != nil {
		log.Printf("Batch query failed: %v", err)
	} else {
		fmt.Printf("Successfully queried %d users\n", len(users))
		for _, user := range users {
			fmt.Printf("User: ID=%d, Name=%s\n", user.ID, user.Name)
		}
	}

	fmt.Println("=== DemoRealWorldExample End ===")
}

// UserInfo 用户信息
type UserInfo struct {
	ID   int64
	Name string
}

// queryUserInfo 模拟查询用户信息
func queryUserInfo(ctx context.Context, uid int64) (*UserInfo, error) {
	traceID := ctx.Value("trace_id")
	requestID := ctx.Value("request_id")

	// 记录日志，traceID 会自动传递
	logx.WithContext(ctx).Infof("Querying user %d, TraceID: %v, RequestID: %v",
		uid, traceID, requestID)

	// 模拟数据库查询
	time.Sleep(50 * time.Millisecond)

	return &UserInfo{
		ID:   uid,
		Name: fmt.Sprintf("User-%d", uid),
	}, nil
}
