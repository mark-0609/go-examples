# go-zero core/threading 包完整功能介绍

## 📦 包概述

`core/threading` 是 go-zero 提供的协程管理工具包，主要解决以下问题：
- **协程安全启动**：自动捕获 panic，防止程序崩溃
- **协程组管理**：优雅地等待一组协程完成
- **协程池控制**：限制并发数，保护系统资源
- **任务调度**：非阻塞的任务队列处理

## 🔧 核心函数列表

### 1. GoSafe - 安全启动协程

#### 函数签名
```go
func GoSafe(fn func())
```

#### 功能说明
- 在新协程中安全执行函数
- 自动捕获并记录 panic，不会导致程序崩溃
- 内部使用 `logx` 记录 panic 信息

#### 应用场景

##### ✅ 场景1：后台异步任务
```go
// 发送邮件通知（不影响主流程）
threading.GoSafe(func() {
    sendEmail(user.Email, "欢迎注册")
})
```

##### ✅ 场景2：不稳定的第三方调用
```go
// 调用第三方 API，可能会 panic
threading.GoSafe(func() {
    thirdPartyAPI.Report(data)
})
```

##### ✅ 场景3：HTTP 服务中的后台处理
```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // 立即返回响应
    w.Write([]byte("OK"))
    
    // 后台异步处理
    threading.GoSafe(func() {
        processData(r.Body)
    })
}
```

#### 最佳实践
- ✅ 用于不需要等待结果的异步任务
- ✅ 用于可能 panic 的不稳定代码
- ❌ 不要用于需要获取返回值的场景
- ❌ 不要用于需要等待完成的场景

---

### 2. RunSafe - 当前协程安全执行

#### 函数签名
```go
func RunSafe(fn func())
```

#### 功能说明
- 在**当前协程**中安全执行函数
- 自动捕获并记录 panic
- 与 GoSafe 的区别：不创建新协程

#### 应用场景

##### ✅ 场景1：HTTP Handler 防护
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    threading.RunSafe(func() {
        // 业务逻辑可能 panic
        result := riskyBusinessLogic(r)
        json.NewEncoder(w).Encode(result)
    })
}
```

##### ✅ 场景2：中间件防护
```go
func RecoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        threading.RunSafe(func() {
            next.ServeHTTP(w, r)
        })
    })
}
```

##### ✅ 场景3：定时任务防护
```go
func cronJob() {
    ticker := time.NewTicker(time.Minute)
    for range ticker.C {
        threading.RunSafe(func() {
            // 定时任务逻辑
            processScheduledTask()
        })
    }
}
```

#### 最佳实践
- ✅ 用于需要在当前协程执行的场景
- ✅ 用于 HTTP Handler、中间件等需要防护的入口
- ❌ 不要嵌套使用（外层已经有 recover 就不需要内层再包一次）

---

### 3. RoutineGroup - 协程组管理

#### 函数签名
```go
func NewRoutineGroup() *RoutineGroup

type RoutineGroup struct {
    // 内部封装 sync.WaitGroup
}

func (g *RoutineGroup) Run(fn func())
func (g *RoutineGroup) RunSafe(fn func())
func (g *RoutineGroup) Wait()
```

#### 功能说明
- 封装 `sync.WaitGroup`，提供更友好的 API
- `Run`：启动协程并增加计数
- `RunSafe`：启动协程、增加计数、自动捕获 panic
- `Wait`：等待所有协程完成

#### 应用场景

##### ✅ 场景1：批量并发处理
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
    fmt.Println("所有用户处理完成")
}
```

##### ✅ 场景2：微服务聚合调用
```go
func AggregateData(ctx context.Context, id string) (*Response, error) {
    group := threading.NewRoutineGroup()
    
    var userInfo *User
    var orderInfo *Order
    var productInfo *Product
    
    // 并发调用多个服务
    group.RunSafe(func() {
        userInfo = userService.GetUser(ctx, id)
    })
    
    group.RunSafe(func() {
        orderInfo = orderService.GetOrders(ctx, id)
    })
    
    group.RunSafe(func() {
        productInfo = productService.GetProducts(ctx, id)
    })
    
    group.Wait()
    
    return &Response{
        User:    userInfo,
        Orders:  orderInfo,
        Products: productInfo,
    }, nil
}
```

