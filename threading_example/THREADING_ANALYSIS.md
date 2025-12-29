# go-zero core/threading 包完整功能分析

## 📦 包概述

基于 go-zero 源码分析，`core/threading` 包提供以下核心功能：

## 🔧 核心函数列表

### 1. GoSafe / GoSafeCtx - 安全启动协程

#### 函数签名
```go
func GoSafe(fn func())
func GoSafeCtx(ctx context.Context, fn func())
```

#### 功能说明
- 在新协程中安全执行函数
- 自动捕获并记录 panic
- `GoSafeCtx` 支持 context 传递

#### 应用场景

✅ **场景1：异步发送通知**
```go
// 用户注册后发送欢迎邮件
threading.GoSafe(func() {
    sendWelcomeEmail(user.Email)
})
```

✅ **场景2：异步记录日志**
```go
// 不阻塞主流程
threading.GoSafe(func() {
    logService.Record(operation)
})
```

✅ **场景3：异步更新缓存**
```go
// 数据更新后刷新缓存
threading.GoSafe(func() {
    cache.Refresh(key)
})
```

✅ **场景4：第三方API调用**
```go
// 上报统计数据，不影响主流程
threading.GoSafeCtx(ctx, func() {
    analytics.Report(event)
})
```

---

### 2. RunSafe / RunSafeCtx - 当前协程安全执行

#### 函数签名
```go
func RunSafe(fn func())
func RunSafeCtx(ctx context.Context, fn func())
```

#### 功能说明
- 在**当前协程**中安全执行函数
- 自动捕获并记录 panic
- 不创建新协程

#### 应用场景

✅ **场景1：HTTP Handler 防护**
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    threading.RunSafe(func() {
        result := processRequest(r)
        json.NewEncoder(w).Encode(result)
    })
}
```

✅ **场景2：中间件防护**
```go
func RecoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        threading.RunSafe(func() {
            next.ServeHTTP(w, r)
        })
    })
}
```

✅ **场景3：定时任务防护**
```go
func cronJob() {
    ticker := time.NewTicker(time.Minute)
    for range ticker.C {
        threading.RunSafe(func() {
            processScheduledTask()
        })
    }
}
```

---

### 3. RoutineGroup - 协程组管理

#### 函数签名
```go
func NewRoutineGroup() *RoutineGroup

type RoutineGroup struct {}

func (g *RoutineGroup) Run(fn func())
func (g *RoutineGroup) RunSafe(fn func())
func (g *RoutineGroup) Wait()
```

#### 功能说明
- 封装 `sync.WaitGroup`
- `Run`：启动协程
- `RunSafe`：启动协程并捕获 panic
- `Wait`：等待所有协程完成

#### 应用场景

✅ **场景1：批量并发处理**
```go
func BatchProcessUsers(userIDs []int) {
    group := threading.NewRoutineGroup()
    
    for _, uid := range userIDs {
        userID := uid
        group.RunSafe(func() {
            processUser(userID)
        })
    }
    
    group.Wait()
}
```

✅ **场景2：微服务聚合调用**
```go
func AggregateData(ctx context.Context, id string) (*Response, error) {
    group := threading.NewRoutineGroup()
    
    var userInfo *User
    var orderInfo *Order
    
    group.RunSafe(func() {
        userInfo = userService.GetUser(ctx, id)
    })
    
    group.RunSafe(func() {
        orderInfo = orderService.GetOrders(ctx, id)
    })
    
    group.Wait()
    
    return &Response{User: userInfo, Orders: orderInfo}, nil
}
```

✅ **场景3：数据导出**
```go
func ExportData(records []Record) {
    group := threading.NewRoutineGroup()
    
    batchSize := 100
    for i := 0; i < len(records); i += batchSize {
        batch := records[i:min(i+batchSize, len(records))]
        group.RunSafe(func() {
            exportBatch(batch)
        })
    }
    
    group.Wait()
}
```

---

### 4. TaskRunner - 任务调度器

#### 函数签名
```go
func NewTaskRunner(concurrency int) *TaskRunner

