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

## 12. executors - 执行器

### 包说明
提供各种任务执行器，用于批量处理、延迟执行、定期执行等场景。

### 核心类型

#### **BulkExecutor**
批量执行器，当任务数达到阈值或时间间隔到达时批量执行。

#### **ChunkExecutor**
分块执行器，按数据大小分块执行。

#### **PeriodicalExecutor**
周期执行器，定期批量执行任务。

#### **DelayExecutor**
延迟执行器，延迟执行任务，多次触发只执行一次。

#### **LessExecutor**
限制执行器，在时间间隔内最多执行一次。

### 主要函数

#### **NewBulkExecutor(execute Execute, opts ...BulkOption) *BulkExecutor**
- **作用**: 创建批量执行器
- **选项**:
    - `WithBulkTasks(n)`: 设置批量大小
    - `WithBulkInterval(d)`: 设置刷新间隔
- **应用场景**:
  ```go
  // 场景1: 批量写入数据库
  executor := executors.NewBulkExecutor(func(items []any) {
      var records []Record
      for _, item := range items {
          records = append(records, item.(Record))
      }
      db.BatchInsert(records)
  }, executors.WithBulkTasks(100), executors.WithBulkInterval(time.Second))
  
  // 添加任务
  for _, record := range records {
      executor.Add(record)
  }
  executor.Wait()
  
  // 场景2: 批量发送消息
  executor := executors.NewBulkExecutor(func(items []any) {
      var messages []Message
      for _, item := range items {
          messages = append(messages, item.(Message))
      }
      kafka.SendBatch(messages)
  }, executors.WithBulkTasks(50))
  ```

#### **NewChunkExecutor(execute Execute, opts ...ChunkOption) *ChunkExecutor**
- **作用**: 创建分块执行器（按字节大小分块）
- **选项**:
    - `WithChunkBytes(n)`: 设置块大小（字节）
    - `WithFlushInterval(d)`: 设置刷新间隔
- **应用场景**:
  ```go
  // 场景: 批量上传文件（按大小分块）
  executor := executors.NewChunkExecutor(func(items []any) {
      var files []File
      for _, item := range items {
          files = append(files, item.(File))
      }
      uploadBatch(files)
  }, executors.WithChunkBytes(1024*1024)) // 1MB
  
  for _, file := range files {
      executor.Add(file, len(file.Content))
  }
  executor.Wait()
  ```

#### **NewDelayExecutor(fn func(), delay time.Duration) *DelayExecutor**
- **作用**: 创建延迟执行器
- **应用场景**:
  ```go
  // 场景1: 搜索框防抖
  executor := executors.NewDelayExecutor(func() {
      performSearch(keyword)
  }, 300*time.Millisecond)
  
  // 用户每次输入都触发，但只在停止输入300ms后执行
  onKeyPress := func() {
      executor.Trigger()
  }
  
  // 场景2: 配置文件变更延迟重载
  executor := executors.NewDelayExecutor(func() {
      reloadConfig()
  }, time.Second)
  
  fileWatcher.OnChange(func() {
      executor.Trigger() // 多次变更只重载一次
  })
  ```

#### **NewLessExecutor(threshold time.Duration) *LessExecutor**
- **作用**: 创建限制执行器（时间间隔内最多执行一次）
- **应用场景**:
  ```go
  // 场景1: 限制日志输出频率
  executor := executors.NewLessExecutor(time.Minute)
  
  for _, event := range events {
      executor.DoOrDiscard(func() {
          log.Printf("High frequency event occurred")
      })
  }
  
  // 场景2: 限制告警频率
  executor := executors.NewLessExecutor(5*time.Minute)
  
  if cpuUsage > 90 {
      executor.DoOrDiscard(func() {
          sendAlert("CPU usage too high")
      })
  }
  ```

#### **NewPeriodicalExecutor(interval time.Duration, container TaskContainer) *PeriodicalExecutor**
- **作用**: 创建周期执行器
- **应用场景**:
  ```go
  // 场景: 自定义批量处理逻辑
  container := &MyContainer{
      tasks: make([]Task, 0),
  }
  executor := executors.NewPeriodicalExecutor(time.Second, container)
  
  for _, task := range tasks {
      executor.Add(task)
  }
  executor.Wait()
  ```

### 典型应用场景

1. **BulkExecutor**: 批量数据库操作、批量消息发送
2. **ChunkExecutor**: 大文件分块上传、按大小批量处理
3. **DelayExecutor**: 搜索防抖、配置延迟重载
4. **LessExecutor**: 限制日志频率、限制告警频率
5. **PeriodicalExecutor**: 定期批量处理、定时任务

---

## 13. filex - 文件扩展

### 包说明
提供文件操作的扩展功能。

### 主要函数

#### **RangeReader(file *os.File, start, stop int64) io.ReadCloser**
- **作用**: 创建范围读取器，读取文件的指定范围
- **应用场景**:
  ```go
  // 场景1: 断点续传
  file, _ := os.Open("large_file.dat")
  reader := filex.RangeReader(file, 1024*1024, 2*1024*1024)
  io.Copy(conn, reader)
  
  // 场景2: 分片下载
  reader := filex.RangeReader(file, offset, offset+chunkSize)
  ```

---

## 14. fs - 文件系统

### 包说明
提供文件系统相关的工具函数。

### 主要函数

#### **TempFileWithText(text string) (string, error)**
- **作用**: 创建包含指定文本的临时文件
- **应用场景**:
  ```go
  // 场景: 单元测试中创建临时配置文件
  configFile, _ := fs.TempFileWithText(`
      host: localhost
      port: 8080
  `)
  defer os.Remove(configFile)
  
  conf.Load(configFile, &config)
  ```

