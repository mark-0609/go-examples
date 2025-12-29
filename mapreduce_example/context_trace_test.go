package mapreduce_example

import (
	"context"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// TestContextTrace 测试 context 传递
func TestContextTrace(t *testing.T) {
	// 初始化日志
	logx.DisableStat()

	t.Run("MapReduce+WithContext", func(t *testing.T) {
		DemoContextTrace()
	})

	t.Run("ContextSeparation", func(t *testing.T) {
		DemoContextSeparation()
	})

	t.Run("ManualContextCopy", func(t *testing.T) {
		DemoManualContextCopy()
	})

	t.Run("RealWorldExample", func(t *testing.T) {
		DemoRealWorldExample()
	})
}

// TestDetachContext 测试 context 分离
func TestDetachContext(t *testing.T) {
	// 创建带超时的 context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "test-trace-123")
	ctx = context.WithValue(ctx, "user_id", "user-456")
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	// 分离 context
	detachedCtx := DetachContext(ctx, "trace_id", "user_id")

	// 等待原 context 超时
	time.Sleep(200 * time.Millisecond)

	// 原 context 应该已经超时
	if ctx.Err() == nil {
		t.Error("Expected original context to be cancelled")
	}

	// 分离的 context 不应该超时
	if detachedCtx.Err() != nil {
		t.Error("Expected detached context to not be cancelled")
	}

	// 分离的 context 应该保留值
	if traceID := detachedCtx.Value("trace_id"); traceID != "test-trace-123" {
		t.Errorf("Expected trace_id to be test-trace-123, got %v", traceID)
	}

	if userID := detachedCtx.Value("user_id"); userID != "user-456" {
		t.Errorf("Expected user_id to be user-456, got %v", userID)
	}
}

// TestContextValueExtractor 测试 context 值提取器
func TestContextValueExtractor(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-999")
	ctx = context.WithValue(ctx, "request_id", "req-888")
	ctx = context.WithValue(ctx, "user_id", "user-777")

	extractor := NewContextValueExtractor(ctx)

	if traceID := extractor.GetTraceID(); traceID != "trace-999" {
		t.Errorf("Expected trace_id to be trace-999, got %s", traceID)
	}

	if requestID := extractor.GetRequestID(); requestID != "req-888" {
		t.Errorf("Expected request_id to be req-888, got %s", requestID)
	}

	if userID := extractor.GetUserID(); userID != "user-777" {
		t.Errorf("Expected user_id to be user-777, got %s", userID)
	}

	allInfo := extractor.GetAllTraceInfo()
	if len(allInfo) != 3 {
		t.Errorf("Expected 3 trace info items, got %d", len(allInfo))
	}
}
