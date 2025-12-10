`fx.Module` 是 fx 的**模块化组织工具**，用来把相关的 Provider 打包在一起，让代码结构更清晰。

## 🎯 核心概念

```go
// 没有 Module：所有东西堆在一起
fx.New(
    fx.Provide(
        NewLogger,
        NewDatabase,
        NewCache,
        NewEmailSender,
        NewSMSSender,
        NewUserService,
        NewNotificationService,
        // ... 100 个 Provider
    ),
)
// ❌ 混乱，不知道哪些是一组的

// 使用 Module：按功能分组
var DatabaseModule = fx.Module("database",
    fx.Provide(NewDatabase, NewCache),
)

var NotificationModule = fx.Module("notification",
    fx.Provide(NewEmailSender, NewSMSSender),
)

fx.New(
    DatabaseModule,
    NotificationModule,
    // ...
)
// ✅ 清晰，一眼看出模块划分
```

## 📊 你的代码解析

```go
// 1. 定义模块（只是声明，不执行）
var LoggingModule = fx.Module("logging",  // 模块名（用于调试）
    fx.Provide(NewSimpleLogger),          // 这个模块提供什么
)

var DataModule = fx.Module("data",
    fx.Provide(
        NewDatabase,   // 数据相关的东西放一起
        NewCache,
    ),
)

// 2. 使用模块（组装应用）
fx.New(
    LoggingModule,        // 引入日志模块
    DataModule,           // 引入数据模块
    ServiceModule,        // 引入服务模块
    NotificationModule,   // 引入通知模块
)
```

## 🔍 Module 的作用

### 1️⃣ **代码组织**（最主要）

```go
// 项目结构
project/
├── logging/
│   └── module.go       → var LoggingModule = fx.Module(...)
├── database/
│   └── module.go       → var DatabaseModule = fx.Module(...)
├── notification/
│   └── module.go       → var NotificationModule = fx.Module(...)
└── main.go
    └── fx.New(LoggingModule, DatabaseModule, NotificationModule)

// 每个模块独立，职责清晰
```

### 2️⃣ **可复用性**

```go
// 定义一个可复用的模块
var AuthModule = fx.Module("auth",
    fx.Provide(
        NewJWTService,
        NewAuthMiddleware,
        NewUserRepository,
    ),
)

// 在多个应用中使用
// app1/main.go
fx.New(AuthModule, ...)

// app2/main.go
fx.New(AuthModule, ...)
```

### 3️⃣ **命名空间隔离**（高级用法）

```go
// 两个模块可以有同名的类型
var Module1 = fx.Module("module1",
    fx.Provide(fx.Annotate(
        NewLogger,
        fx.ResultTags(`name:"module1-logger"`),
    )),
)

var Module2 = fx.Module("module2",
    fx.Provide(fx.Annotate(
        NewLogger,
        fx.ResultTags(`name:"module2-logger"`),
    )),
)
// 通过命名区分
```

## 🆚 对比：有无 Module

### 没有 Module（小项目可以）

```go
// main.go - 所有东西都在这
fx.New(
    fx.Provide(
        NewLogger,
        NewDB,
        NewCache,
        NewUserService,
        NewOrderService,
        NewPaymentService,
        NewEmailService,
        NewSMSService,
        // ... 30 个
    ),
)
// 😵 超过 20 个就开始混乱
```

### 使用 Module（推荐）

```go
// infrastructure/module.go
var InfraModule = fx.Module("infra",
    fx.Provide(NewLogger, NewDB, NewCache),
)

// user/module.go
var UserModule = fx.Module("user",
    fx.Provide(NewUserService, NewUserRepo),
)

// order/module.go
var OrderModule = fx.Module("order",
    fx.Provide(NewOrderService, NewOrderRepo),
)

// main.go
fx.New(
    InfraModule,
    UserModule,
    OrderModule,
    PaymentModule,
    NotificationModule,
)
// ✨ 清晰！
```

