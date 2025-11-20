package main

import (
	"context"
	"fmt"
	"time"
)

/*
Context 包详解

Context 是 Go 语言中用于在 goroutine 之间传递截止时间、取消信号和请求范围值的标准方式。
主要用于控制并发操作的生命周期。

主要用途：
1. 取消信号传递：通知 goroutine 停止工作
2. 超时控制：设置操作的最大执行时间
3. 传递请求范围的值：如请求 ID、用户信息等

常用的 Context 类型：
1. context.Background()    - 根 context，通常用于 main 函数、初始化和测试
2. context.TODO()          - 当不确定使用什么 context 时使用
3. context.WithCancel()    - 可以手动取消的 context
4. context.WithTimeout()   - 超时自动取消的 context
5. context.WithDeadline()  - 指定截止时间的 context
6. context.WithValue()     - 携带键值对数据的 context
*/

// 定义 context key 类型（包级别，避免不同函数中类型不一致）
type contextKey string

// 定义常量 key
const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
)

func main() {
	fmt.Println("=== Context 学习示例 ===")
	fmt.Println()

	// 示例 1：WithCancel - 手动取消
	fmt.Println("1. WithCancel 示例：手动取消操作")
	example1WithCancel()
	time.Sleep(time.Second)

	// 示例 2：WithTimeout - 超时控制
	fmt.Println("\n2. WithTimeout 示例：超时控制")
	example2WithTimeout()
	time.Sleep(time.Second)

	// 示例 3：WithDeadline - 截止时间
	fmt.Println("\n3. WithDeadline 示例：截止时间控制")
	example3WithDeadline()
	time.Sleep(time.Second)

	// 示例 4：WithValue - 传递请求范围的值
	fmt.Println("\n4. WithValue 示例：传递请求范围的值")
	example4WithValue()
	time.Sleep(time.Second)

	// 示例 5：级联取消
	fmt.Println("\n5. 级联取消示例：父 context 取消会影响所有子 context")
	example5CascadeCancel()
	time.Sleep(time.Second)

	// 示例 6：实际应用 - 模拟 HTTP 请求处理
	fmt.Println("\n6. 实际应用：模拟 HTTP 请求处理")
	example6HTTPRequest()
}

// 示例 1：WithCancel - 手动取消
func example1WithCancel() {
	// 创建一个可以取消的 context
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// 模拟长时间运行的任务
		for {
			select {
			case <-ctx.Done():
				// context 被取消时，Done() 会返回
				fmt.Println("   goroutine 收到取消信号，退出")
				return
			default:
				fmt.Println("   goroutine 正在工作...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// 主 goroutine 等待 2 秒后取消
	time.Sleep(2 * time.Second)
	fmt.Println("   主程序调用 cancel()，取消操作")
	cancel() // 手动取消
	time.Sleep(500 * time.Millisecond)
}

// 示例 2：WithTimeout - 超时控制
func example2WithTimeout() {
	// 创建一个 3 秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // 养成好习惯：即使超时也要调用 cancel 释放资源

	// 模拟一个可能耗时的操作
	result := make(chan string)
	go func() {
		// 模拟耗时操作（4 秒）
		time.Sleep(4 * time.Second)
		result <- "操作完成"
	}()

	select {
	case <-ctx.Done():
		// 超时或被取消
		fmt.Printf("   操作超时：%v\n", ctx.Err())
	case res := <-result:
		fmt.Printf("   %s\n", res)
	}
}

// 示例 3：WithDeadline - 截止时间
func example3WithDeadline() {
	// 设置 2 秒后的截止时间
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go processWithDeadline(ctx)

	time.Sleep(3 * time.Second)
}

func processWithDeadline(ctx context.Context) {
	// 检查截止时间
	if deadline, ok := ctx.Deadline(); ok {
		fmt.Printf("   任务截止时间：%v\n", deadline.Format("15:04:05"))
	}

	select {
	case <-time.After(3 * time.Second):
		fmt.Println("   任务完成")
	case <-ctx.Done():
		fmt.Printf("   任务因超过截止时间而终止：%v\n", ctx.Err())
	}
}

// 示例 4：WithValue - 传递请求范围的值
func example4WithValue() {
	// 创建携带值的 context（使用包级别定义的 contextKey 类型）
	ctx := context.WithValue(context.Background(), userIDKey, "12345")
	ctx = context.WithValue(ctx, requestIDKey, "req-abc-123")

	// 在函数中读取 context 中的值
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	// 从 context 中获取值（使用安全的类型断言）
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		fmt.Printf("   处理用户请求，userID: %s\n", userID)
	}

	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Printf("   请求 ID: %s\n", requestID)
	}

	// 调用其他函数，传递 context
	logRequest(ctx)
}

func logRequest(ctx context.Context) {
	// 使用安全的类型断言，避免 panic
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Printf("   [日志] 记录请求：%s\n", requestID)
	} else {
		fmt.Println("   [日志] 无法获取请求 ID")
	}
}

