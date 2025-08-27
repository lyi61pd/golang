package main

import (
	"fmt"
	"time"
)

// goroutine 是 Go 的并发执行单元，使用 go 关键字启动。
// channel 是 goroutine 之间通信的管道，可以安全地在多个 goroutine 之间传递数据。

// 定义一个函数，稍后用 goroutine 启动
func sayHello() {
	fmt.Println("你好，世界！") // 这行会在新的 goroutine 中执行
}

func main() {
	// 启动一个新的 goroutine，异步执行 sayHello
	go sayHello()

	// 主 goroutine 继续执行，如果不等待，主程序可能提前退出
	time.Sleep(time.Second) // 等待1秒，确保上面的 goroutine 有机会运行

	// channel 示例
	ch := make(chan int) // 创建一个整型 channel

	// 启动一个 goroutine，向 channel 发送数据
	go func() {
		ch <- 42 // 发送数字42到 channel
	}()

	// 从 channel 接收数据（会阻塞直到有数据到来）
	value := <-ch
	fmt.Println("从channel收到：", value)
}