#### **TempFilenameWithText(text string) (string, error)**
- **作用**: 创建临时文件并返回文件名
- **应用场景**: 同上

---

## 15. fx - 函数式编程

### 包说明
提供函数式编程工具，支持流式数据处理。

### 主要函数

#### **From(generate func(source chan<- any)) Stream**
- **作用**: 从生成函数创建流
- **应用场景**:
  ```go
  // 场景: 流式处理数据
  fx.From(func(source chan<- any) {
      for i := 0; i < 100; i++ {
          source <- i
      }
  }).Filter(func(item any) bool {
      return item.(int) % 2 == 0
  }).Map(func(item any) any {
      return item.(int) * 2
  }).ForEach(func(item any) {
      fmt.Println(item)
  })
  ```

#### **Just(items ...any) Stream**
- **作用**: 从元素创建流
- **应用场景**:
  ```go
  // 场景: 处理固定元素
  fx.Just(1, 2, 3, 4, 5).
      Filter(func(item any) bool {
          return item.(int) > 2
      }).
      ForEach(func(item any) {
          fmt.Println(item) // 输出: 3, 4, 5
      })
  ```

#### **Range(start, stop int) Stream**
- **作用**: 创建范围流
- **应用场景**:
  ```go
  // 场景: 批量处理
  fx.Range(1, 100).
      Map(func(item any) any {
          return processItem(item.(int))
      }).
      ForEach(func(item any) {
          saveResult(item)
      })
  ```

### Stream 方法

#### **Filter(fn FilterFunc) Stream**
- **作用**: 过滤元素
- **应用场景**: 数据筛选

#### **Map(fn MapFunc) Stream**
- **作用**: 转换元素
- **应用场景**: 数据转换

#### **Reduce(fn ReduceFunc) (any, error)**
- **作用**: 聚合元素
- **应用场景**:
  ```go
  // 场景: 求和
  sum, _ := fx.Range(1, 101).Reduce(func(a, b any) any {
      return a.(int) + b.(int)
  })
  fmt.Println(sum) // 5050
  ```

#### **ForEach(fn ForEachFunc)**
- **作用**: 遍历元素
- **应用场景**: 执行副作用操作

#### **Parallel(fn ParallelFunc, opts ...Option) Stream**
- **作用**: 并行处理元素
- **应用场景**:
  ```go
  // 场景: 并行HTTP请求
  fx.Just(urls...).
      Parallel(func(item any) any {
          return fetchURL(item.(string))
      }).
      ForEach(func(item any) {
          processResponse(item)
      })
  ```

### 典型应用场景

1. **数据转换**: 流式数据处理和转换
2. **并行处理**: 并行处理大量数据
3. **数据聚合**: 统计、求和、求平均值
4. **数据过滤**: 筛选符合条件的数据

---

## 16. hash - 哈希算法

### 包说明
提供一致性哈希算法实现。

### 主要函数

#### **NewConsistentHash() *ConsistentHash**
- **作用**: 创建一致性哈希实例
- **应用场景**:
  ```go
  // 场景1: 分布式缓存
  hash := hash.NewConsistentHash()
  hash.Add("cache-server-1")
  hash.Add("cache-server-2")
  hash.Add("cache-server-3")
  
  server, _ := hash.Get("user:1001")
  // 总是路由到同一台服务器
  
  // 场景2: 负载均衡
  hash := hash.NewConsistentHash()
  for _, server := range servers {
      hash.Add(server.Address)
  }
  
  targetServer, _ := hash.Get(requestID)
  ```

#### **Add(node string)**
- **作用**: 添加节点
- **应用场景**: 动态添加服务器节点

#### **Get(key string) (string, bool)**
- **作用**: 获取key对应的节点
- **应用场景**: 路由请求到指定节点

#### **Remove(node string)**
- **作用**: 移除节点
- **应用场景**: 服务器下线

### 典型应用场景

1. **分布式缓存**: Redis集群、Memcached集群
2. **负载均衡**: 请求路由、会话保持
3. **分布式存储**: 数据分片、副本分布

---

## 18. jsonx - JSON扩展

### 包说明
提供JSON处理的扩展功能。

### 主要函数

#### **Marshal(v any) ([]byte, error)**
- **作用**: JSON序列化（支持更多类型）
- **应用场景**:
  ```go
  // 场景: 序列化复杂对象
  data, _ := jsonx.Marshal(complexObject)
  ```

#### **Unmarshal(data []byte, v any) error**
- **作用**: JSON反序列化（更宽松的解析）
- **应用场景**:
  ```go
  // 场景: 解析JSON
  var obj Object
  jsonx.Unmarshal(data, &obj)
  ```

---

## 19. lang - 语言工具

### 包说明
提供Go语言相关的基础工具。

### 主要类型和常量

#### **PlaceholderType**
- **作用**: 空结构体类型，用于channel信号传递
- **应用场景**:
  ```go
  // 场景: 信号channel
  done := make(chan lang.PlaceholderType)
  go func() {
      doWork()
      done <- lang.Placeholder
  }()
  <-done
  ```

#### **Placeholder**
- **作用**: PlaceholderType的实例
- **应用场景**: 同上

---

## 20. limit - 限流器

### 包说明
提供多种限流算法实现。

### 核心类型

#### **PeriodLimit**
周期限流器，基于Redis实现。

#### **TokenLimitHandler**
令牌桶限流器。

