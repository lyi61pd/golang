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

// Database 模拟数据库
type Database struct {
	logger Logger
}

func NewDatabase(logger Logger) *Database {
	logger.Log("✓ Database 已创建")
	return &Database{logger: logger}
}

func (db *Database) Query(sql string) string {
	return "查询结果"
}

// Cache 模拟缓存
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

// UserService 用户服务
type UserService struct {
	db     *Database
	cache  *Cache
	logger Logger
}

func NewUserService(db *Database, cache *Cache, logger Logger) *UserService {
	logger.Log("✓ UserService 已创建")
	return &UserService{db: db, cache: cache, logger: logger}
}

func (s *UserService) GetUser(id string) {
	s.logger.Log(fmt.Sprintf("👤 获取用户: %s", id))
}

// EmailSender 邮件发送器
type EmailSender struct {
	logger Logger
}

func NewEmailSender(logger Logger) *EmailSender {
	logger.Log("✓ EmailSender 已创建")
	return &EmailSender{logger: logger}
}

func (e *EmailSender) Send(message string) error {
	e.logger.Log(fmt.Sprintf("📧 发送邮件: %s", message))
	return nil
}

// SMSSender 短信发送器
type SMSSender struct {
	logger Logger
}

func NewSMSSender(logger Logger) *SMSSender {
	logger.Log("✓ SMSSender 已创建")
	return &SMSSender{logger: logger}
}

func (s *SMSSender) Send(message string) error {
	s.logger.Log(fmt.Sprintf("📱 发送短信: %s", message))
	return nil
}

// NotificationService 通知服务
type NotificationService struct {
	emailSender *EmailSender
	smsSender   *SMSSender
	logger      Logger
}

func NewNotificationService(email *EmailSender, sms *SMSSender, logger Logger) *NotificationService {
	logger.Log("✓ NotificationService 已创建")
	return &NotificationService{
		emailSender: email,
		smsSender:   sms,
		logger:      logger,
	}
}

func (n *NotificationService) NotifyUser(message string) {
	n.logger.Log("\n--- 发送通知 ---")
	n.emailSender.Send(message)
	n.smsSender.Send(message)
}

// 定义模块
var LoggingModule = fx.Module("logging",
	fx.Provide(NewSimpleLogger),
)

var DataModule = fx.Module("data",
	fx.Provide(
		NewDatabase,
		NewCache,
	),
)

var ServiceModule = fx.Module("service",
	fx.Provide(NewUserService),
)

var NotificationModule = fx.Module("notification",
	fx.Provide(
		NewEmailSender,
		NewSMSSender,
		NewNotificationService,
	),
)

func main() {
	fmt.Println("=== fx 模块化示例 ===\n")

	app := fx.New(
		// 使用模块
		LoggingModule,
		DataModule,
		ServiceModule,
		NotificationModule,

		fx.Invoke(func(
			userService *UserService,
			notificationService *NotificationService,
			lc fx.Lifecycle,
			shutdowner fx.Shutdowner,
		) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					fmt.Println("\n=== 应用程序运行中 ===")

					// 使用用户服务
					userService.GetUser("user789")

					// 使用通知服务
					notificationService.NotifyUser("欢迎使用我们的服务！")

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
	fmt.Println("  • fx.Module() 将相关功能组织在一起")
	fmt.Println("  • 模块可以包含多个 Provider")
	fmt.Println("  • 模块让代码更清晰、可维护")
}
