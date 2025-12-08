# Go fx 框架学习指南

## 📚 简介

fx 是 Uber 开发的一个 Go 依赖注入框架，它可以帮助你：
- 自动管理依赖关系
- 控制组件的生命周期
- 构建模块化、可测试的应用程序
- 避免全局变量和手动初始化代码

## 🎯 核心概念

### 1. 依赖注入（Dependency Injection）
fx 会自动分析构造函数的参数，并注入所需的依赖。

```go
// fx 会自动注入 Logger
func NewGreeter(logger Logger) *Greeter {
    return &Greeter{logger: logger}
}
```

### 2. 提供者（Providers）
使用 `fx.Provide()` 注册构造函数，告诉 fx 如何创建对象。

```go
fx.Provide(
    NewLogger,    // 提供 Logger
    NewGreeter,   // 提供 Greeter
)
```

### 3. 调用者（Invokers）
使用 `fx.Invoke()` 请求依赖并执行代码。

```go
fx.Invoke(func(greeter *Greeter) {
    greeter.Greet("World")
})
```

### 4. 生命周期（Lifecycle）
fx.Lifecycle 允许你注册启动和关闭钩子。

```go
func NewServer(lc fx.Lifecycle) *Server {
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            // 启动服务器
            return nil
        },
        OnStop: func(ctx context.Context) error {
            // 关闭服务器
            return nil
        },
    })
    return srv
}
```

### 5. 模块（Modules）
使用 `fx.Module()` 将相关的提供者组织在一起。

```go
var DatabaseModule = fx.Module("database",
    fx.Provide(
        NewDatabase,
        NewCache,
    ),
)
```

## 📖 示例说明

### `main.go` - 基础示例
最简单的 fx 应用，演示：
- 如何定义构造函数
- 如何使用 fx.Provide 注册依赖
- 如何使用 fx.Invoke 执行代码
- fx 如何自动解决依赖关系

**运行：**
```bash
go run main.go
```

### `02_http_server.go` - HTTP 服务器示例
演示如何使用 fx 构建 HTTP 服务器：
- 生命周期管理（启动/关闭）
- 优雅关闭
- 配置管理

**运行：**
```bash
# 1. 取消注释文件末尾的 main 函数
# 2. 注释掉 main.go 中的 main 函数
go run 02_http_server.go main.go

# 在另一个终端测试：
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/time
```

### `03_lifecycle.go` - 生命周期管理
演示：
- 如何注册 OnStart 和 OnStop 钩子
- 组件启动和关闭的顺序
- 资源清理

**运行：**
```bash
# 取消注释该文件的 main 函数
go run 03_lifecycle.go main.go
```

### `04_modules.go` - 模块化
演示：
- 如何使用 fx.Module 组织代码
- 模块化架构
- 依赖多个实现

**运行：**
```bash
# 取消注释该文件的 main 函数
go run 04_modules.go main.go 03_lifecycle.go
```

### `05_advanced.go` - 高级特性
演示高级用法：
- **结构体参数** (`fx.In`)：避免参数过多
- **结构体结果** (`fx.Out`)：提供多个返回值
- **可选依赖** (`optional:"true"`)：依赖可能不存在
- **值组** (`group:"name"`)：收集同类型的多个实现
- **命名依赖** (`name:"xyz"`)：区分同类型的不同实例

**运行：**
```bash
# 取消注释该文件的 main 函数
go run 05_advanced.go main.go 02_http_server.go 03_lifecycle.go
```

## 🚀 快速开始

### 1. 安装依赖
```bash
go get go.uber.org/fx
```

### 2. 运行基础示例
```bash
cd /Users/ybbj100324/code/self/golang_demo/packages/fxx
go run main.go
```

### 3. 查看输出
你会看到：
- 依赖的创建顺序
- 自动注入过程
- 应用程序运行过程

## 📝 学习路径

1. **第一步**：运行 `main.go`，理解基础的依赖注入
2. **第二步**：学习 `03_lifecycle.go`，掌握生命周期管理
3. **第三步**：实践 `02_http_server.go`，构建真实的服务
4. **第四步**：探索 `04_modules.go`，学习模块化设计
5. **第五步**：研究 `05_advanced.go`，掌握高级特性

## 🎓 最佳实践

### ✅ 推荐做法

```go
// 1. 构造函数返回接口而不是具体类型
func NewLogger() Logger {
    return &SimpleLogger{}
}

// 2. 使用 fx.In 避免参数过多
type Params struct {
    fx.In
    Logger Logger
    Config *Config
    DB     *Database
}

func NewService(p Params) *Service {
    // ...
}

// 3. 使用生命周期管理资源
func NewServer(lc fx.Lifecycle) *Server {
    srv := &Server{}
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            return srv.Start()
        },
        OnStop: func(ctx context.Context) error {
            return srv.Stop()
        },
    })
    return srv
}

// 4. 使用模块组织代码
var MyModule = fx.Module("mymodule",
    fx.Provide(NewA, NewB, NewC),
)
```

### ❌ 避免做法

```go
// ❌ 不要在构造函数中做耗时操作
func NewService() *Service {
    // 错误：应该在 OnStart 钩子中执行
    time.Sleep(10 * time.Second)
    return &Service{}
}

// ❌ 不要使用全局变量
var globalLogger Logger // 应该通过依赖注入

// ❌ 不要手动管理依赖
func main() {
    logger := NewLogger()
    db := NewDatabase(logger)
    service := NewService(logger, db)
    // 应该让 fx 自动处理
}
```

## 🔍 常见问题

### Q: fx 如何知道依赖的创建顺序？
A: fx 会分析构造函数的参数，自动计算依赖图，并按拓扑排序创建对象。

### Q: 如果有循环依赖怎么办？
A: fx 会在启动时检测并报错，你需要重新设计避免循环依赖。

### Q: 如何测试使用 fx 的代码？
A: 可以创建测试专用的 fx.App，或者直接测试构造函数（它们是普通函数）。

### Q: fx.Invoke 和直接调用有什么区别？
A: fx.Invoke 会在所有依赖创建完成后调用，并且能自动注入参数。

## 📚 更多资源

- [fx 官方文档](https://uber-go.github.io/fx/)
- [fx GitHub 仓库](https://github.com/uber-go/fx)
- [Uber Go 风格指南](https://github.com/uber-go/guide)

## 💡 实战建议

1. **从小开始**：先在小项目中使用 fx，熟悉后再应用到大项目
2. **合理拆分**：将应用拆分为多个模块，每个模块负责一个领域
3. **测试优先**：fx 让测试更容易，充分利用这个优势
4. **文档清晰**：为每个提供者添加注释，说明它的作用和依赖
5. **生命周期管理**：所有需要清理的资源都应该使用生命周期钩子

## 🎯 下一步

尝试：
1. 修改示例代码，添加自己的组件
2. 构建一个完整的 REST API 服务
3. 添加数据库、缓存、消息队列等真实依赖
4. 实现配置热加载
5. 添加监控和指标收集

祝学习愉快！🚀
