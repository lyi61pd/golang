package main

import "fmt"

// Person 结构体用于描述一个人的基本信息
// 包含姓名和年龄两个字段
type Person struct {
	Name string // 姓名
	Age  int    // 年龄
}

// Greet 方法是 Person 结构体关联的函数
// 用于打印问候语，包括姓名和年龄
func (p Person) Greet() {
	fmt.Printf("你好，我是 %s，%d 岁。\n", p.Name, p.Age)
}

func main() {
	// 创建 Person 结构体实例并初始化字段
	p := Person{Name: "小明", Age: 18}
	// 调用结构体方法 Greet
	p.Greet()
}
