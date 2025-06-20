package main

import "fmt"

// 指针是存储变量地址的变量，可以用来间接修改数据
func main() {
	a := 10
	var p *int = &a // p 保存 a 的地址
	fmt.Println("a 的值:", a)
	fmt.Println("p 指向的值:", *p) // 通过指针访问 a 的值

	*p = 20 // 通过指针修改 a 的值
	fmt.Println("a 修改后:", a)
}