type TaskRunner struct {}

func (r *TaskRunner) Schedule(task func())
func (r *TaskRunner) ScheduleImmediately(task func()) error
func (r *TaskRunner) Wait()
```

#### 功能说明
- 控制协程并发数
- `Schedule`：提交任务（非阻塞，队列满时丢弃）
- `ScheduleImmediately`：立即提交任务（阻塞，返回错误如果 runner 已关闭）
- `Wait`：等待所有任务完成

#### 应用场景

✅ **场景1：限流处理**
```go
func ProcessRequests() {
    // 限制并发数为 10
    runner := threading.NewTaskRunner(10)
    
    for _, req := range requests {
        request := req
        runner.Schedule(func() {
            handleRequest(request)
        })
    }
    
    runner.Wait()
}
```

✅ **场景2：批量发送通知**
```go
func SendNotifications(users []User) {
    runner := threading.NewTaskRunner(20)
    
    for _, user := range users {
        u := user
        runner.Schedule(func() {
            sendNotification(u.ID)
        })
    }
    
    runner.Wait()
}
```

✅ **场景3：图片处理**
```go
func ProcessImages(images []string) {
    runner := threading.NewTaskRunner(runtime.NumCPU())
    
    for _, img := range images {
        imagePath := img
        runner.Schedule(func() {
            processImage(imagePath)
        })
    }
    
    runner.Wait()
}
```

✅ **场景4：HTTP 服务背压控制**
```go
func HandleTask(w http.ResponseWriter, r *http.Request) {
    err := runner.ScheduleImmediately(func() {
        processTask(r.Body)
    })
    
    if err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    
    w.Write([]byte("Task accepted"))
}
```

---

### 5. WorkerGroup - 工作组

#### 函数签名
```go
func NewWorkerGroup(job func(), workers int) WorkerGroup

type WorkerGroup struct {}

func (wg WorkerGroup) Start()
```

#### 功能说明
- 创建固定数量的 worker 执行相同任务
- 所有 worker 执行同一个 job 函数
- 适合长期运行的后台任务

#### 应用场景

✅ **场景1：消息队列消费**
```go
func ConsumeMessages(queue chan Message) {
    job := func() {
        for msg := range queue {
            processMessage(msg)
        }
    }
    
    // 启动 10 个 worker 消费消息
    wg := threading.NewWorkerGroup(job, 10)
    wg.Start()
}
```

✅ **场景2：日志处理**
```go
func ProcessLogs(logChan chan LogEntry) {
    job := func() {
        for log := range logChan {
            parseAndStore(log)
        }
    }
    
    wg := threading.NewWorkerGroup(job, 5)
    wg.Start()
}
```

✅ **场景3：数据同步**
```go
func SyncData(dataChan chan Data) {
    job := func() {
        for data := range dataChan {
            syncToDatabase(data)
        }
    }
    
    wg := threading.NewWorkerGroup(job, 3)
    wg.Start()
}
```

---

### 6. StableRunner - 稳定运行器

#### 函数签名
```go
func NewStableRunner[I, O any](fn func(I) O) *StableRunner[I, O]

type StableRunner[I, O any] struct {}

func (r *StableRunner[I, O]) Run(input I) O
```

#### 功能说明
- 泛型实现的稳定运行器
- 确保函数执行的稳定性
- 自动处理 panic

#### 应用场景

✅ **场景1：数据转换**
```go
// 创建一个稳定的数据转换器
converter := threading.NewStableRunner(func(data string) int {
    // 可能 panic 的转换逻辑
    return parseToInt(data)
})

result := converter.Run("123")
```

✅ **场景2：API 调用封装**
```go
apiCaller := threading.NewStableRunner(func(req Request) Response {
    // 可能失败的 API 调用
    return callExternalAPI(req)
})

