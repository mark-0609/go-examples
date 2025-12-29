# MapReduce Context 传递最佳实践

## 问题描述

在使用 go-zero 的 MapReduce 时，经常遇到以下矛盾：
- **需求1**：子协程需要继承父 context 的 traceID，以便日志可以通过同一个 traceID 查询
- **需求2**：子协程的执行时间可能超过父 context 的超时时间，不希望被父 context 的超时影响

如果直接传递带超时的 context，子协程会因为超时而被取消；如果不传递 context，又无法继承 traceID。

## 解决方案

### 方案1：使用 MapReduce + WithContext 选项（推荐）⭐

go-zero 提供了 `WithContext` 选项，配合 `MapReduce` 函数使用，可以自动处理 context 的传递。

```go
ctx := context.WithValue(context.Background(), "trace_id", "trace-12345")
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

result, err := mr.MapReduce(
    func(source chan<- int) {
        // Generator 可以访问 ctx
        traceID := ctx.Value("trace_id")
        for _, v := range nums {
            source <- v
        }
    },
    func(item int, writer mr.Writer[int], cancel func(error)) {
        // Mapper 中的 ctx 已经被处理过
        // 会继承 traceID，但不会因为父 context 超时而立即取消
        logx.WithContext(ctx).Infof("Processing: %d", item)
        writer.Write(item * item)
    },
    func(pipe <-chan int, writer mr.Writer[[]int], cancel func(error)) {
        // Reducer 聚合结果
        var results []int
        for v := range pipe {
            results = append(results, v)
        }
        writer.Write(results)
    },
    mr.WithContext(ctx), // 关键：通过选项传递 context
)
```

**优点**：
- ✅ 自动处理 context 传递
- ✅ 子协程继承 traceID
- ✅ 不受父 context 超时影响
- ✅ 代码简洁

**适用场景**：
- 所有使用 MapReduce 的场景（强烈推荐）

---

### 方案2：使用 context.WithoutCancel（Go 1.21+）

Go 1.21 引入了 `context.WithoutCancel` 函数，可以创建一个继承所有值但不继承取消信号的 context。

```go
ctx := context.WithValue(context.Background(), "trace_id", "trace-12345")
ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()

// 分离 context：继承值，但不继承取消信号
detachedCtx := context.WithoutCancel(ctx)

result, err := mr.MapReduce(
    func(source chan<- int) {
        for _, v := range nums {
            source <- v
        }
    },
    func(item int, writer mr.Writer[int], cancel func(error)) {
        // 使用分离后的 context
        logx.WithContext(detachedCtx).Infof("Processing: %d", item)
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
```

**优点**：
- ✅ 官方标准库支持
- ✅ 语义清晰
- ✅ 性能好

**缺点**：
- ❌ 需要 Go 1.21+

**适用场景**：
- Go 版本 >= 1.21
- 需要手动控制 context 分离的场景

---

### 方案3：手动复制 context 值（Go 1.20 及以下）

对于 Go 1.20 及以下版本，可以手动复制 context 中的值到新的 context。

```go
// 工具函数
func DetachContext(ctx context.Context, keys ...interface{}) context.Context {
    newCtx := context.Background()
    for _, key := range keys {
        if val := ctx.Value(key); val != nil {
            newCtx = context.WithValue(newCtx, key, val)
        }
    }
    return newCtx
}

// 使用
ctx := context.WithValue(context.Background(), "trace_id", "trace-12345")
ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()

// 手动复制值
newCtx := DetachContext(ctx, "trace_id", "user_id")

result, err := mr.MapReduce(
    func(source chan<- int) {
        for _, v := range nums {
            source <- v
        }
    },
    func(item int, writer mr.Writer[int], cancel func(error)) {
        logx.WithContext(newCtx).Infof("Processing: %d", item)
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
```

**优点**：
- ✅ 兼容所有 Go 版本
- ✅ 灵活可控

**缺点**：
- ❌ 需要手动指定要复制的 key
- ❌ 代码略显冗余

**适用场景**：
- Go 版本 < 1.21
- 需要精确控制哪些值被复制

---

## 真实场景示例

### 场景：批量查询用户信息

