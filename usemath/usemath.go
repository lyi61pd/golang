package main

import (
	"fmt"
	"golang-learning/mathutil" // 导入自定义包
)

func main() {
	sum := mathutil.Add(3, 5) // 调用包中的 Add 函数
	fmt.Println("3 + 5 =", sum)
}
