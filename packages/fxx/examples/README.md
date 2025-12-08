# fx 框架示例集合

每个示例都是独立的，可以单独运行。

## 📂 目录结构

```
examples/
├── 01-basic/          # 基础依赖注入
├── 02-http-server/    # HTTP 服务器 + 生命周期
├── 03-lifecycle/      # 生命周期管理详解
├── 04-modules/        # 模块化架构
└── 05-advanced/       # 高级特性（可选依赖、值组）
```

## 🚀 运行方式

### 1. 基础示例
```bash
cd 01-basic
go mod tidy
go run main.go
```

**学习内容：**
- fx.Provide() 注册依赖
- fx.Invoke() 执行逻辑
- 自动依赖注入
- shutdowner.Shutdown() 优雅退出

### 2. HTTP 服务器示例
```bash
cd 02-http-server
go mod tidy
go run main.go

# 在另一个终端测试：
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/time

# 按 Ctrl+C 停止服务器
```

**学习内容：**
- 生命周期钩子（OnStart/OnStop）
- HTTP 服务器的启动和优雅关闭
- 长时间运行的服务

### 3. 生命周期管理示例
```bash
cd 03-lifecycle
go mod tidy
go run main.go
```

**学习内容：**
- 数据库/缓存的启动和关闭
- 组件启动顺序
- 资源清理

### 4. 模块化示例
```bash
cd 04-modules
go mod tidy
go run main.go
```

**学习内容：**
- fx.Module() 组织代码
- 模块化架构设计
- 多个模块协作

### 5. 高级特性示例
```bash
cd 05-advanced
go mod tidy
go run main.go
```

**学习内容：**
- 可选依赖 `optional:"true"`
- 值组 `group:"name"`
- fx.In/fx.Out 结构体
- 复杂依赖管理

## 📚 学习路径

建议按以下顺序学习：

1. **01-basic** → 理解基本概念
2. **03-lifecycle** → 掌握生命周期
3. **02-http-server** → 实践真实场景
4. **04-modules** → 学习架构设计
5. **05-advanced** → 探索高级特性

## 💡 关键概念

### fx.Provide()
注册构造函数，告诉 fx 如何创建对象。

```go
fx.Provide(
    NewLogger,
    NewDatabase,
    NewService,
)
```

### fx.Invoke()
请求依赖并执行代码。

```go
fx.Invoke(func(service *Service) {
    service.Run()
})
```

### fx.Lifecycle
管理组件的启动和关闭。

```go
lc.Append(fx.Hook{
    OnStart: func(ctx context.Context) error {
        // 启动逻辑
        return nil
    },
    OnStop: func(ctx context.Context) error {
        // 关闭逻辑
        return nil
    },
})
```

### fx.Shutdowner
程序化地关闭应用。

```go
shutdowner.Shutdown()
```

## 🎯 常见问题

### Q: 为什么 app.Run() 不退出？
A: 这是正常的。fx 应用会一直运行，等待系统信号（Ctrl+C）。如果要自动退出，使用 `shutdowner.Shutdown()`。

### Q: 如何查看详细的依赖注入过程？
A: 运行时会看到 `[Fx]` 前缀的日志，显示依赖的创建过程。

### Q: 如何调试依赖注入问题？
A: 检查构造函数的参数和返回值类型，确保它们匹配。

## 🔗 更多资源

- [fx 官方文档](https://uber-go.github.io/fx/)
- [fx GitHub](https://github.com/uber-go/fx)
- 上级目录的 README.md

祝学习愉快！🚀
