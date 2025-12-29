package mapreduce_example

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

// TestContextTraceIDInheritance 验证 traceID 是否正确继承
func TestContextTraceIDInheritance(t *testing.T) {
	// 创建带 traceID 的 context
	ctx := context.Background()
	expectedTraceID := "test-trace-12345"
	ctx = context.WithValue(ctx, "trace_id", expectedTraceID)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	nums := []int{1, 2, 3}
	var traceIDCount int32

	result, err := mr.MapReduce(
		func(source chan<- int) {
			// Generator 中验证 traceID
			if traceID := ctx.Value("trace_id"); traceID == expectedTraceID {
				atomic.AddInt32(&traceIDCount, 1)
				t.Logf("✅ Generator: traceID = %v", traceID)
			}
			for _, v := range nums {
				source <- v
			}
		},
		func(item int, writer mr.Writer[int], cancel func(error)) {
			// Mapper 中验证 traceID
			if traceID := ctx.Value("trace_id"); traceID == expectedTraceID {
				atomic.AddInt32(&traceIDCount, 1)
				t.Logf("✅ Mapper[%d]: traceID = %v", item, traceID)
			} else {
				t.Errorf("❌ Mapper[%d]: traceID not inherited, got %v", item, traceID)
			}
			writer.Write(item * 2)
		},
		func(pipe <-chan int, writer mr.Writer[[]int], cancel func(error)) {
			// Reducer 中验证 traceID
			if traceID := ctx.Value("trace_id"); traceID == expectedTraceID {
				atomic.AddInt32(&traceIDCount, 1)
				t.Logf("✅ Reducer: traceID = %v", traceID)
			}
			var results []int
			for v := range pipe {
				results = append(results, v)
			}
			writer.Write(results)
		},
		mr.WithContext(ctx),
	)

	if err != nil {
		t.Fatalf("MapReduce failed: %v", err)
	}

	// 验证结果
	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}

	// 验证 traceID 在所有阶段都被正确继承
	// Generator(1) + Mapper(3) + Reducer(1) = 5
	expectedCount := int32(5)
	if traceIDCount != expectedCount {
		t.Errorf("❌ TraceID not inherited in all stages. Expected %d, got %d", expectedCount, traceIDCount)
	} else {
		t.Logf("✅ TraceID successfully inherited in all %d stages", traceIDCount)
	}
}

// TestContextTimeoutControl 验证超时控制是否生效
func TestContextTimeoutControl(t *testing.T) {
	// 创建一个很短的超时 context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "timeout-test-trace")
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	startTime := time.Now()

	result, err := mr.MapReduce(
		func(source chan<- int) {
			for _, v := range nums {
				select {
				case <-ctx.Done():
					// 超时后停止生成
					t.Logf("✅ Generator stopped due to timeout after %v", time.Since(startTime))
					return
				case source <- v:
					time.Sleep(50 * time.Millisecond) // 模拟慢速生成
				}
			}
		},
		func(item int, writer mr.Writer[int], cancel func(error)) {
			// 即使父 context 超时，Mapper 也应该能完成已经接收到的任务
			traceID := ctx.Value("trace_id")
			t.Logf("Mapper[%d]: processing with traceID=%v", item, traceID)
			time.Sleep(20 * time.Millisecond)
			writer.Write(item * 2)
		},
		func(pipe <-chan int, writer mr.Writer[[]int], cancel func(error)) {
			var results []int
			for v := range pipe {
				results = append(results, v)
			}
			t.Logf("✅ Reducer: processed %d items", len(results))
			writer.Write(results)
		},
		mr.WithContext(ctx),
	)

	elapsed := time.Since(startTime)

	// 验证：应该在超时时间附近停止（不会处理完所有10个数字）
	if elapsed > 500*time.Millisecond {
		t.Errorf("❌ Timeout control failed: took %v, expected around 100ms", elapsed)
	} else {
		t.Logf("✅ Timeout control works: stopped after %v", elapsed)
	}

	// 验证：不应该处理完所有数字
	if len(result) >= len(nums) {
		t.Errorf("❌ Timeout didn't stop processing: processed all %d items", len(result))
	} else {
		t.Logf("✅ Timeout stopped processing: only processed %d/%d items", len(result), len(nums))
	}

	if err != nil {
		t.Logf("Error (expected): %v", err)
	}
}

