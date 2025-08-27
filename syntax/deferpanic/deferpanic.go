package main

import "fmt"

// defer：用于延迟执行一段代码（通常用于资源释放、收尾工作）
// panic：用于主动抛出异常，使程序中断执行
// recover：用于捕获 panic，防止程序崩溃，实现异常处理

func main() {
	fmt.Println("程序开始")

	// defer 声明的语句会在 main 函数返回前（包括发生 panic 时）按“后进先出”顺序执行
	defer fmt.Println("defer 1：一定会执行（最后）")
	defer fmt.Println("defer 2：也会执行（先于 defer 1）")

	// 用于捕获 panic 的 defer
	defer func() {
		// recover 只能在 defer 函数中调用
		// 如果 main 函数中发生了 panic，这里可以捕获到
		if r := recover(); r != nil {
			fmt.Println("捕获到 panic：", r)
		}
	}()

	fmt.Println("即将触发 panic")
	panic("发生了一个严重错误") // 主动抛出异常，程序会中断到最近的 defer

	// 这行不会被执行，因为 panic 之后 main 函数会中断
	fmt.Println("程序结束")
}

/*
执行顺序说明（建议运行观察）：
1. 打印“程序开始”
2. defer 语句被注册（但暂不执行）
3. 打印“即将触发 panic”
4. panic 触发，程序中断，进入所有 defer
5. 先执行最后注册的 defer（defer 2），再执行 defer 1
6. defer func() 里的 recover 捕获到 panic，打印异常信息
7. 程序优雅退出，不会崩溃
*/
