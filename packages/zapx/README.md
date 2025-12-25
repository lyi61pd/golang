# Zap 日志库学习指南

Uber 开源的高性能 Go 日志库学习示例

## 📚 目录结构

- `main.go` - 基础用法教程（从入门到进阶）
- `advanced.go` - 实战技巧和最佳实践

## 🚀 快速开始

```bash
# 运行基础教程
go run main.go

# 运行实战示例
go run advanced.go
```

## 📖 学习路径

### 1. 基础篇 (main.go)

1. **预设 Logger**
   - Development Logger - 开发环境用
   - Production Logger - 生产环境用
   - Example Logger - 演示用

2. **自定义配置**
   - 使用 Config 结构体
   - 自定义编码器
   - 配置时间格式

3. **结构化日志**
   - 基本类型字段：String, Int, Bool, Float64
   - 时间和持续时间：Time, Duration
   - 复杂类型：Strings, Ints, Any
   - 错误字段：Error

4. **日志级别**
   - Debug - 调试信息
   - Info - 重要信息
   - Warn - 警告
   - Error - 错误
   - Fatal/Panic - 致命错误

5. **SugaredLogger**
   - Printf 风格：`Infof()`
   - 键值对风格：`Infow()`
   - 性能对比

6. **实际应用**
   - HTTP 请求日志
   - 数据库操作日志
   - 使用 With() 添加上下文

7. **输出配置**
   - 控制台输出
   - 文件输出
   - 自定义编码器

### 2. 进阶篇 (advanced.go)

1. **全局 Logger 模式**
   - 单例模式
   - 子 Logger 创建

2. **日志分级输出**
   - 不同级别输出到不同目标
   - Info → 控制台
   - Error → 文件

3. **动态调整级别**
   - AtomicLevel 使用
   - 运行时调整

4. **上下文字段**
   - Request ID
   - User ID
   - 链式添加

5. **性能优化**
   - Logger vs SugaredLogger
   - 性能测试对比

6. **日志采样**
   - 高频场景优化
   - 避免日志爆炸

7. **错误处理**
   - 业务错误
   - 系统错误
   - 堆栈跟踪
   - 错误恢复

## 💡 核心概念

### Logger vs SugaredLogger

```go
// Logger - 类型安全，零内存分配，最快
logger.Info("user login",
    zap.String("username", "zhangsan"),
    zap.Int("user_id", 123))

// SugaredLogger - 灵活方便，有少量分配，稍慢
sugar.Infow("user login",
    "username", "zhangsan",
    "user_id", 123)
```

### 结构化 vs 非结构化

```go
// ❌ 不推荐：非结构化
logger.Info("User zhangsan login with ID 123")

// ✅ 推荐：结构化
logger.Info("user login",
    zap.String("username", "zhangsan"),
    zap.Int("user_id", 123))
```

### 日志级别选择

| 级别 | 使用场景 | 示例 |
|------|----------|------|
| Debug | 开发调试 | 变量值、方法调用 |
| Info | 重要流程 | 用户登录、订单创建 |
| Warn | 可恢复异常 | 重试成功、降级处理 |
| Error | 需要关注的错误 | 数据库错误、API 失败 |
| Fatal | 程序无法运行 | 配置错误、启动失败 |

## 🎯 最佳实践

### 1. 初始化

```go
var Logger *zap.Logger

func init() {
    var err error
    if os.Getenv("ENV") == "production" {
        Logger, err = zap.NewProduction()
    } else {
        Logger, err = zap.NewDevelopment()
    }
    if err != nil {
        panic(err)
    }
}
```

### 2. 添加上下文

```go
// 为每个请求创建带上下文的 Logger
func HandleRequest(ctx context.Context) {
    requestLogger := Logger.With(
        zap.String("request_id", GetRequestID(ctx)),
        zap.String("user_id", GetUserID(ctx)),
    )
    
    requestLogger.Info("processing request")
    // 后续所有日志自动带上 request_id 和 user_id
}
```

### 3. 错误日志

```go
// 记录错误时提供足够的上下文
if err != nil {
    logger.Error("failed to query database",
        zap.Error(err),
        zap.String("query", sql),
        zap.String("table", "users"),
        zap.Duration("elapsed", time.Since(start)))
    return err
}
```

### 4. 性能敏感场景

```go
// 使用 Check 避免不必要的字段构造
if ce := logger.Check(zap.DebugLevel, "debug info"); ce != nil {
    ce.Write(
        zap.String("key", expensiveOperation()),
    )
}
```

## 📦 依赖

```bash
go get -u go.uber.org/zap
```

## 🔗 参考资料

- [Zap GitHub](https://github.com/uber-go/zap)
- [官方文档](https://pkg.go.dev/go.uber.org/zap)
- [性能对比](https://github.com/uber-go/zap#performance)

## 📝 常见问题

### Q: Logger 和 SugaredLogger 该选哪个？

**A**: 
- 性能敏感场景（热路径）→ Logger
- 一般业务代码 → SugaredLogger
- 推荐：90% 场景用 Logger，习惯后非常自然

### Q: 如何输出到文件？

**A**: 
```go
config := zap.NewProductionConfig()
config.OutputPaths = []string{"stdout", "./logs/app.log"}
logger, _ := config.Build()
```

推荐使用 [lumberjack](https://github.com/natefinch/lumberjack) 实现日志轮转。

### Q: 如何彩色输出？

**A**:
```go
config := zap.NewDevelopmentEncoderConfig()
config.EncodeLevel = zapcore.CapitalColorLevelEncoder
encoder := zapcore.NewConsoleEncoder(config)
```

### Q: 生产环境推荐配置？

**A**:
```go
logger, _ := zap.NewProduction() // JSON 格式，性能优化
defer logger.Sync()
```

## 🎓 学习建议

1. 先运行 `main.go`，理解基本概念
2. 阅读代码注释，动手修改参数
3. 运行 `advanced.go`，学习实战技巧
4. 在自己的项目中实践
5. 关注性能和结构化日志的平衡

---

Happy Logging! 🚀
