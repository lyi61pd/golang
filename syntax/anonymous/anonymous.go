package main

import "fmt"

// 匿名函数：没有名字的函数，可以像变量一样赋值、传递和调用
// 闭包：匿名函数引用了其外部作用域的变量，变量会被“捕获”并随函数一起存在

func main() {
	// 1. 匿名函数赋值给变量
	// 这里定义了一个没有名字的函数，并赋值给变量 add
	add := func(a, b int) int {
		return a + b
	}
	// 现在可以像普通函数一样调用 add
	fmt.Println("add(2, 3) =", add(2, 3)) // 输出 5

	// 2. 匿名函数作为参数传递
	// do 是一个函数，接收另一个函数作为参数
	do := func(op func(int, int) int, x, y int) int {
		return op(x, y)
	}
	// 这里直接把匿名函数作为参数传递，实现乘法
	result := do(func(a, b int) int { return a * b }, 4, 5)
	fmt.Println("do(4*5) =", result) // 输出 20

	// 3. 匿名函数作为返回值（闭包）
	// makeAdder 返回一个函数，这个返回的函数会“记住”外部变量 x
	makeAdder := func(x int) func(int) int {
		// 这里返回的匿名函数引用了外部的 x
		return func(y int) int {
			return x + y // x 被“捕获”，形成闭包
		}
	}
	add10 := makeAdder(10)              // add10 是一个函数，每次调用都会把 10 加到参数上
	fmt.Println("add10(5) =", add10(5)) // 输出 15

	// 4. 闭包常见场景：累加器
	// acc 返回一个函数，这个函数内部有一个 sum 变量，每次调用都会累加
	acc := func() func(int) int {
		sum := 0 // sum 只会初始化一次，后续调用会记住它的值
		return func(x int) int {
			sum += x // sum 被闭包捕获并持续累加
			return sum
		}
	}
	accumulator := acc()
	fmt.Println("accumulator(1):", accumulator(1)) // 输出 1
	fmt.Println("accumulator(2):", accumulator(2)) // 输出 3
	fmt.Println("accumulator(3):", accumulator(3)) // 输出 6

	// 总结：
	// - 匿名函数可以像变量一样使用
	// - 闭包可以让函数记住并操作其外部作用域的变量
	// - 常用于回调、工厂函数、数据封装等场景
}
