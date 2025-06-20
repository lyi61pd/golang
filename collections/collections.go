package main

import "fmt"

// 本例演示数组、切片和 map 的基本用法
func main() {
	// 数组：长度固定，元素类型相同
	var arr [3]int = [3]int{1, 2, 3}
	fmt.Println("数组:", arr)

	// 切片：长度可变，常用于动态数据
	s := []string{"Go", "Python", "Java"}
	s = append(s, "Rust") // 追加元素
	fmt.Println("切片:", s)

	// map：键值对集合，类似字典
	m := make(map[string]int) // 创建 map
	m["apple"] = 5
	m["banana"] = 3
	fmt.Println("map:", m)
	fmt.Println("apple数量:", m["apple"]) // 取 key 为 "apple" 的值
}