### 主要函数

#### **NewPeriodLimit(period, quota int, limitStore *redis.Redis, keyPrefix string) *PeriodLimit**
- **作用**: 创建周期限流器
- **参数**:
    - `period`: 时间窗口（秒）
    - `quota`: 配额
    - `limitStore`: Redis客户端
    - `keyPrefix`: key前缀
- **应用场景**:
  ```go
  // 场景1: API限流（每分钟100次）
  limiter := limit.NewPeriodLimit(60, 100, rds, "api-limit")
  
  func handleRequest(w http.ResponseWriter, r *http.Request) {
      code, err := limiter.Take(getUserID(r))
      if code == limit.OverQuota {
          http.Error(w, "Too many requests", 429)
          return
      }
      // 处理请求
  }
  
  // 场景2: 短信发送限流（每天10条）
  limiter := limit.NewPeriodLimit(86400, 10, rds, "sms-limit")
  
  code, _ := limiter.Take(phoneNumber)
  if code == limit.Allowed {
      sendSMS(phoneNumber, message)
  }
  ```

#### **Take(key string) (int, error)**
- **作用**: 尝试获取令牌
- **返回值**:
    - `limit.Allowed`: 允许
    - `limit.HitQuota`: 达到配额
    - `limit.OverQuota`: 超过配额
- **应用场景**: 见上述示例

### 典型应用场景

1. **API限流**: 限制用户API调用频率
2. **短信限流**: 限制短信发送次数
3. **登录限流**: 防止暴力破解
4. **下载限流**: 限制下载次数

---

## 21. load - 负载统计

### 包说明
提供自适应负载统计和过载保护。

### 主要函数

#### **NewAdaptiveShedder(opts ...ShedderOption) Shedder**
- **作用**: 创建自适应过载保护器
- **应用场景**:
  ```go
  // 场景: 服务过载保护
  shedder := load.NewAdaptiveShedder()
  
  func handleRequest(w http.ResponseWriter, r *http.Request) {
      promise, err := shedder.Allow()
      if err != nil {
          http.Error(w, "Service overloaded", 503)
          return
      }
      defer promise.Pass() // 或 promise.Fail()
      
      // 处理请求
      processRequest(r)
  }
  ```

### 典型应用场景

1. **服务保护**: 防止服务过载崩溃
2. **降级处理**: 高负载时自动降级
3. **流量控制**: 动态调整处理能力

---

## 22. logc - 日志上下文

### 包说明
提供带上下文的日志功能。

### 主要函数

#### **Info(ctx context.Context, v ...any)**
- **作用**: 输出Info级别日志（带上下文）
- **应用场景**:
  ```go
  // 场景: 带trace ID的日志
  logc.Info(ctx, "User logged in")
  // 输出: [trace_id] User logged in
  ```

#### **Error(ctx context.Context, v ...any)**
- **作用**: 输出Error级别日志（带上下文）
- **应用场景**: 同上

---

## 23. logx - 日志系统

### 包说明
提供完整的日志系统，支持多种输出格式和级别。

### 主要函数

#### **Info(v ...any)**
- **作用**: 输出Info级别日志
- **应用场景**:
  ```go
  // 场景: 记录信息
  logx.Info("Server started on port 8080")
  ```

#### **Error(v ...any)**
- **作用**: 输出Error级别日志
- **应用场景**:
  ```go
  // 场景: 记录错误
  logx.Error("Failed to connect to database:", err)
  ```

#### **Infof(format string, v ...any)**
- **作用**: 格式化输出Info日志
- **应用场景**:
  ```go
  logx.Infof("User %s logged in from %s", username, ip)
  ```

#### **Errorf(format string, v ...any)**
- **作用**: 格式化输出Error日志
- **应用场景**: 同上

#### **Slow(v ...any)**
- **作用**: 输出慢日志
- **应用场景**:
  ```go
  // 场景: 记录慢查询
  if duration > time.Second {
      logx.Slow("Slow query:", sql, "duration:", duration)
  }
  ```

#### **Stat(v ...any)**
- **作用**: 输出统计日志
- **应用场景**:
  ```go
  // 场景: 记录统计信息
  logx.Stat("Request count:", count, "avg duration:", avgDuration)
  ```

#### **WithDuration(duration time.Duration) Logger**
- **作用**: 创建带持续时间的日志器
- **应用场景**:
  ```go
  // 场景: 记录请求耗时
  start := time.Now()
  processRequest()
  logx.WithDuration(time.Since(start)).Info("Request processed")
  ```

#### **MustSetup(c LogConf)**
- **作用**: 设置日志配置（失败则panic）
- **应用场景**:
  ```go
  // 场景: 应用启动时配置日志
  logx.MustSetup(logx.LogConf{
      ServiceName: "user-service",
      Mode:        "file",
      Path:        "/var/log/app",
      Level:       "info",
  })
  ```

### 典型应用场景

1. **应用日志**: 记录应用运行信息
2. **错误追踪**: 记录和追踪错误
3. **性能监控**: 记录慢查询、慢请求
4. **统计分析**: 记录统计数据

---

## 24. mapping - 映射工具

### 包说明
提供结构体映射和数据绑定功能。

### 主要函数

#### **UnmarshalKey(m map[string]any, v any) error**
- **作用**: 将map映射到结构体
- **应用场景**:
  ```go
  // 场景: 配置解析
  data := map[string]any{
      "host": "localhost",
      "port": 8080,
  }
  var config Config
  mapping.UnmarshalKey(data, &config)
  ```

