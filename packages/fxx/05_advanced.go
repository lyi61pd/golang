package main

import (
	"fmt"

	"go.uber.org/fx"
)

// 5. 高级特性示例

// === 1. 使用结构体参数（避免参数过多）===

// ServerParams 服务器参数
type ServerParams struct {
	fx.In // 标记这是一个依赖注入参数结构体

	Logger  Logger
	Config  *Config
	Handler *HTTPHandler
}

func NewServerWithParams(params ServerParams) *HTTPServer {
	params.Logger.Log("使用结构体参数创建服务器")
	return &HTTPServer{
		logger: params.Logger,
	}
}

// === 2. 使用结果结构体（提供多个返回值）===

// LoggerResult 日志器结果
type LoggerResult struct {
	fx.Out // 标记这是一个依赖注入结果结构体

	Logger        Logger
	ConsoleLogger Logger `name:"console"` // 带名称的依赖
}

func NewLoggers() LoggerResult {
	return LoggerResult{
		Logger:        NewSimpleLogger(),
		ConsoleLogger: NewSimpleLogger(),
	}
}

// === 3. 可选依赖 ===

// OptionalParams 可选依赖参数
type OptionalParams struct {
	fx.In

	Logger        Logger
	Cache         *Cache         `optional:"true"` // 标记为可选
	MetricsClient *MetricsClient `optional:"true"`
}

// MetricsClient 指标客户端
type MetricsClient struct {
	logger Logger
}

func NewMetricsClient(logger Logger) *MetricsClient {
	logger.Log("✓ MetricsClient 已创建")
	return &MetricsClient{logger: logger}
}

func (m *MetricsClient) RecordMetric(name string, value float64) {
	m.logger.Log(fmt.Sprintf("📊 记录指标: %s = %.2f", name, value))
}

// ServiceWithOptional 带可选依赖的服务
type ServiceWithOptional struct {
	logger        Logger
	cache         *Cache
	metricsClient *MetricsClient
}

func NewServiceWithOptional(params OptionalParams) *ServiceWithOptional {
	params.Logger.Log("✓ ServiceWithOptional 已创建")

	if params.Cache != nil {
		params.Logger.Log("  - 检测到 Cache 依赖")
	} else {
		params.Logger.Log("  - Cache 依赖不存在（可选）")
	}

	if params.MetricsClient != nil {
		params.Logger.Log("  - 检测到 MetricsClient 依赖")
	} else {
		params.Logger.Log("  - MetricsClient 依赖不存在（可选）")
	}

	return &ServiceWithOptional{
		logger:        params.Logger,
		cache:         params.Cache,
		metricsClient: params.MetricsClient,
	}
}

func (s *ServiceWithOptional) DoWork() {
	s.logger.Log("\n--- 执行工作 ---")

	if s.cache != nil {
		s.cache.Get("some-key")
	}

	if s.metricsClient != nil {
		s.metricsClient.RecordMetric("work.completed", 1.0)
	}
}

// === 4. 值组（Value Groups）===

// Handler 通用处理器接口
type Handler interface {
	Name() string
	Handle()
}

// HandlerA 处理器A
type HandlerA struct {
	logger Logger
}

func NewHandlerA(logger Logger) Handler {
	logger.Log("✓ HandlerA 已创建")
	return &HandlerA{logger: logger}
}

func (h *HandlerA) Name() string { return "HandlerA" }
func (h *HandlerA) Handle() {
	h.logger.Log("HandlerA 正在处理")
}

// HandlerB 处理器B
type HandlerB struct {
	logger Logger
}

func NewHandlerB(logger Logger) Handler {
	logger.Log("✓ HandlerB 已创建")
	return &HandlerB{logger: logger}
}

func (h *HandlerB) Name() string { return "HandlerB" }
func (h *HandlerB) Handle() {
	h.logger.Log("HandlerB 正在处理")
}

// HandlersResult 处理器结果（使用值组）
type HandlersResult struct {
	fx.Out

	HandlerA Handler `group:"handlers"` // 加入 handlers 组
	HandlerB Handler `group:"handlers"` // 加入 handlers 组
}

func NewHandlers(logger Logger) HandlersResult {
	return HandlersResult{
		HandlerA: NewHandlerA(logger),
		HandlerB: NewHandlerB(logger),
	}
}

// HandlersParam 处理器参数（接收值组）
type HandlersParam struct {
	fx.In

	Handlers []Handler `group:"handlers"` // 接收 handlers 组的所有值
}

// Router 路由器，处理所有处理器
type Router struct {
	handlers []Handler
	logger   Logger
}

func NewRouter(params HandlersParam, logger Logger) *Router {
	logger.Log(fmt.Sprintf("✓ Router 已创建，注册了 %d 个处理器", len(params.Handlers)))
	return &Router{
		handlers: params.Handlers,
		logger:   logger,
	}
}

func (r *Router) RouteAll() {
	r.logger.Log("\n--- 路由所有处理器 ---")
	for _, handler := range r.handlers {
		r.logger.Log(fmt.Sprintf("路由到: %s", handler.Name()))
		handler.Handle()
	}
}

// 要运行此示例：
// 创建新目录并复制必要的代码，然后取消注释下面的 main 函数

/*
func main() {
	fmt.Println("=== fx 高级特性示例 ===\n")

	app := fx.New(
		fx.Provide(
			NewSimpleLogger,
			NewConfig,

			// 提供可选依赖
			NewCache,
			NewMetricsClient,

			// 使用结构体参数
			NewServiceWithOptional,

			// 使用值组
			NewHandlers,
			NewRouter,
		),
		fx.Invoke(func(
			service *ServiceWithOptional,
			router *Router,
		) {
			fmt.Println("\n=== 应用程序运行中 ===")

			// 测试可选依赖
			service.DoWork()

			// 测试值组
			router.RouteAll()
		}),
	)

	app.Run()
}
*/
