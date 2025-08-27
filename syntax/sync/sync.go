package main

import (
	"fmt"
	"sync"
)

// WaitGroup 用于等待一组 goroutine 完成
func main() {
	var wg sync.WaitGroup

	// 启动3个并发任务
	for i := 1; i <= 3; i++ {
		wg.Add(1) // 计数器+1
		go func(id int) {
			defer wg.Done() // 任务完成，计数器-1
			fmt.Printf("任务 %d 完成\n", id)
		}(i)
	}

	wg.Wait() // 等待所有任务完成
	fmt.Println("所有任务已完成")
}
