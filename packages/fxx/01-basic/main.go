package main

import (
	"fmt"

	"go.uber.org/fx"
)

// 1. 基础示例：简单的依赖注入

// Logger 是一个简单的日志接口
type Logger interface {
	Log(message string)
}

// SimpleLogger 实现了 Logger 接口
type SimpleLogger struct{}

func NewSimpleLogger() Logger {
	fmt.Println("✓ SimpleLogger 已创建")
	return &SimpleLogger{}
}

func (l *SimpleLogger) Log(message string) {
	fmt.Printf("[LOG] %s\n", message)
}

// Greeter 依赖于 Logger
type Greeter struct {
	logger Logger
}

func NewGreeter(logger Logger) *Greeter {
	fmt.Println("✓ Greeter 已创建")
	return &Greeter{logger: logger}
}

func (g *Greeter) Greet(name string) {
	message := fmt.Sprintf("Hello, %s!", name)
	g.logger.Log(message)
}

// Application 是主应用程序
type Application struct {
	greeter *Greeter
}

func NewApplication(greeter *Greeter) *Application {
	fmt.Println("✓ Application 已创建")
	return &Application{greeter: greeter}
}

func (a *Application) Run() {
	fmt.Println("\n=== 应用程序开始运行 ===")
	a.greeter.Greet("fx 学习者")
	a.greeter.Greet("Go 开发者")
	fmt.Println("=== 应用程序运行完成 ===\n")
}

func main() {
	fmt.Println("=== fx 依赖注入框架基础示例 ===\n")
	fmt.Println("--- 依赖构建阶段 ---")

	app := fx.New(
		// 提供依赖
		fx.Provide(
			NewSimpleLogger, // 提供 Logger
			NewGreeter,      // 提供 Greeter（自动注入 Logger）
			NewApplication,  // 提供 Application（自动注入 Greeter）
		),
		// 调用应用程序的 Run 方法
		fx.Invoke(func(app *Application, shutdowner fx.Shutdowner) {
			app.Run()
			// 程序运行完成后自动关闭
			shutdowner.Shutdown()
		}),
	)

	fmt.Println("\n--- 应用程序启动阶段 ---")
	app.Run()

	fmt.Println("\n💡 学习要点：")
	fmt.Println("  • fx.Provide() 注册构造函数")
	fmt.Println("  • fx 自动解析依赖：Logger → Greeter → Application")
	fmt.Println("  • fx.Invoke() 执行业务逻辑")
	fmt.Println("  • shutdowner.Shutdown() 让程序自动退出")
}
