package main

import "fmt"

// 本例演示 if、for、switch 三种流程控制语句
func main() {
	// if 语句：条件判断
	x := 10
	if x > 5 {
		fmt.Println("x 大于 5")
	} else {
		fmt.Println("x 小于等于 5")
	}

	// for 循环：计算 1 到 5 的和
	sum := 0
	for i := 1; i <= 5; i++ {
		sum += i
	}
	fmt.Println("sum =", sum)

	// switch 语句：多分支选择
	day := 3
	switch day {
	case 1:
		fmt.Println("星期一")
	case 2:
		fmt.Println("星期二")
	case 3:
		fmt.Println("星期三")
	default:
		fmt.Println("其他")
	}
}
