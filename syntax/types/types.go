package main

import "fmt"

// 本例演示 Go 的基本数据类型：int, float64, string, bool
func main() {
	var a int = 10          // 整型变量，赋值为10
	var b float64 = 3.14    // 浮点型变量，赋值为3.14
	var c string = "Golang" // 字符串变量
	var d bool = true       // 布尔型变量

	// 打印变量的值
	fmt.Println("a =", a)
	fmt.Println("b =", b)
	fmt.Println("c =", c)
	fmt.Println("d =", d)
}
