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

// ========== 接口嵌入示例 ==========

// Mover 接口：会移动
type Mover interface {
	Move() string
}

// Eater 接口：会吃东西
type Eater interface {
	Eat() string
}

// LivingBeing 接口嵌入：组合了 Animal, Mover, Eater
// 任何实现 LivingBeing 的类型必须实现所有三个接口的方法
type LivingBeing interface {
	Animal         // 嵌入 Animal 接口
	Mover          // 嵌入 Mover 接口
	Eater          // 嵌入 Eater 接口
	Sleep() string // 自己额外定义的方法
}

// 完整的狗类型，实现了 LivingBeing 接口
type CompleteDog struct {
	Name string
}

func (d CompleteDog) Speak() string {
	return "汪汪"
}

func (d CompleteDog) Move() string {
	return "四条腿跑"
}

func (d CompleteDog) Eat() string {
	return "吃狗粮"
}

func (d CompleteDog) Sleep() string {
	return "趴着睡"
}

// 演示函数：接受 LivingBeing 接口
func describeLife(lb LivingBeing) {
	fmt.Printf("说话: %s\n", lb.Speak())
	fmt.Printf("移动: %s\n", lb.Move())
	fmt.Printf("吃: %s\n", lb.Eat())
	fmt.Printf("睡觉: %s\n", lb.Sleep())
}

func main() {
	fmt.Println("=== 基础接口示例 ===")
	var a Animal // 声明接口类型变量

	a = Dog{} // 可以赋值为实现了接口的类型
	fmt.Println("狗叫:", a.Speak())

	a = Cat{}
	fmt.Println("猫叫:", a.Speak())

	fmt.Println("\n=== 接口嵌入示例 ===")
	dog := CompleteDog{Name: "旺财"}
	describeLife(dog)

	// 接口类型断言
	var lb LivingBeing = dog
	// lb 同时也是 Animal, Mover, Eater 类型
	var animal Animal = lb
	var mover Mover = lb
	fmt.Printf("\n类型转换: %s, %s\n", animal.Speak(), mover.Move())
}
