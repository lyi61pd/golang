package main

import (
	"errors"
	"fmt"
)

// Go 推荐通过返回 error 类型值进行错误处理
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除数不能为0")
	}
	return a / b, nil
}

func main() {
	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("发生错误:", err)
	} else {
		fmt.Println("结果:", result)
	}
}