---

## 25. mathx - 数学工具

### 包说明
提供数学计算相关的工具函数。

### 主要函数

#### **CalcPercent(val, total int64) float64**
- **作用**: 计算百分比
- **应用场景**:
  ```go
  // 场景: 计算成功率
  percent := mathx.CalcPercent(successCount, totalCount)
  fmt.Printf("Success rate: %.2f%%\n", percent)
  ```

#### **Max(a, b int) int**
- **作用**: 返回最大值
- **应用场景**:
  ```go
  // 场景: 取最大值
  maxValue := mathx.Max(value1, value2)
  ```

#### **Min(a, b int) int**
- **作用**: 返回最小值
- **应用场景**: 同上

---

## 26. metric - 指标监控

### 包说明
提供指标收集和监控功能。

### 主要函数

#### **NewHistogramVec(cfg *HistogramVecOpts) *HistogramVec**
- **作用**: 创建直方图指标
- **应用场景**:
  ```go
  // 场景: 监控请求耗时
  histogram := metric.NewHistogramVec(&metric.HistogramVecOpts{
      Namespace: "http",
      Subsystem: "requests",
      Name:      "duration_ms",
      Help:      "HTTP request duration in milliseconds",
      Labels:    []string{"method", "path"},
  })
  
  start := time.Now()
  processRequest()
  histogram.Observe(int64(time.Since(start)/time.Millisecond), method, path)
  ```

---

## 27. mr - MapReduce

### 包说明
提供进程内MapReduce并发处理框架。

### 主要函数

#### **MapReduce[T, U, V any](generate GenerateFunc[T], mapper MapperFunc[T, U], reducer ReducerFunc[U, V], opts ...Option) (V, error)**
- **作用**: 执行MapReduce操作
- **类型参数**:
    - `T`: 输入类型
    - `U`: 中间类型
    - `V`: 输出类型
- **应用场景**:
  ```go
  // 场景1: 并发查询商品详情
  type ProductID int
  type ProductDetail struct {
      ID    int
      Name  string
      Price float64
  }
  
  result, _ := mr.MapReduce(
      // Generate: 生成商品ID
      func(source chan<- ProductID) {
          for _, id := range productIDs {
              source <- id
          }
      },
      // Mapper: 并发查询商品详情
      func(id ProductID, writer mr.Writer[ProductDetail], cancel func(error)) {
          detail, err := queryProductDetail(id)
          if err != nil {
              cancel(err)
              return
          }
          writer.Write(detail)
      },
      // Reducer: 聚合结果
      func(pipe <-chan ProductDetail, writer mr.Writer[[]ProductDetail], cancel func(error)) {
          var products []ProductDetail
          for product := range pipe {
              products = append(products, product)
          }
          writer.Write(products)
      },
      mr.WithWorkers(10),
  )
  
  // 场景2: 并发计算
  sum, _ := mr.MapReduce(
      func(source chan<- int) {
          for i := 1; i <= 100; i++ {
              source <- i
          }
      },
      func(i int, writer mr.Writer[int], cancel func(error)) {
          writer.Write(i * i) // 计算平方
      },
      func(pipe <-chan int, writer mr.Writer[int], cancel func(error)) {
          var sum int
          for v := range pipe {
              sum += v
          }
          writer.Write(sum)
      },
  )
  ```

#### **MapReduceVoid[T, U any](generate GenerateFunc[T], mapper MapperFunc[T, U], reducer VoidReducerFunc[U], opts ...Option) error**
- **作用**: 执行MapReduce操作（无返回值）
- **应用场景**:
  ```go
  // 场景: 并发处理数据（无需返回值）
  mr.MapReduceVoid(
      func(source chan<- string) {
          for _, url := range urls {
              source <- url
          }
      },
      func(url string, writer mr.Writer[Response], cancel func(error)) {
          resp, err := http.Get(url)
          if err != nil {
              cancel(err)
              return
          }
          writer.Write(resp)
      },
      func(pipe <-chan Response, cancel func(error)) {
          for resp := range pipe {
              processResponse(resp)
          }
      },
  )
  ```

#### **ForEach[T any](generate GenerateFunc[T], mapper ForEachFunc[T], opts ...Option)**
- **作用**: 并发遍历处理（无输出）
- **应用场景**:
  ```go
  // 场景: 并发发送通知
  mr.ForEach(
      func(source chan<- User) {
          for _, user := range users {
              source <- user
          }
      },
      func(user User) {
          sendNotification(user)
      },
      mr.WithWorkers(20),
  )
  ```

#### **WithWorkers(workers int) Option**
- **作用**: 设置并发worker数量
- **应用场景**: 控制并发度

#### **WithContext(ctx context.Context) Option**
- **作用**: 设置上下文
- **应用场景**: 支持取消操作

### 典型应用场景

1. **并发RPC调用**: 并发查询多个服务组装数据
2. **批量数据处理**: 并发处理大量数据
3. **并发计算**: 并行计算任务
4. **数据聚合**: 并发查询后聚合结果

---

## 28. naming - 命名工具

### 包说明
提供服务命名相关的工具。

### 主要函数

#### **BuildTarget(endpoints []string) string**
- **作用**: 构建服务目标地址
- **应用场景**:
  ```go
  // 场景: 构建etcd服务地址
  target := naming.BuildTarget([]string{"etcd1:2379", "etcd2:2379"})
  ```

---

## 29. netx - 网络工具

### 包说明
提供网络相关的工具函数。

### 主要函数