## 💡 实际项目示例

### 中型项目结构

```go
project/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── module.go       // ConfigModule
│   ├── database/
│   │   └── module.go       // DatabaseModule
│   ├── cache/
│   │   └── module.go       // CacheModule
│   ├── user/
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── module.go       // UserModule
│   ├── order/
│   │   └── module.go       // OrderModule
│   └── notification/
│       └── module.go       // NotificationModule

// cmd/main.go
func main() {
    fx.New(
        config.Module,
        database.Module,
        cache.Module,
        user.Module,
        order.Module,
        notification.Module,
    ).Run()
}
```

### 每个模块的定义

```go
// internal/user/module.go
package user

import "go.uber.org/fx"

var Module = fx.Module("user",
    fx.Provide(
        NewService,      // 用户服务
        NewRepository,   // 用户仓储
        NewHandler,      // HTTP 处理器
    ),
)

// internal/notification/module.go
package notification

var Module = fx.Module("notification",
    fx.Provide(
        NewEmailService,
        NewSMSService,
        NewPushService,
        NewNotificationService,
    ),
)
```

## 🔧 Module 的高级用法

### 1️⃣ 模块组合

```go
// 小模块
var EmailModule = fx.Module("email",
    fx.Provide(NewEmailSender),
)

var SMSModule = fx.Module("sms",
    fx.Provide(NewSMSSender),
)

// 大模块（组合小模块）
var NotificationModule = fx.Module("notification",
    fx.Options(
        EmailModule,      // 包含 Email 模块
        SMSModule,        // 包含 SMS 模块
        fx.Provide(NewNotificationService),
    ),
)
```

### 2️⃣ 模块配置

```go
// 可配置的模块
func NewDatabaseModule(config DBConfig) fx.Option {
    return fx.Module("database",
        fx.Supply(config),  // 提供配置
        fx.Provide(NewDatabase),
    )
}

// 使用
fx.New(
    NewDatabaseModule(DBConfig{Host: "localhost"}),
)
```

### 3️⃣ 条件模块

```go
func GetModules(env string) []fx.Option {
    modules := []fx.Option{
        CoreModule,
        DatabaseModule,
    }
    
    if env == "production" {
        modules = append(modules, MonitoringModule)
    } else {
        modules = append(modules, MockModule)
    }
    
    return modules
}

fx.New(GetModules(os.Getenv("ENV"))...)
```

## 📝 Module vs Provide 的区别

| 对比 | fx.Provide | fx.Module |
|------|-----------|-----------|
| **作用** | 注册单个构造函数 | 打包多个 Provider |
| **返回** | Option | Option |
| **嵌套** | 不能嵌套 | 可以包含其他 Module |
| **命名** | 无名字 | 有名字（调试用） |
| **适用** | 单个依赖 | 一组相关依赖 |

```go
// fx.Provide：单个
fx.Provide(NewLogger)

// fx.Module：一组
fx.Module("logging",
    fx.Provide(
        NewLogger,
        NewLogRotator,
        NewLogFormatter,
    ),
)
```

## ⚖️ 何时使用 Module？

### ✅ 应该用 Module

```
• 项目超过 10 个 Provider
• 多个开发者协作
• 需要复用某组功能
• 想要清晰的架构
```

### ❌ 不需要 Module

```
• 小脚本（< 5 个 Provider）
• 快速原型
• 学习阶段
```

## 💡 总结

| 概念 | 解释 |
|------|------|
| **fx.Module** | 把相关 Provider 打包成一个逻辑单元 |
| **核心价值** | 代码组织 + 可维护性 + 可复用 |
| **本质** | 就是 `fx.Options()` 的语法糖，加了个名字 |
| **类比** | Go 的 package，Node.js 的 module |

**一句话：Module 就是把一堆 Provider 打包，让代码不像一锅粥！**