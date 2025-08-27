package main

import "fmt"

// 泛型函数，T 可以是任何类型
func PrintSlice[T any](s []T) {
	for _, v := range s {
		fmt.Println(v)
	}
}

// 泛型类型示例
type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

func main() {
	// 泛型函数示例
	ints := []int{1, 2, 3}
	strs := []string{"a", "b", "c"}
	PrintSlice(ints)
	PrintSlice(strs)

	// 泛型结构体示例
	p := Pair[int, string]{First: 1, Second: "one"}
	fmt.Println("Pair:", p.First, p.Second)
}