#### **InternalIp() string**
- **作用**: 获取内网IP
- **应用场景**:
  ```go
  // 场景: 服务注册时获取本机IP
  ip := netx.InternalIp()
  registerService(ip, port)
  ```

---

## 30. proc - 进程管理

### 包说明
提供进程生命周期管理功能。

### 主要函数

#### **AddShutdownListener(fn func())**
- **作用**: 添加关闭监听器
- **应用场景**:
  ```go
  // 场景: 优雅关闭
  proc.AddShutdownListener(func() {
      log.Println("Shutting down...")
      db.Close()
      cache.Close()
  })
  ```

#### **AddWrapUpListener(fn func())**
- **作用**: 添加清理监听器
- **应用场景**: 同上

#### **Shutdown()**
- **作用**: 触发关闭流程
- **应用场景**:
  ```go
  // 场景: 手动触发关闭
  if criticalError {
      proc.Shutdown()
  }
  ```

---

## 31. prof - 性能分析

### 包说明
提供性能分析工具。

### 主要函数

#### **StartProfile() Stopper**
- **作用**: 开始性能分析
- **应用场景**:
  ```go
  // 场景: 性能分析
  stopper := prof.StartProfile()
  defer stopper.Stop()
  
  // 执行需要分析的代码
  performanceTest()
  ```

---

## 32. prometheus - Prometheus集成

### 包说明
提供Prometheus指标集成。

### 主要函数

#### **StartAgent(c Config)**
- **作用**: 启动Prometheus agent
- **应用场景**:
  ```go
  // 场景: 暴露metrics端点
  prometheus.StartAgent(prometheus.Config{
      Host: "0.0.0.0",
      Port: 9090,
      Path: "/metrics",
  })
  ```

---

## 33. queue - 队列

### 包说明
提供队列实现。

### 主要函数

#### **NewQueue(size int) *Queue**
- **作用**: 创建队列
- **应用场景**:
  ```go
  // 场景: 任务队列
  queue := queue.NewQueue(1000)
  
  // 生产者
  go func() {
      for task := range taskSource {
          queue.Put(task)
      }
  }()
  
  // 消费者
  go func() {
      for {
          task := queue.Take()
          processTask(task)
      }
  }()
  ```

---

## 34. rescue - 异常恢复

### 包说明
提供panic恢复功能。

### 主要函数

#### **Recover(cleanups ...func())**
- **作用**: 恢复panic
- **应用场景**:
  ```go
  // 场景: HTTP handler中恢复panic
  func handleRequest(w http.ResponseWriter, r *http.Request) {
      defer rescue.Recover(func() {
          log.Println("Recovered from panic")
      })
      
      // 可能panic的代码
      riskyOperation()
  }
  ```

---

## 35. search - 搜索工具

### 包说明
提供搜索相关的工具。

---

## 36. service - 服务框架

### 包说明
提供服务框架基础功能。

### 主要函数

#### **NewServiceGroup() *ServiceGroup**
- **作用**: 创建服务组
- **应用场景**:
  ```go
  // 场景: 管理多个服务
  group := service.NewServiceGroup()
  group.Add(httpServer)
  group.Add(grpcServer)
  group.Start()
  ```

---

## 37. stat - 统计工具

### 包说明
提供统计功能。

### 主要函数

#### **NewMetrics(name string) *Metrics**
- **作用**: 创建指标统计
- **应用场景**:
  ```go
  // 场景: 统计请求
  metrics := stat.NewMetrics("http_requests")
  metrics.Add(stat.Task{
      Duration: duration,
  })
  ```

---

## 38. stores - 存储

### 包说明
提供统一的存储接口，包括Redis、SQL、MongoDB等。

### 核心功能

1. **Redis**: Redis客户端封装
2. **SQL**: 数据库操作封装
3. **Cache**: 缓存封装
4. **MongoDB**: MongoDB客户端封装

---

## 39. stringx - 字符串工具

### 包说明
提供字符串处理工具函数。

### 主要函数

#### **Contains(list []string, str string) bool**
- **作用**: 检查字符串是否在列表中
- **应用场景**:
  ```go
  // 场景: 权限检查
  if stringx.Contains(allowedRoles, userRole) {
      // 允许访问
  }
  ```

#### **Filter(s string, filter func(r rune) bool) string**
- **作用**: 过滤字符串中的字符
- **应用场景**:
  ```go
  // 场景: 移除特殊字符
  cleaned := stringx.Filter(input, func(r rune) bool {
      return unicode.IsLetter(r) || unicode.IsDigit(r)
  })
  ```

#### **FirstN(s string, n int, ellipsis ...string) string**
- **作用**: 获取前N个字符
- **应用场景**:
  ```go
  // 场景: 文本截断
  preview := stringx.FirstN(content, 100, "...")
  ```

#### **HasEmpty(args ...string) bool**
- **作用**: 检查是否有空字符串
- **应用场景**:
  ```go
  // 场景: 参数验证
  if stringx.HasEmpty(username, password, email) {
      return errors.New("missing required fields")
  }
  ```

#### **NotEmpty(args ...string) bool**
- **作用**: 检查所有字符串都不为空
- **应用场景**: 同上

#### **Remove(strings []string, strs ...string) []string**
- **作用**: 从列表中移除指定字符串
- **应用场景**:
  ```go
  // 场景: 移除黑名单
  cleaned := stringx.Remove(allUsers, bannedUsers...)
  ```

#### **Reverse(s string) string**
- **作用**: 反转字符串
- **应用场景**:
  ```go
  // 场景: 字符串反转
  reversed := stringx.Reverse("hello") // "olleh"
  ```