```go
func BatchQueryUsers(ctx context.Context, userIDs []int64) ([]*UserInfo, error) {
    // ctx 来自 HTTP 请求，带有 traceID 和 3 秒超时
    
    return mr.MapReduce(
        func(source chan<- int64) {
            for _, uid := range userIDs {
                select {
                case <-ctx.Done():
                    // 检查是否超时
                    return
                case source <- uid:
                }
            }
        },
        func(uid int64, writer mr.Writer[*UserInfo], cancel func(error)) {
            // 查询单个用户（可能需要 500ms）
            user, err := queryUserFromDB(ctx, uid)
            if err != nil {
                logx.WithContext(ctx).Errorf("Query user %d failed: %v", uid, err)
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
        mr.WithContext(ctx), // 关键：通过选项传递 context
    )
}

func queryUserFromDB(ctx context.Context, uid int64) (*UserInfo, error) {
    // 日志会自动带上 traceID
    logx.WithContext(ctx).Infof("Querying user: %d", uid)
    
    // 模拟数据库查询
    time.Sleep(500 * time.Millisecond)
    
    return &UserInfo{ID: uid, Name: fmt.Sprintf("User-%d", uid)}, nil
}
```

### 关键点说明

1. **使用 MapReduce + WithContext 选项**：自动处理 context 传递
2. **在 Generator 中检查超时**：避免继续生成任务
3. **在 Mapper 中使用 logx.WithContext**：日志自动带上 traceID
4. **错误处理**：单个任务失败不影响其他任务

---

## 常见问题

### Q1: 为什么不能直接传递带超时的 context？

**A**: 如果直接传递带超时的 context，当父 context 超时时，所有子协程都会被取消，即使它们还没有完成。这在批量处理场景中是不合理的。

### Q2: MapReduce + WithContext 内部是如何处理的？

**A**: go-zero 的 `MapReduce` 配合 `WithContext` 选项内部会：
1. 在 Generator 中使用原始 context（可以被超时取消）
2. 在 Mapper 和 Reducer 中使用分离后的 context（继承值但不继承取消信号）

### Q3: 如何在 Mapper 中设置独立的超时？

**A**: 可以在 Mapper 中创建新的带超时的 context：

```go
func(item int, writer mr.Writer[int], cancel func(error)) {
    // 为每个任务设置独立的 5 秒超时
    taskCtx, taskCancel := context.WithTimeout(detachedCtx, 5*time.Second)
    defer taskCancel()
    
    result, err := doSomethingWithTimeout(taskCtx, item)
    if err != nil {
        logx.WithContext(taskCtx).Errorf("Task failed: %v", err)
        return
    }
    writer.Write(result)
}
```

### Q4: 如何确保日志可以通过 traceID 查询？

**A**: 使用 `logx.WithContext(ctx)` 记录日志，go-zero 会自动从 context 中提取 traceID 并添加到日志中。

```go
// 正确方式
logx.WithContext(ctx).Infof("Processing item: %d", item)

// 错误方式（不会带 traceID）
logx.Infof("Processing item: %d", item)
```

---

## 最佳实践总结

### ✅ 推荐做法

1. **优先使用 MapReduce + WithContext 选项**
   - 这是最简单、最安全的方式
   - go-zero 已经帮你处理好了 context 传递

2. **使用 logx.WithContext 记录日志**
   - 确保日志带上 traceID
   - 方便问题排查

3. **在 Generator 中检查 context 是否取消**
   - 避免在父 context 超时后继续生成任务

4. **合理处理错误**
   - 单个任务失败不应该影响其他任务
   - 使用 logx 记录错误，而不是 cancel 整个任务

### ❌ 避免的做法

1. **不要直接传递带超时的 context 给 MapReduce**
   - 会导致子协程被提前取消

2. **不要在 Mapper 中使用 logx.Infof**
   - 应该使用 logx.WithContext(ctx).Infof

3. **不要忽略错误处理**
   - 应该记录错误日志，方便排查问题

4. **不要在循环中创建过多的 context**
   - 会增加 GC 压力

---

## 性能对比

| 方案 | 性能 | 兼容性 | 推荐度 |
|------|------|--------|--------|
| MapReduce + WithContext | ⭐⭐⭐⭐⭐ | Go 1.11+ | ⭐⭐⭐⭐⭐ |
| context.WithoutCancel | ⭐⭐⭐⭐⭐ | Go 1.21+ | ⭐⭐⭐⭐ |
| 手动复制 context 值 | ⭐⭐⭐⭐ | 所有版本 | ⭐⭐⭐ |

---

## 参考资料

- [go-zero MapReduce 文档](https://go-zero.dev/docs/tutorials/go-zero/mapreduce)
- [Go Context 官方文档](https://pkg.go.dev/context)
- [context.WithoutCancel (Go 1.21+)](https://pkg.go.dev/context#WithoutCancel)

---

## 示例代码

完整的示例代码请参考：
- `context_trace.go` - 各种方案的完整示例
- `context_utils.go` - 工具函数
- `context_trace_test.go` - 单元测试