response := apiCaller.Run(request)
```

---

### 7. RoutineId - 获取协程ID

#### 函数签名
```go
func RoutineId() uint64
```

#### 功能说明
- 获取当前协程的唯一 ID
- 用于调试和追踪

#### 应用场景

✅ **场景1：日志追踪**
```go
func ProcessTask() {
    routineID := threading.RoutineId()
    log.Printf("[Routine %d] Processing task...", routineID)
}
```

✅ **场景2：协程监控**
```go
func MonitorRoutine() {
    id := threading.RoutineId()
    metrics.RecordRoutine(id)
}
```

---

## 📊 功能对比表

| 功能 | GoSafe | RunSafe | RoutineGroup | TaskRunner | WorkerGroup | StableRunner |
|------|--------|---------|--------------|------------|-------------|--------------|
| 创建新协程 | ✅ | ❌ | ✅ | ✅ | ✅ | ❌ |
| 捕获 panic | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| 等待完成 | ❌ | ❌ | ✅ | ✅ | ❌ | ❌ |
| 限制并发数 | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ |
| Context 支持 | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| 泛型支持 | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |

## 🎯 选择指南

### 何时使用 GoSafe？
- ✅ 简单的异步任务（发邮件、记日志）
- ✅ 不需要等待结果
- ✅ 需要 context 传递

### 何时使用 RunSafe？
- ✅ 当前协程需要防护
- ✅ HTTP Handler、中间件
- ✅ 不需要创建新协程

### 何时使用 RoutineGroup？
- ✅ 需要等待一组任务完成
- ✅ 任务量不大（< 1000）
- ✅ 不需要限制并发数

### 何时使用 TaskRunner？
- ✅ 大量任务需要并发处理
- ✅ 需要限制并发数
- ✅ 需要等待所有任务完成
- ✅ 批量处理场景

### 何时使用 WorkerGroup？
- ✅ 长期运行的后台服务
- ✅ 固定数量的 worker
- ✅ 所有 worker 执行相同任务
- ✅ 消息队列消费场景

### 何时使用 StableRunner？
- ✅ 需要泛型支持
- ✅ 需要返回值
- ✅ 数据转换场景

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
    index := i
    threading.GoSafe(func() {
        fmt.Println(index)
    })
}
```

### 2. TaskRunner 并发数设置
```go
// CPU 密集型
runner := threading.NewTaskRunner(runtime.NumCPU())

// IO 密集型
runner := threading.NewTaskRunner(runtime.NumCPU() * 2)

// 网络请求（根据下游承载能力）
runner := threading.NewTaskRunner(50)
```

### 3. WorkerGroup 使用注意
```go
// ❌ 错误：job 函数会立即返回
wg := threading.NewWorkerGroup(func() {
    msg := <-queue
    process(msg)
}, 10)

// ✅ 正确：job 函数应该是循环
wg := threading.NewWorkerGroup(func() {
    for msg := range queue {
        process(msg)
    }
}, 10)
```

## 💡 最佳实践

1. **优先使用 RunSafe/GoSafe** - 简单场景不要过度设计
2. **合理设置并发数** - 根据任务类型和资源情况
3. **注意闭包变量** - 循环中要复制变量
4. **使用 Context** - GoSafeCtx/RunSafeCtx 支持超时和取消
5. **监控和日志** - go-zero 会自动记录 panic 日志

## 🔗 相关资源

- [go-zero 官方文档](https://go-zero.dev/)
- [go-zero GitHub](https://github.com/zeromicro/go-zero)
- [threading 源码](https://github.com/zeromicro/go-zero/tree/master/core/threading)

## 📝 总结

go-zero threading 包提供了 **7 个核心功能**：

1. **GoSafe/GoSafeCtx** - 安全启动协程
2. **RunSafe/RunSafeCtx** - 当前协程安全执行
3. **RoutineGroup** - 协程组管理
4. **TaskRunner** - 任务调度器（限制并发）
5. **WorkerGroup** - 工作组（固定 worker）
6. **StableRunner** - 稳定运行器（泛型）
7. **RoutineId** - 获取协程 ID

这些工具覆盖了从简单异步任务到复杂并发控制的各种场景，是构建高并发 Go 应用的利器！
