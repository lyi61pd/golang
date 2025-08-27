package main

import "fmt"

// add 是一个带返回值的函数，返回 a+b 的结果
func add(a int, b int) int {
	return a + b
}

// sayHello 是一个无返回值的函数，打印问候语
func sayHello(name string) {
	fmt.Println("Hello,", name)
}

func main() {
	sum := add(2, 3) // 调用 add 函数，传入 2 和 3
	fmt.Println("2 + 3 =", sum)
	sayHello("Gopher") // 调用 sayHello 函数，传入 "Gopher"
}