// 示例 5：级联取消
func example5CascadeCancel() {
	// 创建父 context
	parentCtx, parentCancel := context.WithCancel(context.Background())

	// 创建子 context
	childCtx1, childCancel1 := context.WithCancel(parentCtx)
	defer childCancel1() // 确保释放资源
	childCtx2, childCancel2 := context.WithCancel(parentCtx)
	defer childCancel2() // 确保释放资源

	// 启动多个 goroutine
	go worker(childCtx1, "Worker-1")
	go worker(childCtx2, "Worker-2")

	// 2 秒后取消父 context
	time.Sleep(2 * time.Second)
	fmt.Println("   取消父 context，所有子 context 也会被取消")
	parentCancel()

	time.Sleep(500 * time.Millisecond)
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("   %s 收到取消信号，退出\n", name)
			return
		default:
			fmt.Printf("   %s 正在工作...\n", name)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// 示例 6：实际应用 - 模拟 HTTP 请求处理
func example6HTTPRequest() {
	// 模拟 HTTP 请求，设置 5 秒超时
	ctx := context.WithValue(context.Background(), requestIDKey, "req-001")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fmt.Println("   开始处理 HTTP 请求...")

	// 模拟调用数据库
	if err := queryDatabase(ctx); err != nil {
		fmt.Printf("   数据库查询失败：%v\n", err)
		return
	}

	// 模拟调用外部 API
	if err := callExternalAPI(ctx); err != nil {
		fmt.Printf("   调用外部 API 失败：%v\n", err)
		return
	}

	fmt.Println("   HTTP 请求处理完成")
}

func queryDatabase(ctx context.Context) error {
	// 模拟数据库查询（2 秒）
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("   数据库查询成功")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("数据库查询被取消：%w", ctx.Err())
	}
}

func callExternalAPI(ctx context.Context) error {
	// 模拟外部 API 调用（1 秒）
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("   外部 API 调用成功")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("外部 API 调用被取消：%w", ctx.Err())
	}
}

/*
Context 最佳实践：

1. Context 应该作为函数的第一个参数，命名为 ctx
   func DoSomething(ctx context.Context, arg string) error

2. 不要将 Context 存储在结构体中，应该显式传递

3. 不要传递 nil Context，如果不确定用什么，使用 context.TODO()

4. Context.Value 仅用于传递请求范围的数据，不要用于传递可选参数
   ⚠️ 重要：使用 WithValue 时，key 的类型必须在包级别定义
   在不同函数中重复定义相同名称的类型，它们是不同的类型，会导致取值失败

5. 使用类型断言时，务必使用逗号 ok 模式，避免 panic
   正确：if val, ok := ctx.Value(key).(string); ok { ... }
   错误：val := ctx.Value(key).(string)  // 如果 key 不存在会 panic

6. Context 是并发安全的，可以在多个 goroutine 中安全使用

7. 调用 WithCancel、WithTimeout、WithDeadline 返回的 cancel 函数必须被调用
   使用 defer cancel() 确保资源释放

常见错误类型：
- context.Canceled：context 被手动取消
- context.DeadlineExceeded：超过截止时间

常见陷阱：
- 在不同作用域重复定义 contextKey 类型，导致取值失败
- 直接类型断言不检查 ok，导致 panic
- 忘记调用 cancel 函数，导致资源泄漏
*/
