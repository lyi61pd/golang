package main

import (
	"fmt"
	"reflect"
)

// 反射用于在运行时动态获取变量的类型和值，常用于通用库、序列化等场景
func main() {
	var x float64 = 3.14

	// 获取变量的类型
	t := reflect.TypeOf(x)
	fmt.Println("类型:", t.Name()) // 输出 float64

	// 获取变量的值
	v := reflect.ValueOf(x)
	fmt.Println("值:", v.Float()) // 输出 3.14

	// 通过反射修改变量的值（需要传递指针）
	var y float64 = 1.23
	p := reflect.ValueOf(&y)
	p.Elem().SetFloat(6.28)
	fmt.Println("y 修改后:", y)

	// 反射获取结构体字段信息
	type Person struct {
		Name string
		Age  int
	}
	p1 := Person{Name: "小明", Age: 18}
	val := reflect.ValueOf(p1)
	typ := reflect.TypeOf(p1)
	for i := 0; i < val.NumField(); i++ {
		fmt.Printf("字段名:%s, 字段值:%v\n", typ.Field(i).Name, val.Field(i).Interface())
	}
}