##### ✅ 场景3：数据导入/导出
```go
func ExportData(records []Record) {
    group := threading.NewRoutineGroup()
    
    // 分批处理
    batchSize := 100
    for i := 0; i < len(records); i += batchSize {
        end := i + batchSize
        if end > len(records) {
            end = len(records)
        }
        
        batch := records[i:end]
        group.RunSafe(func() {
            exportBatch(batch)
        })
    }
    
    group.Wait()
}
```

#### 最佳实践
- ✅ 用于需要等待一组协程完成的场景
- ✅ 优先使用 `RunSafe` 而不是 `Run`
- ⚠️ 注意闭包变量问题（循环中要复制变量）
- ❌ 不适合大量协程（无并发数限制）

---

### 4. WorkerPool - 协程池（限制并发数）

#### 函数签名
```go
func NewWorkerPool(size int) *WorkerPool

type WorkerPool struct {
    // 内部使用 channel 控制并发数
}

func (p *WorkerPool) Schedule(task func())
func (p *WorkerPool) Wait()
```

#### 功能说明
- 创建固定大小的协程池
- 限制最大并发数，保护系统资源
- 自动捕获 panic

#### 应用场景

##### ✅ 场景1：限流保护下游服务
```go
func SendNotifications(users []User) {
    // 限制并发数为 10，避免打爆下游服务
    pool := threading.NewWorkerPool(10)
    
    for _, user := range users {
        u := user
        pool.Schedule(func() {
            notificationService.Send(u.ID, "系统通知")
        })
    }
    
    pool.Wait()
}
```

##### ✅ 场景2：资源密集型任务
```go
func ProcessImages(imagePaths []string) {
    // 限制并发数为 CPU 核心数
    pool := threading.NewWorkerPool(runtime.NumCPU())
    
    for _, path := range imagePaths {
        imagePath := path
        pool.Schedule(func() {
            // 图片处理是 CPU 密集型任务
            processImage(imagePath)
        })
    }
    
    pool.Wait()
}
```

##### ✅ 场景3：数据库批量操作
```go
func BatchInsert(records []Record) {
    // 限制并发数为 20，避免数据库连接池耗尽
    pool := threading.NewWorkerPool(20)
    
    for _, record := range records {
        r := record
        pool.Schedule(func() {
            db.Insert(r)
        })
    }
    
    pool.Wait()
}
```

##### ✅ 场景4：爬虫任务
```go
func CrawlWebsites(urls []string) {
    // 限制并发数为 50，避免被封 IP
    pool := threading.NewWorkerPool(50)
    
    for _, url := range urls {
        targetURL := url
        pool.Schedule(func() {
            content := fetchURL(targetURL)
            saveContent(content)
        })
    }
    
    pool.Wait()
}
```

#### 最佳实践
- ✅ 用于大量任务需要限制并发数的场景
- ✅ CPU 密集型任务：并发数 = CPU 核心数
- ✅ IO 密集型任务：并发数 = CPU 核心数 * 2
- ✅ 网络请求：根据下游服务承载能力设置
- ⚠️ 注意闭包变量问题

---

### 5. TaskRunner - 非阻塞任务调度器

#### 函数签名
```go
func NewTaskRunner(workers int) *TaskRunner

type TaskRunner struct {
    // 内部使用 channel 队列
}

func (r *TaskRunner) Schedule(task func()) bool
func (r *TaskRunner) ScheduleAuto(task func())
```

#### 功能说明
- 创建一个任务队列，后台持续消费
- `Schedule`：非阻塞提交任务，队列满时返回 false
- `ScheduleAuto`：自动扩容，任务一定会被执行

#### 应用场景

##### ✅ 场景1：长期运行的服务
```go
func StartService() {
    // 创建 10 个 worker 处理任务
    runner := threading.NewTaskRunner(10)
    
    // HTTP 接口接收任务
    http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
        ok := runner.Schedule(func() {
            processTask(r.Body)
        })
        
        if ok {
            w.Write([]byte("任务已提交"))
        } else {
            w.WriteHeader(http.StatusTooManyRequests)
            w.Write([]byte("系统繁忙，请稍后重试"))
        }
    })
}
```