#### **Substr(str string, start, stop int) (string, error)**
- **作用**: 获取子字符串
- **应用场景**:
  ```go
  // 场景: 字符串切片
  sub, _ := stringx.Substr("hello world", 0, 5) // "hello"
  ```

#### **TakeOne(valid, or string) string**
- **作用**: 返回第一个非空字符串
- **应用场景**:
  ```go
  // 场景: 默认值
  value := stringx.TakeOne(userInput, defaultValue)
  ```

#### **ToCamelCase(s string) string**
- **作用**: 转换为驼峰命名
- **应用场景**:
  ```go
  // 场景: 命名转换
  camel := stringx.ToCamelCase("HelloWorld") // "helloWorld"
  ```

#### **Union(first, second []string) []string**
- **作用**: 合并字符串列表（去重）
- **应用场景**:
  ```go
  // 场景: 合并标签
  allTags := stringx.Union(tags1, tags2)
  ```

### 典型应用场景

1. **参数验证**: 检查必填字段
2. **文本处理**: 截断、过滤、转换
3. **列表操作**: 合并、去重、移除
4. **字符串工具**: 反转、切片、命名转换

---

## 40. syncx - 同步工具

### 包说明
提供同步原语和并发控制工具。

### 核心类型

#### **Barrier**
屏障，用于保护资源访问。

#### **SpinLock**
自旋锁，用于快速执行的锁。

#### **SingleFlight**
单飞模式，合并并发相同请求。

#### **LockedCalls**
锁定调用，保证相同key的调用顺序执行。

#### **OnceGuard**
一次性守卫，保证资源只被获取一次。

#### **Cond**
条件变量。

#### **DoneChan**
完成channel，可多次关闭。

### 主要函数

#### **NewBarrier() *Barrier**
- **作用**: 创建屏障
- **应用场景**:
  ```go
  // 场景: 保护共享资源
  var barrier syncx.Barrier
  barrier.Guard(func() {
      // 临界区代码
      sharedResource.Update()
  })
  ```

#### **NewSpinLock() *SpinLock**
- **作用**: 创建自旋锁
- **应用场景**:
  ```go
  // 场景: 快速锁定
  var lock syncx.SpinLock
  lock.Lock()
  defer lock.Unlock()
  // 快速操作
  ```

#### **NewSingleFlight() SingleFlight**
- **作用**: 创建单飞实例
- **应用场景**:
  ```go
  // 场景: 缓存击穿防护
  sf := syncx.NewSingleFlight()
  
  func getUser(id string) (*User, error) {
      v, err := sf.Do(id, func() (any, error) {
          // 只有第一个请求会执行
          return db.QueryUser(id)
      })
      return v.(*User), err
  }
  
  // 场景2: 防止缓存雪崩
  v, shared, err := sf.DoEx("key", func() (any, error) {
      return expensiveOperation()
  })
  if shared {
      log.Println("Result was shared from another call")
  }
  ```

#### **NewLockedCalls() LockedCalls**
- **作用**: 创建锁定调用实例
- **应用场景**:
  ```go
  // 场景: 保证相同key的调用顺序执行
  lc := syncx.NewLockedCalls()
  
  func processUser(userID string) error {
      _, err := lc.Do(userID, func() (any, error) {
          // 相同userID的调用会排队执行
          return updateUser(userID)
      })
      return err
  }
  ```

#### **NewOnceGuard() *OnceGuard**
- **作用**: 创建一次性守卫
- **应用场景**:
  ```go
  // 场景: 保证资源只被获取一次
  var guard syncx.OnceGuard
  
  if guard.Take() {
      // 只有第一个调用者会执行
      initializeResource()
  }
  
  if guard.Taken() {
      // 检查资源是否已被获取
  }
  ```

#### **NewCond() *Cond**
- **作用**: 创建条件变量
- **应用场景**:
  ```go
  // 场景: 等待条件满足
  cond := syncx.NewCond()
  
  go func() {
      time.Sleep(time.Second)
      cond.Signal() // 发送信号
  }()
  
  cond.Wait() // 等待信号
  
  // 场景2: 超时等待
  remain, ok := cond.WaitWithTimeout(5*time.Second)
  if !ok {
      log.Println("Timeout")
  }
  ```

#### **NewDoneChan() *DoneChan**
- **作用**: 创建完成channel
- **应用场景**:
  ```go
  // 场景: 可多次关闭的done channel
  done := syncx.NewDoneChan()
  
  go func() {
      <-done.Done()
      cleanup()
  }()
  
  // 可以安全地多次调用
  done.Close()
  done.Close() // 不会panic
  ```

### 典型应用场景

1. **SingleFlight**: 缓存击穿防护、防止重复请求
2. **LockedCalls**: 顺序执行相同key的操作
3. **OnceGuard**: 单例初始化、资源获取
4. **Barrier**: 保护共享资源
5. **SpinLock**: 快速锁定场景
6. **Cond**: 条件等待、信号通知
7. **DoneChan**: 优雅关闭、多次关闭安全

---

## 41. sysx - 系统工具

### 包说明
提供系统相关的工具函数。

### 主要函数

