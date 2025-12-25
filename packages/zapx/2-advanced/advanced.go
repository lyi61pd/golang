package main

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 实战技巧和最佳实践

func main() {
	println("=== Zap 实战技巧 ===\n")

	// 1. 全局 Logger 模式
	println("1. 全局 Logger:")
	globalLoggerPattern()

	// 2. 日志分级输出
	println("\n2. 日志分级输出:")
	levelBasedOutput()

	// 3. 动态调整日志级别
	println("\n3. 动态调整级别:")
	dynamicLevel()

	// 4. 添加上下文字段
	println("\n4. 上下文字段:")
	contextFields()

	// 5. 性能优化技巧
	println("\n5. 性能对比:")
	performanceComparison()

	// 6. 日志采样（高频场景）
	println("\n6. 日志采样:")
	logSampling()

	// 7. 错误处理最佳实践
	println("\n7. 错误处理:")
	errorHandling()
}

// 1. 全局 Logger 模式（推荐用于实际项目）
var globalLogger *zap.Logger

func InitLogger() {
	var err error
	globalLogger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func GetLogger() *zap.Logger {
	if globalLogger == nil {
		InitLogger()
	}
	return globalLogger
}

func globalLoggerPattern() {
	InitLogger()
	defer globalLogger.Sync()

	// 在任何地方使用
	GetLogger().Info("使用全局 Logger",
		zap.String("pattern", "singleton"))

	// 或者创建带上下文的子 Logger
	requestLogger := GetLogger().With(
		zap.String("trace_id", "abc-123"),
		zap.String("user_id", "user-456"),
	)

	requestLogger.Info("处理请求")
	requestLogger.Info("查询数据")
	// 所有日志都会自动带上 trace_id 和 user_id
}

// 2. 日志分级输出（Info 到控制台，Error 到文件）
func levelBasedOutput() {
	// 配置：Info 及以上输出到 stdout，Error 及以上输出到文件
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	// 控制台输出（彩色）
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// 文件输出（JSON）
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// 创建 multi core
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), lowPriority),
		zapcore.NewCore(fileEncoder, zapcore.Lock(os.Stdout), highPriority), // 这里用 stdout 演示，实际应该用文件
	)

	logger := zap.New(core)
	defer logger.Sync()

	logger.Info("这是 Info 日志 - 只输出到控制台")
	logger.Warn("这是 Warn 日志 - 只输出到控制台")
	logger.Error("这是 Error 日志 - 输出到控制台和文件")
}

// 3. 动态调整日志级别（不重启服务）
func dynamicLevel() {
	// 创建可动态调整的级别
	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	config := zap.Config{
		Level:            atomicLevel,
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	defer logger.Sync()

	logger.Debug("这条 Debug 不会显示")
	logger.Info("这条 Info 会显示")

	// 动态调整为 Debug 级别
	atomicLevel.SetLevel(zap.DebugLevel)
	fmt.Println("→ 已将日志级别调整为 Debug")

	logger.Debug("现在 Debug 会显示了")
	logger.Info("Info 依然显示")

	// 在实际项目中，可以通过 HTTP 接口动态调整级别
	// http.HandleFunc("/log/level", func(w http.ResponseWriter, r *http.Request) {
	//     atomicLevel.SetLevel(zap.DebugLevel)
	// })
}

// 4. 添加上下文字段（请求 ID、用户 ID 等）
func contextFields() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 方式1: 使用 With 创建子 Logger
	requestLogger := logger.With(
		zap.String("request_id", "req-123456"),
		zap.String("method", "POST"),
		zap.String("path", "/api/users"),
	)

	// 所有后续日志都会自动带上这些字段
	requestLogger.Info("开始处理请求")
	requestLogger.Info("验证用户权限")
	requestLogger.Info("查询数据库")
	requestLogger.Info("返回响应", zap.Int("status_code", 200))

	// 方式2: 嵌套添加更多上下文
	userLogger := requestLogger.With(
		zap.String("user_id", "user-789"),
		zap.String("role", "admin"),
	)

	userLogger.Info("执行管理员操作")
}