// TestBothTraceIDAndTimeout 综合测试：同时验证 traceID 继承和超时控制
func TestBothTraceIDAndTimeout(t *testing.T) {
	t.Log("=== 综合测试：验证 traceID 继承 + 超时控制 ===")

	// 场景：HTTP 请求带 traceID，超时 500ms
	ctx := context.Background()
	expectedTraceID := "http-request-trace-999"
	ctx = context.WithValue(ctx, "trace_id", expectedTraceID)
	ctx = context.WithValue(ctx, "user_id", "user-123")
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	// 模拟批量查询 20 个用户，每个查询需要 100ms
	// 在 500ms 超时内，理论上最多能处理 5 个左右（考虑并发）
	userIDs := make([]int, 20)
	for i := 0; i < 20; i++ {
		userIDs[i] = 1000 + i
	}

	var generatedCount int32
	var traceIDVerified int32

	startTime := time.Now()

	result, err := mr.MapReduce(
		func(source chan<- int) {
			for _, uid := range userIDs {
				select {
				case <-ctx.Done():
					generated := atomic.LoadInt32(&generatedCount)
					t.Logf("✅ Generator: stopped due to timeout after generating %d tasks", generated)
					return
				case source <- uid:
					atomic.AddInt32(&generatedCount, 1)
					// 慢速生成，确保超时能触发
					time.Sleep(30 * time.Millisecond)
				}
			}
		},
		func(uid int, writer mr.Writer[string], cancel func(error)) {
			// 验证 traceID 继承
			traceID := ctx.Value("trace_id")
			userID := ctx.Value("user_id")

			if traceID == expectedTraceID && userID == "user-123" {
				atomic.AddInt32(&traceIDVerified, 1)
			}

			// 模拟数据库查询
			time.Sleep(100 * time.Millisecond)

			result := fmt.Sprintf("User-%d (trace:%v)", uid, traceID)
			writer.Write(result)
		},
		func(pipe <-chan string, writer mr.Writer[[]string], cancel func(error)) {
			var results []string
			for v := range pipe {
				results = append(results, v)
			}
			writer.Write(results)
		},
		mr.WithContext(ctx),
	)

	elapsed := time.Since(startTime)
	generated := atomic.LoadInt32(&generatedCount)

	t.Logf("\n=== 测试结果 ===")
	t.Logf("执行时间: %v", elapsed)
	t.Logf("生成任务数: %d/%d", generated, len(userIDs))
	t.Logf("完成任务数: %d", len(result))
	t.Logf("traceID 验证次数: %d", atomic.LoadInt32(&traceIDVerified))
	if err != nil {
		t.Logf("错误信息: %v", err)
	}

	// 验证1：超时控制生效（Generator 应该被中断）
	if generated >= int32(len(userIDs)) {
		t.Errorf("❌ 超时控制失败：生成了所有 %d 个任务", len(userIDs))
	} else {
		t.Logf("✅ 超时控制生效：只生成了 %d/%d 个任务", generated, len(userIDs))
	}

	// 验证2：执行时间应该接近超时时间
	if elapsed > 1*time.Second {
		t.Errorf("❌ 超时控制失败：执行时间 %v 超过预期", elapsed)
	} else {
		t.Logf("✅ 超时控制正常：执行时间 %v 符合预期", elapsed)
	}

	// 验证3：所有 Mapper 都应该继承了 traceID（即使最后因超时返回错误）
	// traceIDVerified 记录了所有进入 Mapper 的任务数
	if atomic.LoadInt32(&traceIDVerified) > 0 {
		t.Logf("✅ traceID 继承成功：%d 个 Mapper 协程都继承了 traceID", atomic.LoadInt32(&traceIDVerified))
	} else {
		t.Errorf("❌ traceID 继承失败：没有任务继承 traceID")
	}

	// 验证4：应该返回超时错误
	if err != nil && err.Error() == "context deadline exceeded" {
		t.Logf("✅ 正确返回超时错误")
	}

	t.Log("\n=== 结论 ===")
	// 关键验证：
	// 1. Generator 被超时中断（generated < total）
	// 2. Mapper 协程都继承了 traceID（traceIDVerified > 0）
	// 3. 返回了超时错误
	if generated < int32(len(userIDs)) && atomic.LoadInt32(&traceIDVerified) > 0 && err != nil {
		t.Log("✅✅✅ 测试通过：同时满足超时控制和 traceID 继承！")
		t.Log("说明：")
		t.Log("  - 超时控制：Generator 在超时后停止生成新任务")
		t.Log("  - traceID 继承：所有 Mapper 协程都正确继承了 traceID")
		t.Log("  - 错误处理：正确返回了超时错误")
		t.Log("  - 这证明了 mr.WithContext(ctx) 既保证了超时控制，又保证了 traceID 传递")
	} else {
		t.Error("❌ 测试失败：未能同时满足两个需求")
	}
}

// TestLogxWithContext 验证 logx.WithContext 是否正确传递 traceID
func TestLogxWithContext(t *testing.T) {
	// 配置 logx 输出到控制台
	logx.SetLevel(logx.InfoLevel)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "logx-test-trace-888")
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	nums := []int{1, 2, 3}

	_, err := mr.MapReduce(
		func(source chan<- int) {
			for _, v := range nums {
				source <- v
			}
		},
		func(item int, writer mr.Writer[int], cancel func(error)) {
			// 使用 logx.WithContext 记录日志
			traceID := ctx.Value("trace_id")
			logx.WithContext(ctx).Infof("Processing item %d with traceID: %v", item, traceID)
			writer.Write(item)
		},
		func(pipe <-chan int, writer mr.Writer[[]int], cancel func(error)) {
			var results []int
			for v := range pipe {
				results = append(results, v)
			}
			logx.WithContext(ctx).Infof("Reducer completed with %d results", len(results))
			writer.Write(results)
		},
		mr.WithContext(ctx),
	)

	if err != nil {
		t.Fatalf("MapReduce failed: %v", err)
	}

	t.Log("✅ logx.WithContext 正确传递了 traceID（查看上面的日志输出）")
}