#### **Hostname() string**
- **作用**: 获取主机名
- **应用场景**:
  ```go
  // 场景: 服务标识
  hostname := sysx.Hostname()
  log.Printf("Service running on %s", hostname)
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

## 44. trace - 链路追踪

### 包说明
提供分布式链路追踪功能，集成OpenTelemetry。

### 主要函数

#### **StartServerSpan(ctx context.Context, carrier propagation.TextMapCarrier, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)**
- **作用**: 启动服务端span
- **应用场景**:
  ```go
  // 场景: HTTP服务端追踪
  ctx, span := trace.StartServerSpan(r.Context(), propagation.HeaderCarrier(r.Header), "HandleRequest")
  defer span.End()
  
  // 处理请求
  processRequest(ctx)
  ```

#### **StartClientSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)**
- **作用**: 启动客户端span
- **应用场景**:
  ```go
  // 场景: HTTP客户端追踪
  ctx, span := trace.StartClientSpan(ctx, "CallAPI")
  defer span.End()
  
  resp, err := http.Get(url)
  if err != nil {
      span.RecordError(err)
  }
  ```

### 典型应用场景

1. **分布式追踪**: 追踪请求在微服务间的调用链路
2. **性能分析**: 分析各个环节的耗时
3. **错误追踪**: 追踪错误发生的位置和传播路径

---

## 45. utils - 通用工具

### 包说明
提供通用工具函数。

---

## 46. validation - 数据验证

### 包说明
提供数据验证功能。

### 主要函数

#### **Validate(v any) error**
- **作用**: 验证结构体
- **应用场景**:
  ```go
  // 场景: 请求参数验证
  type CreateUserRequest struct {
      Username string `validate:"required,min=3,max=20"`
      Email    string `validate:"required,email"`
      Age      int    `validate:"required,min=18,max=120"`
  }
  
  req := CreateUserRequest{
      Username: "john",
      Email:    "john@example.com",
      Age:      25,
  }
  
  if err := validation.Validate(req); err != nil {
      return fmt.Errorf("validation failed: %w", err)
  }
  ```

### 典型应用场景

1. **API参数验证**: 验证HTTP请求参数
2. **配置验证**: 验证配置文件的有效性
3. **数据完整性**: 验证数据模型的完整性

---

## 总结

### 核心包分类

#### **基础工具类**
- `lang`: 语言基础工具（PlaceholderType、Placeholder）
- `stringx`: 字符串处理（Filter、FirstN、Reverse、ToCamelCase等）
- `mathx`: 数学计算（CalcPercent、Max、Min）
- `timex`: 时间处理（Now、Since、Time）
- `hash`: 哈希算法（一致性哈希）

#### **并发编程类**
- `threading`: 并发工具（RoutineGroup、TaskRunner、StableRunner、WorkerGroup）
- `syncx`: 同步工具（SingleFlight、LockedCalls、Barrier、SpinLock、OnceGuard、Cond、DoneChan）
- `executors`: 执行器（BulkExecutor、ChunkExecutor、DelayExecutor、LessExecutor、PeriodicalExecutor）
- `mr`: MapReduce（并发数据处理框架）
- `fx`: 函数式编程（Stream流式处理）

#### **IO处理类**
- `iox`: IO扩展（BufferPool、DupReadCloser、TextLineScanner、ReadTextLines等）
- `fs`: 文件系统（TempFileWithText）
- `filex`: 文件扩展（RangeReader）

#### **网络通信类**
- `netx`: 网络工具（InternalIp）
- `discov`: 服务发现（基于etcd的服务注册与发现）
- `naming`: 命名工具（BuildTarget）

#### **数据存储类**
- `stores`: 存储（Redis、SQL、MongoDB、Cache统一接口）
- `collection`: 集合数据结构（Cache、Ring、Set、TimingWheel）
- `bloom`: 布隆过滤器（防缓存穿透、URL去重）

#### **可靠性保障类**
- `breaker`: 熔断器（服务保护、自动降级）
- `limit`: 限流器（PeriodLimit周期限流、TokenLimit令牌桶）
- `load`: 负载统计（自适应过载保护）
- `rescue`: 异常恢复（Recover panic恢复）
- `errorx`: 错误处理（Wrap、Wrapf错误包装）

#### **监控观测类**
- `logx`: 日志系统（Info、Error、Slow、Stat等多级别日志）
- `logc`: 日志上下文（带Context的日志）
- `metric`: 指标监控（HistogramVec直方图指标）
- `stat`: 统计工具（Metrics统计）
- `trace`: 链路追踪（OpenTelemetry集成）
- `prof`: 性能分析（StartProfile）
- `prometheus`: Prometheus集成（StartAgent）

#### **配置管理类**
- `conf`: 配置加载（支持JSON、YAML、TOML多格式）
- `configcenter`: 配置中心（动态配置更新）

#### **编解码类**
- `codec`: 编解码（RSA、AES、HMAC、MD5等加密算法）
- `jsonx`: JSON扩展（Marshal、Unmarshal）
- `mapping`: 映射工具（UnmarshalKey结构体映射）

#### **服务框架类**
- `service`: 服务框架（ServiceGroup服务组管理）
- `proc`: 进程管理（AddShutdownListener优雅关闭）
- `queue`: 队列（Queue任务队列）

#### **其他工具类**
- `cmdline`: 命令行工具（EnterToContinue交互式确认）
- `color`: 终端颜色（WithColor彩色输出）
- `contextx`: Context扩展（ValueOnlyFrom）
- `sysx`: 系统工具（Hostname）
- `search`: 搜索工具
- `utils`: 通用工具
- `validation`: 数据验证（Validate结构体验证）

---

### 使用建议

#### 1. 并发处理场景

**选择指南**：
- **简单并发**：使用 `threading.RoutineGroup`
- **限制并发数**：使用 `threading.TaskRunner`
- **保持顺序**：使用 `threading.StableRunner`
- **固定Worker**：使用 `threading.WorkerGroup`
- **批量处理**：使用 `executors.BulkExecutor`
- **MapReduce**：使用 `mr.MapReduce`
- **流式处理**：使用 `fx.Stream`

#### 2. 限流熔断场景

**选择指南**：
- **API限流**：使用 `limit.PeriodLimit`
- **服务熔断**：使用 `breaker.Breaker`
- **过载保护**：使用 `load.AdaptiveShedder`
- **频率限制**：使用 `executors.LessExecutor`

#### 3. 缓存场景

**选择指南**：
- **LRU缓存**：使用 `collection.Cache`
- **防击穿**：使用 `syncx.SingleFlight`
- **防穿透**：使用 `bloom.Filter`
- **分布式缓存**：使用 `stores.Cache`

#### 4. 日志监控场景

**选择指南**：
- **普通日志**：使用 `logx`
- **带Context**：使用 `logc`
- **链路追踪**：使用 `trace`
- **指标监控**：使用 `metric` + `prometheus`
- **统计分析**：使用 `stat`

#### 5. 数据处理场景

**选择指南**：
- **字符串处理**：使用 `stringx`
- **JSON处理**：使用 `jsonx`
- **IO处理**：使用 `iox`
- **数据验证**：使用 `validation`
- **数据映射**：使用 `mapping`

#### 6. 存储场景

**选择指南**：
- **Redis**：使用 `stores/redis`
- **MySQL**：使用 `stores/sqlx`
- **MongoDB**：使用 `stores/mongo`
- **缓存**：使用 `stores/cache`

---

### 最佳实践

#### 1. 并发控制

```go
// ✅ 推荐：使用TaskRunner限制并发
runner := threading.NewTaskRunner(10)
for _, task := range tasks {
    runner.Schedule(func() {
        processTask(task)
    })
}
runner.Wait()