// 5. 性能对比：Logger vs SugaredLogger
func performanceComparison() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	iterations := 10000

	// 测试 Logger 性能
	start := time.Now()
	for i := 0; i < iterations; i++ {
		logger.Info("性能测试",
			zap.String("type", "logger"),
			zap.Int("iteration", i))
	}
	loggerDuration := time.Since(start)

	// 测试 SugaredLogger 性能
	start = time.Now()
	for i := 0; i < iterations; i++ {
		sugar.Infow("性能测试",
			"type", "sugared",
			"iteration", i)
	}
	sugarDuration := time.Since(start)

	fmt.Printf("Logger 耗时: %v\n", loggerDuration)
	fmt.Printf("SugaredLogger 耗时: %v\n", sugarDuration)
	fmt.Printf("性能差距: %.2fx\n", float64(sugarDuration)/float64(loggerDuration))
}

// 6. 日志采样（高频场景下避免日志爆炸）
func logSampling() {
	// 每秒只记录前 3 条，之后每 5 条记录 1 条
	samplingConfig := &zap.SamplingConfig{
		Initial:    3, // 每秒前 3 条
		Thereafter: 5, // 之后每 5 条记录 1 条
	}

	config := zap.NewProductionConfig()
	config.Sampling = samplingConfig

	logger, _ := config.Build()
	defer logger.Sync()

	fmt.Println("模拟高频日志（20条），但只会记录少数几条：")
	for i := 0; i < 20; i++ {
		logger.Info("高频日志",
			zap.Int("count", i),
			zap.String("message", "这是一条可能频繁出现的日志"))
	}
}

// 7. 错误处理最佳实践
func errorHandling() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 场景1: 业务错误
	err := fmt.Errorf("用户不存在")
	logger.Warn("业务错误",
		zap.Error(err),
		zap.String("user_id", "123"),
		zap.String("action", "query_user"))

	// 场景2: 系统错误（需要堆栈信息）
	systemErr := fmt.Errorf("数据库连接失败")
	logger.Error("系统错误",
		zap.Error(systemErr),
		zap.Stack("stacktrace")) // 添加堆栈跟踪

	// 场景3: 嵌套错误
	rootErr := fmt.Errorf("网络超时")
	wrappedErr := fmt.Errorf("执行查询失败: %w", rootErr)
	logger.Error("操作失败",
		zap.Error(wrappedErr),
		zap.NamedError("root_cause", rootErr)) // 同时记录根因

	// 场景4: 错误恢复
	defer func() {
		if r := recover(); r != nil {
			logger.Error("程序 panic",
				zap.Any("panic", r),
				zap.Stack("stacktrace"))
		}
	}()
}

/*
=== 实战建议 ===

1. 项目初始化：
   - 创建全局 Logger 单例
   - 根据环境（dev/prod）选择不同配置
   - 配置日志文件轮转（使用 lumberjack）

2. 日志级别使用指南：
   - Debug: 开发调试信息
   - Info: 重要业务流程（用户登录、订单创建等）
   - Warn: 可恢复的异常（重试成功、降级处理等）
   - Error: 需要关注的错误（数据库错误、第三方 API 失败等）
   - Fatal/Panic: 程序无法继续运行

3. 性能优化：
   - 热路径使用 Logger 而非 SugaredLogger
   - 高频日志使用采样
   - 避免在日志中使用 Any() 和复杂对象
   - 使用 Check 跳过不必要的日志

4. 结构化日志：
   - 始终使用类型化字段（zap.String, zap.Int 等）
   - 保持字段命名一致性
   - 为重要流程添加 trace_id
   - 使用 With() 避免重复字段

5. 监控告警：
   - Error 级别日志应该触发告警
   - 统计日志量异常增长
   - 关键业务流程添加打点
*/