##### ✅ 场景2：消息队列消费者
```go
func ConsumeMessages() {
    runner := threading.NewTaskRunner(20)
    
    for msg := range messageQueue {
        message := msg
        runner.ScheduleAuto(func() {
            handleMessage(message)
        })
    }
}
```

##### ✅ 场景3：实时日志处理
```go
func LogProcessor() {
    runner := threading.NewTaskRunner(5)
    
    for logEntry := range logChannel {
        entry := logEntry
        runner.ScheduleAuto(func() {
            parseAndStore(entry)
        })
    }
}
```

#### 最佳实践
- ✅ 用于长期运行的服务
- ✅ 用于需要队列缓冲的场景
- ✅ `Schedule` 用于需要背压控制的场景
- ✅ `ScheduleAuto` 用于不能丢失任务的场景
- ⚠️ 注意内存占用（队列可能积压）

---

## 📊 功能对比表

| 功能 | GoSafe | RunSafe | RoutineGroup | WorkerPool | TaskRunner |
|------|--------|---------|--------------|------------|------------|
| 创建新协程 | ✅ | ❌ | ✅ | ✅ | ✅ |
| 捕获 panic | ✅ | ✅ | ✅ | ✅ | ✅ |
| 等待完成 | ❌ | ❌ | ✅ | ✅ | ❌ |
| 限制并发数 | ❌ | ❌ | ❌ | ✅ | ✅ |
| 任务队列 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 背压控制 | ❌ | ❌ | ❌ | ❌ | ✅ |

## 🎯 选择指南

### 何时使用 GoSafe？
- ✅ 简单的异步任务（发邮件、记日志）
- ✅ 不需要等待结果
- ✅ 任务量不大

### 何时使用 RunSafe？
- ✅ 当前协程需要防护
- ✅ HTTP Handler、中间件
- ✅ 不需要创建新协程

### 何时使用 RoutineGroup？
- ✅ 需要等待一组任务完成
- ✅ 任务量不大（< 1000）
- ✅ 不需要限制并发数

### 何时使用 WorkerPool？
- ✅ 大量任务需要并发处理
- ✅ 需要限制并发数
- ✅ 需要等待所有任务完成
- ✅ 批量处理场景

### 何时使用 TaskRunner？
- ✅ 长期运行的服务
- ✅ 需要任务队列缓冲
- ✅ 需要背压控制
- ✅ 不需要等待任务完成

## ⚠️ 常见陷阱

### 1. 闭包变量问题
```go
// ❌ 错误示例
for i := 0; i < 10; i++ {
    threading.GoSafe(func() {
        fmt.Println(i)  // 可能全部打印 10
    })
}

// ✅ 正确示例
for i := 0; i < 10; i++ {
    index := i  // 复制变量
    threading.GoSafe(func() {
        fmt.Println(index)
    })
}
```

### 2. WorkerPool 并发数设置不当
```go
// ❌ 错误：并发数过大
pool := threading.NewWorkerPool(10000)  // 可能耗尽系统资源

// ✅ 正确：根据场景设置
pool := threading.NewWorkerPool(runtime.NumCPU())  // CPU 密集型
pool := threading.NewWorkerPool(runtime.NumCPU() * 2)  // IO 密集型
```

### 3. TaskRunner 内存泄漏
```go
// ❌ 错误：无限制提交任务
runner := threading.NewTaskRunner(10)
for {
    runner.ScheduleAuto(func() {
        // 如果消费速度 < 生产速度，内存会持续增长
    })
}

// ✅ 正确：使用 Schedule 进行背压控制
if !runner.Schedule(task) {
    // 队列满，拒绝任务或等待
}
```

## 📚 完整示例

查看以下文件获取完整示例代码：
- `gosafe_example.go` - GoSafe 示例
- `runsafe_example.go` - RunSafe 示例
- `routine_group_example.go` - RoutineGroup 示例
- `worker_pool_example.go` - WorkerPool 示例
- `task_runner_example.go` - TaskRunner 示例

## 🔗 相关资源

- [go-zero 官方文档](https://go-zero.dev/)
- [go-zero GitHub](https://github.com/zeromicro/go-zero)
- [threading 源码](https://github.com/zeromicro/go-zero/tree/master/core/threading)
