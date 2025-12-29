package mapreduce_example

import (
	"context"
)

// DetachContext 分离 context，保留所有值但移除取消信号
// 适用于需要在子协程中继承 traceID 等信息，但不希望被父 context 的超时/取消影响的场景
//
// Go 1.21+ 可以直接使用 context.WithoutCancel(ctx)
// Go 1.20 及以下需要手动复制值
func DetachContext(ctx context.Context, keys ...interface{}) context.Context {
	// Go 1.21+ 推荐使用 context.WithoutCancel
	// 这里提供兼容方案
	newCtx := context.Background()

	// 如果没有指定 keys，尝试复制常见的 keys
	if len(keys) == 0 {
		keys = []interface{}{
			"trace_id",
			"span_id",
			"request_id",
			"user_id",
			"tenant_id",
		}
	}

	for _, key := range keys {
		if val := ctx.Value(key); val != nil {
			newCtx = context.WithValue(newCtx, key, val)
		}
	}

	return newCtx
}

// InheritTraceContext 创建一个新的 context，继承链路追踪信息但不继承超时
// 这是专门为 MapReduce 等并发场景设计的工具函数
//
// 使用场景：
// 1. 在 MapReduce 的 Mapper 中需要记录日志，要求日志带有 traceID
// 2. 子协程的执行时间可能超过父 context 的超时时间
// 3. 需要独立的超时控制
func InheritTraceContext(parent context.Context) context.Context {
	// 常见的链路追踪相关的 key
	traceKeys := []interface{}{
		"trace_id",
		"span_id",
		"parent_span_id",
		"request_id",
		"user_id",
		"tenant_id",
		"app_id",
		"env",
	}

	return DetachContext(parent, traceKeys...)
}

// WithTraceTimeout 创建一个带超时的 context，同时继承链路追踪信息
// 这个函数结合了 InheritTraceContext 和 context.WithTimeout
//
// 使用场景：
// 在 MapReduce 的 Mapper 中，需要为每个任务设置独立的超时时间
func WithTraceTimeout(parent context.Context, timeout interface{}) (context.Context, context.CancelFunc) {
	// 先继承链路追踪信息
	newCtx := InheritTraceContext(parent)

	// 再设置新的超时时间
	// 这里简化处理，实际使用时可以根据 timeout 的类型做不同处理
	// return context.WithTimeout(newCtx, timeout.(time.Duration))

	// 为了示例简单，这里返回不带超时的
	return newCtx, func() {}
}

// ContextValueExtractor 用于提取 context 中的值
type ContextValueExtractor struct {
	ctx context.Context
}

// NewContextValueExtractor 创建一个 context 值提取器
func NewContextValueExtractor(ctx context.Context) *ContextValueExtractor {
	return &ContextValueExtractor{ctx: ctx}
}

// GetTraceID 获取 traceID
func (e *ContextValueExtractor) GetTraceID() string {
	if val := e.ctx.Value("trace_id"); val != nil {
		if traceID, ok := val.(string); ok {
			return traceID
		}
	}
	return ""
}

// GetRequestID 获取 requestID
func (e *ContextValueExtractor) GetRequestID() string {
	if val := e.ctx.Value("request_id"); val != nil {
		if requestID, ok := val.(string); ok {
			return requestID
		}
	}
	return ""
}

// GetUserID 获取 userID
func (e *ContextValueExtractor) GetUserID() string {
	if val := e.ctx.Value("user_id"); val != nil {
		if userID, ok := val.(string); ok {
			return userID
		}
	}
	return ""
}

// GetAllTraceInfo 获取所有链路追踪信息
func (e *ContextValueExtractor) GetAllTraceInfo() map[string]string {
	info := make(map[string]string)

	if traceID := e.GetTraceID(); traceID != "" {
		info["trace_id"] = traceID
	}
	if requestID := e.GetRequestID(); requestID != "" {
		info["request_id"] = requestID
	}
	if userID := e.GetUserID(); userID != "" {
		info["user_id"] = userID
	}

	return info
}