// ❌ 不推荐：无限制并发
for _, task := range tasks {
    go processTask(task)
}
```

#### 2. 错误处理

```go
// ✅ 推荐：使用errorx包装错误
if err := db.Query(); err != nil {
    return errorx.Wrap(err, "failed to query database")
}

// ✅ 推荐：使用rescue恢复panic
defer rescue.Recover(func() {
    log.Println("Recovered from panic")
})
```

#### 3. 资源管理

```go
// ✅ 推荐：使用proc管理生命周期
proc.AddShutdownListener(func() {
    db.Close()
    cache.Close()
})

// ✅ 推荐：使用iox.BufferPool复用Buffer
var bufPool = iox.NewBufferPool(4096)
buf := bufPool.Get()
defer bufPool.Put(buf)
```

#### 4. 性能优化

```go
// ✅ 推荐：使用SingleFlight防止缓存击穿
sf := syncx.NewSingleFlight()
v, err := sf.Do(key, func() (any, error) {
    return db.Query(key)
})

// ✅ 推荐：使用一致性哈希分布式缓存
hash := hash.NewConsistentHash()
server, _ := hash.Get(key)
```

---

### 性能对比

| 场景 | 传统方式 | go-zero方式 | 性能提升 |
|------|---------|------------|---------|
| 并发处理 | 无限制goroutine | TaskRunner | 资源可控 |
| 缓存击穿 | 加锁 | SingleFlight | 减少90%+请求 |
| 批量操作 | 逐个处理 | BulkExecutor | 提升10倍+ |
| Buffer分配 | 每次new | BufferPool | 减少90%+GC |
| 日志输出 | fmt.Println | logx | 结构化+高性能 |

---

### 常见问题

#### Q1: 什么时候使用StableRunner？
**A**: 当需要并发处理但必须保持输出顺序时，如Kafka消息处理、顺序写入数据库。

#### Q2: SingleFlight和LockedCalls的区别？
**A**:
- `SingleFlight`: 合并并发请求，共享结果
- `LockedCalls`: 串行执行，不共享结果

#### Q3: 如何选择限流器？
**A**:
- 固定窗口限流：`limit.PeriodLimit`
- 令牌桶限流：`limit.TokenLimitHandler`
- 自适应限流：`load.AdaptiveShedder`

#### Q4: 如何实现优雅关闭？
**A**: 使用 `proc.AddShutdownListener` 注册清理函数。

#### Q5: 如何防止缓存穿透？
**A**: 使用 `bloom.Filter` 布隆过滤器。

---

### 学习路径

#### 初级（必学）
1. `threading`: 并发基础
2. `logx`: 日志系统
3. `conf`: 配置加载
4. `errorx`: 错误处理
5. `stringx`: 字符串工具

#### 中级（推荐）
1. `breaker`: 熔断器
2. `limit`: 限流器
3. `syncx`: 同步工具
4. `executors`: 执行器
5. `mr`: MapReduce

#### 高级（进阶）
1. `fx`: 函数式编程
2. `load`: 自适应限流
3. `trace`: 链路追踪
4. `metric`: 指标监控
5. `stores`: 存储抽象

---

### 参考资源

- **官方文档**: https://go-zero.dev/
- **GitHub**: https://github.com/zeromicro/go-zero
- **示例代码**: https://github.com/zeromicro/zero-examples
- **社区讨论**: https://github.com/zeromicro/go-zero/discussions

---

**文档结束**

> **版本**: v1.0  
> **最后更新**: 2025-12-30  
> **维护者**: go-zero 社区  
> **许可**: MIT License
>
> 本文档持续更新中，如有疑问或建议，欢迎通过GitHub Issues反馈。
