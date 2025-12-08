package main

import (
	"fmt"

	"go.uber.org/fx"
)

// 4. 模块化示例：使用 fx.Module 组织代码

// LoggingModule 日志模块
var LoggingModule = fx.Module("logging",
	fx.Provide(NewSimpleLogger),
)

// DataModule 数据访问模块
var DataModule = fx.Module("data",
	fx.Provide(
		NewDatabase,
		NewCache,
	),
)

// ServiceModule 业务服务模块
var ServiceModule = fx.Module("service",
	fx.Provide(NewUserService),
)

// 5. 使用接口和多实现

// MessageSender 消息发送接口
type MessageSender interface {
	Send(message string) error
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

func (e *SMSSender) Send(message string) error {
	e.logger.Log(fmt.Sprintf("📱 发送短信: %s", message))
	return nil
}

// NotificationService 通知服务，依赖多个发送器
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

// NotificationModule 通知模块
var NotificationModule = fx.Module("notification",
	fx.Provide(
		NewEmailSender,
		NewSMSSender,
		NewNotificationService,
	),
)

// 要运行此示例：
// 创建新目录并复制必要的代码，然后取消注释下面的 main 函数

/*
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
		) {
			fmt.Println("\n=== 应用程序运行中 ===")

			// 使用用户服务
			userService.GetUser("user789")

			// 使用通知服务
			notificationService.NotifyUser("欢迎使用我们的服务！")
		}),
	)

	app.Run()
}
*/
