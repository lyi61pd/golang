package main

import (
	"errors"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	println("=== Zap 日志库使用教程 ===\n")

	// 1. 基础使用 - 预设的 Logger
	println("1. 预设 Logger:")
	basicLoggers()

	// 2. 自定义配置
	println("\n2. 自定义配置:")
	customConfig()

	// 3. 结构化日志字段
	println("\n3. 结构化字段:")
	structuredLogging()

	// 4. 日志级别
	println("\n4. 日志级别:")
	logLevels()

	// 5. 性能优化 - SugaredLogger
	println("\n5. SugaredLogger:")
	sugaredLogging()

	// 6. 实际应用场景
	println("\n6. 实际应用:")
	practicalUsage()

	// 7. 日志输出到文件
	println("\n7. 日志输出配置:")
	fileOutput()
}

// 1. 基础使用 - Zap 提供了三种预设的 Logger
func basicLoggers() {
	// 开发环境 Logger - 输出可读性好，性能较低
	devLogger, _ := zap.NewDevelopment()
	defer devLogger.Sync()
	devLogger.Info("这是开发环境的日志",
		zap.String("environment", "development"),
		zap.Int("code", 200))

	// 生产环境 Logger - JSON 格式，性能高
	prodLogger, _ := zap.NewProduction()
	defer prodLogger.Sync()
	prodLogger.Info("这是生产环境的日志",
		zap.String("environment", "production"),
		zap.Int("code", 200))

	// 示例 Logger - 最简单的配置
	exampleLogger := zap.NewExample()
	defer exampleLogger.Sync()
	exampleLogger.Info("这是示例日志",
		zap.String("environment", "example"))
}

// 2. 自定义配置 Logger
func customConfig() {
	// 使用 Config 结构体自定义
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json", // 可选: "json" 或 "console"
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// 自定义时间格式
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := config.Build()
	defer logger.Sync()
	logger.Info("自定义配置的日志", zap.String("custom", "true"))
}

// 3. 结构化日志字段 - Zap 的核心优势
func structuredLogging() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 基本类型字段
	logger.Info("用户操作",
		zap.String("username", "张三"),
		zap.Int("user_id", 12345),
		zap.Bool("is_admin", true),
		zap.Float64("balance", 1234.56),
		zap.Duration("request_time", 123*time.Millisecond),
		zap.Time("created_at", time.Now()),
	)

	// 复杂类型字段
	logger.Info("复杂数据",
		zap.Strings("tags", []string{"golang", "zap", "logging"}),
		zap.Ints("scores", []int{90, 85, 92}),
		zap.Any("metadata", map[string]interface{}{
			"ip":     "192.168.1.1",
			"region": "cn-beijing",
		}),
	)

	// Error 字段
	err := errors.New("数据库连接失败")
	logger.Error("操作失败",
		zap.Error(err),
		zap.String("operation", "db_query"),
	)

	// 嵌套对象
	user := struct {
		Name  string
		Age   int
		Email string
	}{
		Name:  "李四",
		Age:   28,
		Email: "lisi@example.com",
	}
	logger.Info("用户信息", zap.Any("user", user))
}

// 4. 日志级别
func logLevels() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Zap 支持以下日志级别（从低到高）:
	logger.Debug("调试信息 - 最详细的日志")
	logger.Info("普通信息 - 重要的业务流程")
	logger.Warn("警告信息 - 需要注意但不影响运行")
	logger.Error("错误信息 - 出现错误但程序可以继续")

	// Fatal 和 Panic 会终止程序，谨慎使用
	// logger.Fatal("致命错误 - 记录后调用 os.Exit(1)")
	// logger.Panic("恐慌错误 - 记录后调用 panic()")

	// 动态调整日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.WarnLevel) // 只记录 Warn 及以上级别
	logger.Info("这条日志不会显示，因为级别低于 Warn")
}

// 5. SugaredLogger - 更灵活但性能稍低
func sugaredLogging() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// 转换为 SugaredLogger
	sugar := logger.Sugar()

	// 支持 printf 风格
	sugar.Infof("用户 %s 登录成功, ID: %d", "王五", 67890)

	// 支持键值对（自动配对）
	sugar.Infow("用户操作",
		"action", "login",
		"username", "王五",
		"ip", "192.168.1.100",
	)

	// 普通参数
	sugar.Info("这是一条简单日志")
	sugar.Warn("警告:", "磁盘空间不足")

	// 性能对比说明:
	// Logger: 最快，零内存分配，需要使用 zap.String() 等类型化字段
	// SugaredLogger: 稍慢，有少量内存分配，但使用更灵活方便
}

// 6. 实际应用场景
func practicalUsage() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 场景1: HTTP 请求日志
	logHTTPRequest(logger)

	// 场景2: 数据库操作日志
	logDatabaseOperation(logger)

	// 场景3: 使用 With 添加公共字段
	requestLogger := logger.With(
		zap.String("request_id", "req-12345"),
		zap.String("user_id", "user-789"),
	)
	requestLogger.Info("处理用户请求")
	requestLogger.Info("查询数据库")
	requestLogger.Info("返回响应")
}

func logHTTPRequest(logger *zap.Logger) {
	logger.Info("HTTP 请求",
		zap.String("method", "POST"),
		zap.String("path", "/api/users"),
		zap.Int("status_code", 201),
		zap.Duration("latency", 45*time.Millisecond),
		zap.String("client_ip", "192.168.1.1"),
		zap.String("user_agent", "Mozilla/5.0"),
	)
}

func logDatabaseOperation(logger *zap.Logger) {
	start := time.Now()
	// 模拟数据库操作
	time.Sleep(10 * time.Millisecond)
	duration := time.Since(start)

	logger.Info("数据库查询",
		zap.String("query", "SELECT * FROM users WHERE id = ?"),
		zap.Duration("duration", duration),
		zap.Int("rows_affected", 1),
		zap.Bool("cache_hit", false),
	)
}

// 7. 日志输出配置
func fileOutput() {
	// 方式1: 使用 Config
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"stdout", // 标准输出
		// 如需输出到文件，取消下面注释（需要先创建 logs 目录）
		// "./logs/app.log",
	}
	logger, _ := config.Build()
	defer logger.Sync()

	logger.Info("这条日志会输出到控制台")

	// 方式2: 自定义编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 使用 console 编码器（更易读）
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 创建 core - 使用标准输出
	core := zapcore.NewCore(
		consoleEncoder,
		zapcore.Lock(os.Stdout),
		zap.InfoLevel,
	)

	customLogger := zap.New(core, zap.AddCaller())
	defer customLogger.Sync()

	customLogger.Info("自定义编码器的日志",
		zap.String("format", "console"),
		zap.Bool("caller_enabled", true))
}

// 额外提示：
// 1. 在生产环境中，建议使用 zap.NewProduction()
// 2. 使用 defer logger.Sync() 确保日志缓冲区被刷新
// 3. 对于高性能要求，使用 Logger 而不是 SugaredLogger
// 4. 使用 With() 创建带公共字段的子 Logger
// 5. 避免在热路径中使用 Any() 和反射操作
// 6. 考虑日志轮转（使用 lumberjack 库）
