package main

import "fmt"

func printType(i interface{}) {
	// 类型断言：判断 i 是否为 int 类型
	v, ok := i.(int)
	if ok {
		fmt.Println("是 int 类型，值为", v)
	} else {
		fmt.Println("不是 int 类型")
	}
}

func main() {
	printType(123)
	printType("hello")
}
