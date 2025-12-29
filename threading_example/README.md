# go-zero threading 包示例

本目录包含 go-zero `core/threading` 包的完整示例代码和详细说明。

## 📚 文档

- **[THREADING_GUIDE.md](./THREADING_GUIDE.md)** - threading 包完整功能介绍（必读）

## 📁 文件说明

### 核心文件
- `examples_test.go` - 完整的测试用例（✅ 推荐使用）
- `routine_group.go` - RoutineGroup 示例代码
- `safe_run.go` - 基础安全执行示例

### 文档文件
- `README.md` - 本文件，项目说明
- `THREADING_GUIDE.md` - threading 包完整功能介绍
- `THREADING_ANALYSIS.md` - threading 包详细分析
- `TEST_EXECUTION_GUIDE.md` - 测试执行指南

## 🚀 快速开始

### 运行所有测试（推荐）

```bash
cd threading_example

# 运行所有测试
go test -v

# 运行完整功能测试（对应原 main.go 的功能）
go test -v -run TestAllThreadingFeatures
```

### 运行特定功能的测试

```bash
# 运行 GoSafe 简单测试
go test -v -run TestGoSafeSimple

# 运行 RunSafe 简单测试
go test -v -run TestRunSafeSimple

# 运行 RoutineGroup 简单测试
go test -v -run TestRoutineGroupSimple

# 运行单个具体测试
go test -v -run TestGoSafeSimple/AsyncTask
go test -v -run TestRunSafeSimple/PanicRecovery
```

### 查看测试覆盖率

```bash
# 生成覆盖率报告
go test -cover

# 生成详细覆盖率报告
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 在代码中使用示例

```go
package main

import "go-examples/threading_example"

func main() {
    // 运行 GoSafe 示例
    threading_example.GoSafeBasicDemo()
    
    // 运行 RunSafe 示例
    threading_example.RunSafeBasicDemo()
    
    // 运行 RoutineGroup 示例
    threading_example.RoutineGroupDemo()
}
```

## 📖 核心功能概览

### 1. GoSafe - 安全启动协程
```go
threading.GoSafe(func() {
    // 异步任务，自动捕获 panic
    sendEmail(user.Email)
})
```

**适用场景**：
- ✅ 发送邮件、短信等异步通知
- ✅ 记录日志
- ✅ 更新缓存
- ✅ 调用第三方 API

### 2. RunSafe - 当前协程安全执行
```go
threading.RunSafe(func() {
    // 在当前协程执行，自动捕获 panic
    processRequest(r)
})
```

**适用场景**：
- ✅ HTTP Handler 防护
- ✅ 中间件防护
- ✅ 定时任务防护
- ✅ 回调函数防护

### 3. RoutineGroup - 协程组管理
```go
group := threading.NewRoutineGroup()

for _, user := range users {
    u := user
    group.RunSafe(func() {
        processUser(u)
    })
}

group.Wait() // 等待所有协程完成
```

**适用场景**：
- ✅ 批量并发处理
- ✅ 微服务聚合调用
- ✅ 数据导入/导出

### 4. WorkerPool - 协程池（限制并发数）
```go
pool := threading.NewWorkerPool(10) // 限制并发数为 10

for _, task := range tasks {
    t := task
    pool.Schedule(func() {
        processTask(t)
    })
}

pool.Wait() // 等待所有任务完成
```

**适用场景**：
- ✅ 批量发送通知（限流保护）
- ✅ 图片处理（CPU 密集型）
- ✅ 数据库批量操作
- ✅ 爬虫任务

### 5. TaskRunner - 非阻塞任务调度器
```go
runner := threading.NewTaskRunner(10)

// 非阻塞提交任务
ok := runner.Schedule(func() {
    processTask()
})

if !ok {
    // 队列满，拒绝任务
}
```

**适用场景**：
- ✅ 长期运行的服务
- ✅ 消息队列消费
- ✅ 实时日志处理
- ✅ HTTP 接口接收任务

## 🎯 选择指南

| 需求 | 推荐工具 |
|------|---------|
| 简单异步任务 | GoSafe |
| 当前协程防护 | RunSafe |
| 等待一组任务完成 | RoutineGroup |
| 限制并发数 + 等待完成 | WorkerPool |
| 任务队列 + 背压控制 | TaskRunner |

## ⚠️ 常见陷阱

### 1. 闭包变量问题
```go
// ❌ 错误
for i := 0; i < 10; i++ {
    threading.GoSafe(func() {
        fmt.Println(i) // 可能全部打印 10
    })
}

// ✅ 正确
for i := 0; i < 10; i++ {
    index := i // 复制变量
    threading.GoSafe(func() {
        fmt.Println(index)
    })
}
```

### 2. WorkerPool 并发数设置
```go
// CPU 密集型任务
pool := threading.NewWorkerPool(runtime.NumCPU())

// IO 密集型任务
pool := threading.NewWorkerPool(runtime.NumCPU() * 2)

// 网络请求（根据下游承载能力）
pool := threading.NewWorkerPool(50)
```

### 3. TaskRunner 内存控制
```go
// 使用 Schedule 进行背压控制
if !runner.Schedule(task) {
    // 队列满，拒绝任务或等待
    return errors.New("系统繁忙")
}
```

## 📊 性能对比

| 工具 | 任务数 | 并发控制 | 内存占用 | 适用场景 |
|------|--------|---------|---------|---------|
| GoSafe | 少量 | ❌ | 低 | 简单异步 |
| RoutineGroup | < 1000 | ❌ | 中 | 批量处理 |
| WorkerPool | 大量 | ✅ | 中 | 限流场景 |
| TaskRunner | 持续 | ✅ | 高 | 长期服务 |

## 🔗 相关资源

- [go-zero 官方文档](https://go-zero.dev/)
- [go-zero GitHub](https://github.com/zeromicro/go-zero)
- [threading 源码](https://github.com/zeromicro/go-zero/tree/master/core/threading)

## 📝 示例统计

- **GoSafe**: 7 个示例
- **RunSafe**: 8 个示例
- **RoutineGroup**: 1 个示例
- **WorkerPool**: 4 个示例
- **TaskRunner**: 5 个示例

**总计**: 25+ 个实际应用场景示例

## 💡 最佳实践

1. **优先使用 RunSafe/GoSafe** - 简单场景不要过度设计
2. **合理设置并发数** - 根据任务类型和资源情况
3. **注意闭包变量** - 循环中要复制变量
4. **使用背压控制** - TaskRunner.Schedule 防止内存溢出
5. **监控和日志** - go-zero 会自动记录 panic 日志

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可

MIT License
