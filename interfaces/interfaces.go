package main

import "fmt"

// 接口定义了一组方法，任何实现了这些方法的类型都实现了该接口
type Animal interface {
	Speak() string
}

// Dog 类型，实现了 Animal 接口
type Dog struct{}

// Dog 实现 Speak 方法
func (d Dog) Speak() string {
	return "汪汪"
}

// Cat 类型，实现了 Animal 接口
type Cat struct{}

// Cat 实现 Speak 方法
func (c Cat) Speak() string {
	return "喵喵"
}

func main() {
	var a Animal // 声明接口类型变量

	a = Dog{} // 可以赋值为实现了接口的类型
	fmt.Println("狗叫:", a.Speak())

	a = Cat{}
	fmt.Println("猫叫:", a.Speak())
}
