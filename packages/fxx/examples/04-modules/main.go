package modules
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













































































































































}	fmt.Println("  • 模块让代码更清晰、可维护")	fmt.Println("  • 模块可以包含多个 Provider")	fmt.Println("  • fx.Module() 将相关功能组织在一起")	fmt.Println("\n💡 学习要点：")	app.Run()	)		}),			})				},					return nil					}()						shutdowner.Shutdown()						time.Sleep(500 * time.Millisecond)					go func() {					// 自动退出					notificationService.NotifyUser("欢迎使用我们的服务！")					// 使用通知服务					userService.GetUser("user789")					// 使用用户服务					fmt.Println("\n=== 应用程序运行中 ===")				OnStart: func(ctx context.Context) error {			lc.Append(fx.Hook{		) {			shutdowner fx.Shutdowner,			lc fx.Lifecycle,			notificationService *NotificationService,			userService *UserService,		fx.Invoke(func(		NotificationModule,		ServiceModule,		DataModule,		LoggingModule,		// 使用模块	app := fx.New(	fmt.Println("=== fx 模块化示例 ===\n")func main() {)	),		NewNotificationService,		NewSMSSender,		NewEmailSender,	fx.Provide(var NotificationModule = fx.Module("notification",)	fx.Provide(NewUserService),var ServiceModule = fx.Module("service",)	),		NewCache,		NewDatabase,	fx.Provide(var DataModule = fx.Module("data",)	fx.Provide(NewSimpleLogger),var LoggingModule = fx.Module("logging",// 定义模块}	n.smsSender.Send(message)	n.emailSender.Send(message)	n.logger.Log("\n--- 发送通知 ---")func (n *NotificationService) NotifyUser(message string) {}	}		logger:      logger,		smsSender:   sms,		emailSender: email,	return &NotificationService{	logger.Log("✓ NotificationService 已创建")func NewNotificationService(email *EmailSender, sms *SMSSender, logger Logger) *NotificationService {}	logger      Logger	smsSender   *SMSSender	emailSender *EmailSendertype NotificationService struct {// NotificationService 通知服务}	return nil	s.logger.Log(fmt.Sprintf("📱 发送短信: %s", message))func (s *SMSSender) Send(message string) error {}	return &SMSSender{logger: logger}	logger.Log("✓ SMSSender 已创建")func NewSMSSender(logger Logger) *SMSSender {}	logger Loggertype SMSSender struct {// SMSSender 短信发送器}	return nil	e.logger.Log(fmt.Sprintf("📧 发送邮件: %s", message))func (e *EmailSender) Send(message string) error {}	return &EmailSender{logger: logger}	logger.Log("✓ EmailSender 已创建")func NewEmailSender(logger Logger) *EmailSender {}	logger Loggertype EmailSender struct {// EmailSender 邮件发送器}	s.logger.Log(fmt.Sprintf("👤 获取用户: %s", id))func (s *UserService) GetUser(id string) {}	return &UserService{db: db, cache: cache, logger: logger}	logger.Log("✓ UserService 已创建")func NewUserService(db *Database, cache *Cache, logger Logger) *UserService {}	logger Logger	cache  *Cache	db     *Databasetype UserService struct {// UserService 用户服务}	return fmt.Sprintf("缓存值: %s", key)func (c *Cache) Get(key string) string {