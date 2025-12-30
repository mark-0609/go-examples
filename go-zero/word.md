# Go-Zero Core 包函数应用场景完整文档

> **文档版本**: v1.0  
> **生成时间**: 2025-12-30  
> **适用版本**: go-zero latest  
> **文档说明**: 本文档详细介绍 go-zero/core 目录下所有包的函数及其应用场景

---

## 目录

- [1. bloom - 布隆过滤器](#1-bloom---布隆过滤器)
- [2. breaker - 熔断器](#2-breaker---熔断器)
- [3. cmdline - 命令行工具](#3-cmdline---命令行工具)
- [4. codec - 编解码](#4-codec---编解码)
- [5. collection - 集合工具](#5-collection---集合工具)
- [6. color - 终端颜色](#6-color---终端颜色)
- [7. conf - 配置加载](#7-conf---配置加载)
- [8. configcenter - 配置中心](#8-configcenter---配置中心)
- [9. contextx - Context扩展](#9-contextx---context扩展)
- [10. discov - 服务发现](#10-discov---服务发现)
- [11. errorx - 错误处理](#11-errorx---错误处理)
- [12. executors - 执行器](#12-executors---执行器)
- [13. filex - 文件扩展](#13-filex---文件扩展)
- [14. fs - 文件系统](#14-fs---文件系统)
- [15. fx - 函数式编程](#15-fx---函数式编程)
- [16. hash - 哈希算法](#16-hash---哈希算法)
- [17. iox - IO扩展](#17-iox---io扩展)
- [18. jsonx - JSON扩展](#18-jsonx---json扩展)
- [19. lang - 语言工具](#19-lang---语言工具)
- [20. limit - 限流器](#20-limit---限流器)
- [21. load - 负载统计](#21-load---负载统计)
- [22. logc - 日志上下文](#22-logc---日志上下文)
- [23. logx - 日志系统](#23-logx---日志系统)
- [24. mapping - 映射工具](#24-mapping---映射工具)
- [25. mathx - 数学工具](#25-mathx---数学工具)
- [26. metric - 指标监控](#26-metric---指标监控)
- [27. mr - MapReduce](#27-mr---mapreduce)
- [28. naming - 命名工具](#28-naming---命名工具)
- [29. netx - 网络工具](#29-netx---网络工具)
- [30. proc - 进程管理](#30-proc---进程管理)
- [31. prof - 性能分析](#31-prof---性能分析)
- [32. prometheus - Prometheus集成](#32-prometheus---prometheus集成)
- [33. queue - 队列](#33-queue---队列)
- [34. rescue - 异常恢复](#34-rescue---异常恢复)
- [35. search - 搜索工具](#35-search---搜索工具)
- [36. service - 服务框架](#36-service---服务框架)
- [37. stat - 统计工具](#37-stat---统计工具)
- [38. stores - 存储](#38-stores---存储)
- [39. stringx - 字符串工具](#39-stringx---字符串工具)
- [40. syncx - 同步工具](#40-syncx---同步工具)
- [41. sysx - 系统工具](#41-sysx---系统工具)
- [42. threading - 并发工具](#42-threading---并发工具)
- [43. timex - 时间工具](#43-timex---时间工具)
- [44. trace - 链路追踪](#44-trace---链路追踪)
- [45. utils - 通用工具](#45-utils---通用工具)
- [46. validation - 数据验证](#46-validation---数据验证)

---

## 1. bloom - 布隆过滤器

### 包说明
布隆过滤器是一种空间效率极高的概率型数据结构，用于判断一个元素是否在集合中。

### 核心类型

#### **Filter**
布隆过滤器实现，基于 Redis 的位图操作。

### 主要函数

#### **New(store *redis.Redis, key string, bits uint) *Filter**
- **作用**: 创建布隆过滤器实例
- **参数**:
    - `store`: Redis 客户端
    - `key`: Redis 中的键名
    - `bits`: 位图大小
- **应用场景**:
  ```go
  // 场景1: 防止缓存穿透
  filter := bloom.New(rds, "user:bloom", 1024*1024)
  
  // 场景2: 去重检查
  filter := bloom.New(rds, "email:bloom", 10000000)
  ```

#### **Add(data []byte) error**
- **作用**: 添加元素到布隆过滤器
- **应用场景**:
  ```go
  // 场景1: 添加已注册邮箱
  filter.Add([]byte("user@example.com"))
  
  // 场景2: 添加已爬取URL
  filter.Add([]byte("https://example.com/page1"))
  ```

#### **Exists(data []byte) (bool, error)**
- **作用**: 检查元素是否可能存在
- **返回**: true表示可能存在，false表示一定不存在
- **应用场景**:
  ```go
  // 场景1: 检查邮箱是否已注册
  exists, _ := filter.Exists([]byte("user@example.com"))
  if exists {
      // 可能已注册，需要进一步查询数据库确认
  }
  
  // 场景2: 防止重复爬取
  exists, _ := filter.Exists([]byte(url))
  if !exists {
      // 一定未爬取，可以爬取
      crawl(url)
      filter.Add([]byte(url))
  }
  ```

### 典型应用场景

1. **防止缓存穿透**: 将数据库中的所有ID加入布隆过滤器，查询前先检查
2. **URL去重**: 爬虫系统中防止重复爬取
3. **垃圾邮件过滤**: 检查邮件地址是否在黑名单中
4. **推荐系统**: 过滤用户已看过的内容

---

## 2. breaker - 熔断器

### 包说明
实现熔断器模式，防止系统雪崩，提供自动降级和恢复能力。

### 核心类型

#### **Breaker**
熔断器接口，定义熔断器的基本行为。

### 主要函数

#### **NewBreaker(opts ...BreakerOption) Breaker**
- **作用**: 创建熔断器实例
- **选项**:
    - `WithName(name)`: 设置熔断器名称
    - `WithWindow(window)`: 设置统计窗口时间
- **应用场景**:
  ```go
  // 场景1: HTTP客户端熔断
  breaker := breaker.NewBreaker(
      breaker.WithName("api-client"),
  )
  
  // 场景2: 数据库连接熔断
  breaker := breaker.NewBreaker(
      breaker.WithName("mysql"),
  )
  ```

#### **Do(req func() error) error**
- **作用**: 在熔断器保护下执行请求
- **应用场景**:
  ```go
  // 场景1: 保护HTTP请求
  err := breaker.Do(func() error {
      resp, err := http.Get("https://api.example.com")
      return err
  })
  
  // 场景2: 保护数据库查询
  err := breaker.Do(func() error {
      return db.Query("SELECT * FROM users")
  })
  ```

#### **DoWithAcceptable(req func() error, acceptable Acceptable) error**
- **作用**: 执行请求，并自定义哪些错误是可接受的
- **应用场景**:
  ```go
  // 场景: 404错误不触发熔断
  err := breaker.DoWithAcceptable(
      func() error {
          return callAPI()
      },
      func(err error) bool {
          // 404不算失败
          return err == ErrNotFound
      },
  )
  ```

#### **DoWithFallback(req func() error, fallback func(err error) error) error**
- **作用**: 执行请求，失败时执行降级逻辑
- **应用场景**:
  ```go
  // 场景1: API降级到缓存
  err := breaker.DoWithFallback(
      func() error {
          return callAPI()
      },
      func(err error) error {
          // 降级：从缓存读取
          return getFromCache()
      },
  )
  
  // 场景2: 服务降级到默认值
  err := breaker.DoWithFallback(
      func() error {
          return getUserInfo(uid)
      },
      func(err error) error {
          // 返回默认用户信息
          return getDefaultUserInfo()
      },
  )
  ```

### 典型应用场景

1. **微服务调用保护**: 防止下游服务故障导致上游服务雪崩
2. **第三方API调用**: 保护系统不受第三方服务不稳定影响
3. **数据库访问保护**: 数据库故障时自动降级
4. **缓存降级**: 主服务不可用时降级到缓存

---

## 3. cmdline - 命令行工具

### 包说明
提供命令行参数解析和处理工具。

### 主要函数

#### **EnterToContinue()**
- **作用**: 等待用户按回车键继续
- **应用场景**:
  ```go
  // 场景: CLI工具中的交互式确认
  fmt.Println("准备删除所有数据，按回车继续...")
  cmdline.EnterToContinue()
  deleteAllData()
  ```

---

## 4. codec - 编解码

### 包说明
提供各种编解码功能，包括加密、解密、编码等。

### 主要函数

#### **EcbEncrypt(key, src []byte) ([]byte, error)**
- **作用**: ECB模式加密
- **应用场景**:
  ```go
  // 场景: 敏感数据加密
  encrypted, err := codec.EcbEncrypt(key, []byte("sensitive data"))
  ```

#### **EcbDecrypt(key, src []byte) ([]byte, error)**
- **作用**: ECB模式解密
- **应用场景**:
  ```go
  // 场景: 解密敏感数据
  decrypted, err := codec.EcbDecrypt(key, encrypted)
  ```

#### **HmacSha256(key []byte, data string) []byte**
- **作用**: HMAC-SHA256签名
- **应用场景**:
  ```go
  // 场景1: API签名验证
  signature := codec.HmacSha256(secretKey, requestData)
  
  // 场景2: Webhook签名
  signature := codec.HmacSha256(webhookSecret, payload)
  ```

#### **Md5Hex(data []byte) string**
- **作用**: 计算MD5哈希值（十六进制）
- **应用场景**:
  ```go
  // 场景1: 文件完整性校验
  hash := codec.Md5Hex(fileContent)
  
  // 场景2: 密码哈希（不推荐，仅示例）
  hash := codec.Md5Hex([]byte(password))
  ```

#### **RsaDecrypt(cipherText []byte, privateKey string) ([]byte, error)**
- **作用**: RSA解密
- **应用场景**:
  ```go
  // 场景: 解密客户端加密的敏感信息
  plaintext, err := codec.RsaDecrypt(encrypted, privateKey)
  ```

#### **RsaEncrypt(plainText []byte, publicKey string) ([]byte, error)**
- **作用**: RSA加密
- **应用场景**:
  ```go
  // 场景: 加密传输密码
  encrypted, err := codec.RsaEncrypt([]byte(password), publicKey)
  ```

### 典型应用场景

1. **API签名**: 使用HMAC进行请求签名和验证
2. **数据加密**: 使用RSA/AES加密敏感数据
3. **文件校验**: 使用MD5/SHA256校验文件完整性
4. **密码存储**: 使用哈希算法存储密码

---

## 5. collection - 集合工具

### 包说明
提供各种高性能的集合数据结构。

### 核心类型

#### **Cache**
LRU缓存实现。

#### **Ring**
环形缓冲区。

#### **Set**
集合实现。

#### **TimingWheel**
时间轮，用于延迟任务调度。

### 主要函数

#### **NewCache(expire time.Duration, opts ...CacheOption) (*Cache, error)**
- **作用**: 创建LRU缓存
- **应用场景**:
  ```go
  // 场景1: 用户信息缓存
  cache, _ := collection.NewCache(time.Hour)
  cache.Set("user:1001", userInfo)
  
  // 场景2: API响应缓存
  cache, _ := collection.NewCache(5*time.Minute)
  ```

#### **NewRing(n int) *Ring**
- **作用**: 创建环形缓冲区
- **应用场景**:
  ```go
  // 场景: 保存最近N条日志
  ring := collection.NewRing(100)
  ring.Add(logEntry)
  ```

#### **NewSet() *Set**
- **作用**: 创建集合
- **应用场景**:
  ```go
  // 场景: 去重
  set := collection.NewSet()
  set.Add("item1")
  set.Add("item2")
  if set.Contains("item1") {
      // ...
  }
  ```

#### **NewTimingWheel(interval time.Duration, numSlots int, execute Execute) (*TimingWheel, error)**
- **作用**: 创建时间轮
- **应用场景**:
  ```go
  // 场景1: 延迟任务
  tw, _ := collection.NewTimingWheel(time.Second, 60, func(key, value any) {
      // 执行延迟任务
  })
  tw.SetTimer("task1", task, 10*time.Second)
  
  // 场景2: 超时检测
  tw.SetTimer("conn:"+connID, conn, 30*time.Second)
  ```

### 典型应用场景

1. **LRU缓存**: 热点数据缓存
2. **环形缓冲**: 日志、指标数据存储
3. **集合操作**: 去重、交并差集
4. **延迟任务**: 订单超时取消、连接超时检测

---

## 6. color - 终端颜色

### 包说明
提供终端彩色输出功能。

### 主要函数

#### **WithColor(text string, colour color.Color) string**
- **作用**: 给文本添加颜色
- **应用场景**:
  ```go
  // 场景: CLI工具彩色输出
  fmt.Println(color.WithColor("Success", color.FgGreen))
  fmt.Println(color.WithColor("Error", color.FgRed))
  ```

---

## 7. conf - 配置加载

### 包说明
提供配置文件加载和解析功能，支持JSON、YAML、TOML等格式。

### 主要函数

#### **Load(file string, v any, opts ...Option) error**
- **作用**: 从文件加载配置
- **应用场景**:
  ```go
  // 场景1: 加载YAML配置
  var config Config
  conf.Load("config.yaml", &config)
  
  // 场景2: 加载JSON配置
  conf.Load("config.json", &config)
  ```

#### **LoadConfig(file string, v any, opts ...Option) error**
- **作用**: 加载配置（别名）
- **应用场景**: 同上

#### **LoadFromJsonBytes(content []byte, v any) error**
- **作用**: 从JSON字节加载配置
- **应用场景**:
  ```go
  // 场景: 从远程配置中心加载
  jsonData := fetchFromConfigCenter()
  conf.LoadFromJsonBytes(jsonData, &config)
  ```

#### **LoadFromYamlBytes(content []byte, v any) error**
- **作用**: 从YAML字节加载配置
- **应用场景**: 同上

#### **LoadFromTomlBytes(content []byte, v any) error**
- **作用**: 从TOML字节加载配置
- **应用场景**: 同上

#### **MustLoad(file string, v any, opts ...Option)**
- **作用**: 加载配置，失败则panic
- **应用场景**:
  ```go
  // 场景: 应用启动时加载必需配置
  var config Config
  conf.MustLoad("config.yaml", &config)
  ```

### 典型应用场景

1. **应用配置**: 加载数据库、Redis等配置
2. **环境配置**: 根据环境加载不同配置文件
3. **动态配置**: 从配置中心加载配置
4. **配置验证**: 加载时自动验证配置有效性

---

## 8. configcenter - 配置中心

### 包说明
提供配置中心集成，支持动态配置更新。

### 主要函数

#### **MustNewConfigCenter(config Config, ss Subscriber, opts ...Option) *Configurator**
- **作用**: 创建配置中心客户端
- **应用场景**:
  ```go
  // 场景: 集成Apollo/Nacos配置中心
  cc := configcenter.MustNewConfigCenter(config, subscriber)
  cc.AddListener(func() {
      // 配置变更回调
      reloadConfig()
  })
  ```

---

## 9. contextx - Context扩展

### 包说明
提供Context相关的扩展功能。

### 主要函数

#### **ValueOnlyFrom(ctx context.Context) context.Context**
- **作用**: 创建只保留值的Context（不继承取消信号）
- **应用场景**:
  ```go
  // 场景: 异步任务需要原Context的值但不受取消影响
  go func() {
      newCtx := contextx.ValueOnlyFrom(ctx)
      // 即使原ctx被取消，这里也能继续执行
      asyncTask(newCtx)
  }()
  ```

---

## 10. discov - 服务发现

### 包说明
基于etcd的服务注册与发现。

### 主要函数

#### **NewPublisher(endpoints []string, key, val string, opts ...PubOption) *Publisher**
- **作用**: 创建服务发布者（注册服务）
- **应用场景**:
  ```go
  // 场景: 微服务注册
  publisher := discov.NewPublisher(
      []string{"etcd:2379"},
      "services/user/192.168.1.100:8080",
      `{"host":"192.168.1.100","port":8080}`,
  )
  defer publisher.Stop()
  ```

#### **NewSubscriber(endpoints []string, key string, opts ...SubOption) (*Subscriber, error)**
- **作用**: 创建服务订阅者（发现服务）
- **应用场景**:
  ```go
  // 场景: 服务发现
  subscriber, _ := discov.NewSubscriber(
      []string{"etcd:2379"},
      "services/user",
  )
  subscriber.AddListener(func() {
      // 服务列表变更
      services := subscriber.Values()
      updateServiceList(services)
  })
  ```

### 典型应用场景

1. **微服务注册**: 服务启动时注册到etcd
2. **服务发现**: 动态发现可用服务实例
3. **负载均衡**: 基于服务列表进行负载均衡
4. **健康检查**: 自动剔除不健康的服务实例

---

## 11. errorx - 错误处理

### 包说明
提供错误处理相关的工具函数。

### 主要函数

#### **Wrap(err error, message string) error**
- **作用**: 包装错误并添加上下文信息
- **应用场景**:
  ```go
  // 场景: 添加错误上下文
  if err := db.Query(); err != nil {
      return errorx.Wrap(err, "failed to query database")
  }
  ```

#### **Wrapf(err error, format string, args ...any) error**
- **作用**: 格式化包装错误
- **应用场景**:
  ```go
  // 场景: 添加详细错误信息
  if err := processUser(uid); err != nil {
      return errorx.Wrapf(err, "failed to process user %d", uid)
  }
  ```

---

## 42. threading - 并发工具

### 包说明
提供并发编程的各种工具，包括协程组、任务执行器、稳定执行器等。

### 核心类型

#### **RoutineGroup**
协程组，用于管理和等待多个goroutine完成。

#### **TaskRunner**
任务执行器，控制并发数量。

#### **StableRunner**
稳定执行器，保证消息按推入顺序取出。

#### **WorkerGroup**
工作组，运行固定数量的worker处理相同任务。

### 主要函数

#### **NewRoutineGroup() *RoutineGroup**
- **作用**: 创建协程组
- **应用场景**:
  ```go
  // 场景1: 并发处理任务
  group := threading.NewRoutineGroup()
  for _, task := range tasks {
      task := task
      group.Run(func() {
          processTask(task)
      })
  }
  group.Wait()
  
  // 场景2: 并发HTTP请求
  group := threading.NewRoutineGroup()
  for _, url := range urls {
      url := url
      group.RunSafe(func() {
          fetchURL(url)
      })
  }
  group.Wait()
  ```

#### **Run(fn func())**
- **作用**: 在新goroutine中执行函数（不提供panic保护）
- **应用场景**:
  ```go
  // 场景: 可控代码的并发执行
  group.Run(func() {
      reliableFunction()
  })
  ```

#### **RunSafe(fn func())**
- **作用**: 在新goroutine中安全执行函数（自动捕获panic）
- **应用场景**:
  ```go
  // 场景: 不可控代码的并发执行
  group.RunSafe(func() {
      thirdPartyLib.DoSomething()
  })
  ```

#### **Wait()**
- **作用**: 等待所有goroutine完成
- **应用场景**: 见上述示例

#### **NewTaskRunner(concurrency int) *TaskRunner**
- **作用**: 创建任务执行器，限制并发数
- **应用场景**:
  ```go
  // 场景1: 限制HTTP请求并发数
  runner := threading.NewTaskRunner(10)
  for _, url := range urls {
      url := url
      runner.Schedule(func() {
          fetchURL(url)
      })
  }
  runner.Wait()
  
  // 场景2: 限制数据库操作并发
  runner := threading.NewTaskRunner(5)
  for _, record := range records {
      record := record
      runner.Schedule(func() {
          db.Insert(record)
      })
  }
  runner.Wait()
  ```

#### **Schedule(task func())**
- **作用**: 调度任务执行（阻塞式，并发满时等待）
- **应用场景**: 见上述示例

#### **ScheduleImmediately(task func()) error**
- **作用**: 立即调度任务（非阻塞，并发满时返回错误）
- **应用场景**:
  ```go
  // 场景: 需要快速失败的场景
  err := runner.ScheduleImmediately(func() {
      processTask()
  })
  if err == threading.ErrTaskRunnerBusy {
      // 系统繁忙，降级处理
      handleBusy()
  }
  ```

#### **NewStableRunner[I, O any](fn func(I) O) *StableRunner[I, O]**
- **作用**: 创建稳定执行器，保证按推入顺序输出
- **应用场景**:
  ```go
  // 场景1: Kafka消息处理
  runner := threading.NewStableRunner(func(msg KafkaMessage) Result {
      return processMessage(msg)
  })
  
  // 生产者
  go func() {
      for msg := range consumer.Messages() {
          runner.Push(msg)
      }
      runner.Wait()
  }()
  
  // 消费者（按顺序）
  for {
      result, err := runner.Get()
      if err != nil {
          break
      }
      saveToDatabase(result)
  }
  
  // 场景2: 并发数据转换，保持顺序
  runner := threading.NewStableRunner(func(data RawData) ProcessedData {
      return transform(data)
  })
  ```

#### **Push(v I) error**
- **作用**: 推入数据进行并发处理
- **应用场景**: 见上述示例

#### **Get() (O, error)**
- **作用**: 按推入顺序获取处理结果
- **应用场景**: 见上述示例

#### **NewWorkerGroup(job func(), workers int) WorkerGroup**
- **作用**: 创建工作组
- **应用场景**:
  ```go
  // 场景1: 消息队列消费者
  wg := threading.NewWorkerGroup(func() {
      for msg := range msgQueue {
          processMessage(msg)
      }
  }, 10)
  wg.Start()
  
  // 场景2: 爬虫worker池
  wg := threading.NewWorkerGroup(func() {
      for url := range urlQueue {
          crawl(url)
      }
  }, 20)
  wg.Start()
  ```

#### **Start()**
- **作用**: 启动工作组
- **应用场景**: 见上述示例

#### **GoSafe(fn func())**
- **作用**: 安全启动goroutine（自动捕获panic）
- **应用场景**:
  ```go
  // 场景: 启动后台任务
  threading.GoSafe(func() {
      backgroundTask()
  })
  ```

#### **RunSafe(fn func())**
- **作用**: 安全执行函数（捕获panic）
- **应用场景**:
  ```go
  // 场景: 执行不可控代码
  threading.RunSafe(func() {
      thirdPartyLib.DoSomething()
  })
  ```

### 典型应用场景

1. **RoutineGroup**: 批量任务并发处理、并发HTTP请求
2. **TaskRunner**: 限流、控制并发数、资源保护
3. **StableRunner**: Kafka消费、保序处理、流式数据处理
4. **WorkerGroup**: 消息队列消费、爬虫、长期运行的worker池

---

## 17. iox - IO扩展

### 包说明
提供IO操作的增强工具，包括Buffer池、流复制、文本处理等。

### 主要函数

#### **NewBufferPool(capability int) *BufferPool**
- **作用**: 创建Buffer对象池
- **应用场景**:
  ```go
  // 场景1: HTTP请求处理
  var bufPool = iox.NewBufferPool(4096)
  buf := bufPool.Get()
  defer bufPool.Put(buf)
  io.Copy(buf, r.Body)
  
  // 场景2: JSON序列化
  buf := bufPool.Get()
  defer bufPool.Put(buf)
  json.NewEncoder(buf).Encode(data)
  ```

#### **NopCloser(w io.Writer) io.WriteCloser**
- **作用**: 将Writer包装成WriteCloser（Close为空操作）
- **应用场景**:
  ```go
  // 场景: 适配接口
  var buf bytes.Buffer
  writer := iox.NopCloser(&buf)
  defer writer.Close()  // 不会真正关闭
  ```

#### **DupReadCloser(reader io.ReadCloser) (io.ReadCloser, io.ReadCloser)**
- **作用**: 复制ReadCloser，返回两个独立的Reader
- **应用场景**:
  ```go
  // 场景: HTTP请求体多次读取
  reader1, reader2 := iox.DupReadCloser(r.Body)
  defer reader1.Close()
  defer reader2.Close()
  
  // 第一次：记录日志
  body1, _ := io.ReadAll(reader1)
  log.Printf("Request: %s", body1)
  
  // 第二次：业务处理
  body2, _ := io.ReadAll(reader2)
  processData(body2)
  ```

#### **LimitDupReadCloser(reader io.ReadCloser, n int64) (io.ReadCloser, io.ReadCloser)**
- **作用**: 复制ReadCloser，第二个Reader限制读取n字节
- **应用场景**:
  ```go
  // 场景: 大文件日志记录（只记录前1KB）
  fullReader, previewReader := iox.LimitDupReadCloser(file, 1024)
  defer fullReader.Close()
  defer previewReader.Close()
  
  preview, _ := io.ReadAll(previewReader)
  log.Printf("Preview: %s", preview)
  
  fullData, _ := io.ReadAll(fullReader)
  processData(fullData)
  ```

#### **ReadBytes(reader io.Reader, buf []byte) error**
- **作用**: 精确读取指定长度的字节
- **应用场景**:
  ```go
  // 场景1: 协议头解析
  headerBuf := make([]byte, 16)
  iox.ReadBytes(conn, headerBuf)
  
  // 场景2: 二进制文件读取
  recordBuf := make([]byte, 128)
  iox.ReadBytes(file, recordBuf)
  ```

#### **ReadText(filename string) (string, error)**
- **作用**: 读取文件内容并去除首尾空格
- **应用场景**:
  ```go
  // 场景1: 读取Token
  token, _ := iox.ReadText("/etc/secrets/api_token")
  
  // 场景2: 读取版本号
  version, _ := iox.ReadText("VERSION")
  ```

#### **ReadTextLines(filename string, opts ...TextReadOption) ([]string, error)**
- **作用**: 按行读取文本文件
- **选项**:
    - `KeepSpace()`: 保留首尾空格
    - `WithoutBlank()`: 忽略空行
    - `OmitWithPrefix(prefix)`: 忽略指定前缀的行
- **应用场景**:
  ```go
  // 场景1: 读取配置文件（忽略注释）
  lines, _ := iox.ReadTextLines("config.txt",
      iox.WithoutBlank(),
      iox.OmitWithPrefix("#"),
  )
  
  // 场景2: 读取主机列表
  hosts, _ := iox.ReadTextLines("/etc/hosts",
      iox.WithoutBlank(),
      iox.OmitWithPrefix("#"),
  )
  ```

#### **LimitTeeReader(r io.Reader, w io.Writer, n int64) io.Reader**
- **作用**: 类似TeeReader，但限制写入字节数
- **应用场景**:
  ```go
  // 场景: 大文件日志记录（只记录前N字节）
  var logBuf bytes.Buffer
  limitedReader := iox.LimitTeeReader(file, &logBuf, 1024)
  
  data, _ := io.ReadAll(limitedReader)
  log.Printf("Preview: %s", logBuf.String())
  processData(data)
  ```

#### **CountLines(file string) (int, error)**
- **作用**: 统计文件行数
- **应用场景**:
  ```go
  // 场景1: 日志文件统计
  lines, _ := iox.CountLines("app.log")
  fmt.Printf("日志共 %d 行\n", lines)
  
  // 场景2: 进度显示
  totalLines, _ := iox.CountLines("data.csv")
  fmt.Printf("总共需要处理 %d 行\n", totalLines)
  ```

#### **NewTextLineScanner(reader io.Reader) *TextLineScanner**
- **作用**: 创建文本行扫描器
- **应用场景**:
  ```go
  // 场景1: 逐行处理日志
  scanner := iox.NewTextLineScanner(file)
  for scanner.Scan() {
      line, _ := scanner.Line()
      if strings.Contains(line, "ERROR") {
          handleError(line)
      }
  }
  
  // 场景2: 流式处理HTTP响应
  scanner := iox.NewTextLineScanner(resp.Body)
  for scanner.Scan() {
      line, _ := scanner.Line()
      processLine(line)
  }
  ```

### 典型应用场景

1. **BufferPool**: HTTP处理、JSON序列化、字符串拼接
2. **DupReadCloser**: 请求体多次读取、数据验证和处理
3. **ReadTextLines**: 配置文件读取、日志分析
4. **TextLineScanner**: 流式文本处理、大文件处理

---

## 43. timex - 时间工具

### 包说明
提供时间相关的工具函数。

### 主要函数

#### **Now() time.Time**
- **作用**: 获取当前时间（可mock）
- **应用场景**:
  ```go
  // 场景: 单元测试中mock时间
  now := timex.Now()
  ```

#### **Since(t time.Time) time.Duration**
- **作用**: 计算从t到现在的时间间隔
- **应用场景**:
  ```go
  // 场景: 性能统计
  start := timex.Now()
  doSomething()
  duration := timex.Since(start)
  ```

#### **Time() time.Duration**
- **作用**: 获取当前时间戳（纳秒）
- **应用场景**:
  ```go
  // 场景: 高精度计时
  start := timex.Time()
  process()
  elapsed := timex.Time() - start
  ```

---

## 总结

### 核心包分类

#### **基础工具类**
- `lang`: 语言基础工具
- `stringx`: 字符串处理
- `mathx`: 数学计算
- `timex`: 时间处理
- `hash`: 哈希算法

#### **并发编程类**
- `threading`: 并发工具（RoutineGroup、TaskRunner、StableRunner）
- `syncx`: 同步工具
- `executors`: 执行器

#### **IO处理类**
- `iox`: IO扩展
- `fs`: 文件系统
- `filex`: 文件扩展

#### **网络通信类**
- `netx`: 网络工具
- `discov`: 服务发现
- `naming`: 命名服务

#### **数据存储类**
- `stores`: 存储（Redis、SQL、MongoDB等）
- `collection`: 集合数据结构
- `bloom`: 布隆过滤器

#### **可靠性保障类**
- `breaker`: 熔断器
- `limit`: 限流器
- `rescue`: 异常恢复
- `errorx`: 错误处理

#### **监控观测类**
- `logx`: 日志系统
- `metric`: 指标监控
- `trace`: 链路追踪
- `stat`: 统计工具
- `prof`: 性能分析

#### **配置管理类**
- `conf`: 配置加载
- `configcenter`: 配置中心

#### **编解码类**
- `codec`: 编解码
- `jsonx`: JSON扩展
- `mapping`: 映射工具

#### **服务框架类**
- `service`: 服务框架
- `proc`: 进程管理
- `queue`: 队列

### 使用建议

1. **并发处理**: 优先使用 `threading` 包的工具
2. **限流熔断**: 使用 `limit` 和 `breaker` 保护系统
3. **日志监控**: 使用 `logx`、`metric`、`trace` 构建可观测性
4. **配置管理**: 使用 `conf` 和 `configcenter` 管理配置
5. **数据存储**: 使用 `stores` 包的统一接口
6. **IO操作**: 使用 `iox` 提高IO效率

---

**文档结束**

> 本文档持续更新中，如有疑问或建议，欢迎反馈。
