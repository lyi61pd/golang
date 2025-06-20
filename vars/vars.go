package main

import "fmt"

// 本例演示变量声明和常量定义的多种方式
func main() {
	// 变量声明
	var a int = 100 // 显式声明类型为 int，赋值为 100
	var b = 3.14    // 自动类型推断为 float64
	c := "Hello Go" // 简短声明，自动推断为 string

	// 常量定义
	const pi = 3.1415            // 常量，类型自动推断为 float64
	const name string = "Gopher" // 常量，显式指定类型为 string

	// 输出变量和常量的值
	fmt.Println("a =", a)
	fmt.Println("b =", b)
	fmt.Println("c =", c)
	fmt.Println("pi =", pi)
	fmt.Println("name =", name)
}
