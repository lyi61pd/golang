package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

// Logger 接口
type Logger interface {
	Log(message string)
}

type SimpleLogger struct{}

func NewSimpleLogger() Logger {
	return &SimpleLogger{}
}

func (l *SimpleLogger) Log(message string) {
	fmt.Printf("[LOG] %s\n", message)
}

// === 1. 可选依赖示例 ===

type Cache struct {
	logger Logger
}

func NewCache(logger Logger) *Cache {
	logger.Log("✓ Cache 已创建")
	return &Cache{logger: logger}
}

func (c *Cache) Get(key string) string {
	return fmt.Sprintf("缓存值: %s", key)
}

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

// OptionalParams 可选依赖参数
type OptionalParams struct {
	fx.In

	Logger        Logger
	Cache         *Cache         `optional:"true"` // 可选
	MetricsClient *MetricsClient `optional:"true"` // 可选
}

type ServiceWithOptional struct {
	logger        Logger
	cache         *Cache
	metricsClient *MetricsClient
}

func NewServiceWithOptional(params OptionalParams) *ServiceWithOptional {
	params.Logger.Log("✓ ServiceWithOptional 已创建")

	if params.Cache != nil {
		params.Logger.Log("  ✓ 检测到 Cache 依赖")
	} else {
		params.Logger.Log("  ⚠️  Cache 依赖不存在（可选）")
	}

	if params.MetricsClient != nil {
		params.Logger.Log("  ✓ 检测到 MetricsClient 依赖")
	} else {
		params.Logger.Log("  ⚠️  MetricsClient 依赖不存在（可选）")
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
		result := s.cache.Get("some-key")
		s.logger.Log(result)
	}

	if s.metricsClient != nil {
		s.metricsClient.RecordMetric("work.completed", 1.0)
	}
}

// === 2. 值组（Value Groups）示例 ===

type Handler interface {
	Name() string
	Handle()
}

type HandlerA struct {
	logger Logger
}

func (h *HandlerA) Name() string { return "HandlerA" }
func (h *HandlerA) Handle() {
	h.logger.Log("  ➤ HandlerA 正在处理")
}

type HandlerB struct {
	logger Logger
}

func (h *HandlerB) Name() string { return "HandlerB" }
func (h *HandlerB) Handle() {
	h.logger.Log("  ➤ HandlerB 正在处理")
}

// HandlersResult 使用值组
type HandlersResult struct {
	fx.Out

	HandlerA Handler `group:"handlers"`
	HandlerB Handler `group:"handlers"`
}

func NewHandlers(logger Logger) HandlersResult {
	logger.Log("✓ 创建处理器组")
	return HandlersResult{
		HandlerA: &HandlerA{logger: logger},
		HandlerB: &HandlerB{logger: logger},
	}
}

// HandlersParam 接收值组
type HandlersParam struct {
	fx.In

	Handlers []Handler `group:"handlers"`
}

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

func main() {
	fmt.Println("=== fx 高级特性示例 ===\n")

	app := fx.New(
		fx.Provide(
			NewSimpleLogger,

			// 提供可选依赖
			NewCache,
			NewMetricsClient,

			// 服务
			NewServiceWithOptional,

			// 值组
			NewHandlers,
			NewRouter,
		),
		fx.Invoke(func(
			service *ServiceWithOptional,
			router *Router,
			lc fx.Lifecycle,
			shutdowner fx.Shutdowner,
		) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					fmt.Println("\n=== 应用程序运行中 ===")

					// 测试可选依赖
					service.DoWork()

					// 测试值组
					router.RouteAll()

					// 自动退出
					go func() {
						time.Sleep(500 * time.Millisecond)
						shutdowner.Shutdown()
					}()
					return nil
				},
			})
		}),
	)

	app.Run()

	fmt.Println("\n💡 学习要点：")
	fmt.Println("  • optional:\"true\" 标记可选依赖")
	fmt.Println("  • group:\"name\" 收集同类型的多个实现")
	fmt.Println("  • fx.In/fx.Out 结构体用于复杂参数传递")
}
