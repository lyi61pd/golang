// Package learn 提供expr包的学习示例
// 这个文件是学习系列的第一部分：expr基本用法
package learn

import (
	"fmt"
	"log"

	"github.com/expr-lang/expr"
)

// BasicUsage 展示expr包的基本用法
// 这是学习expr的第一步，展示如何编译和执行一个简单的表达式n
func BasicUsage() {
	fmt.Println("=== 1. expr基本用法 ===")

	// 定义环境变量，这些变量可以在表达式中使用
	// 环境变量是表达式可以访问的数据和函数集合
	env := map[string]interface{}{
		"name": "张三",
		"age":  25,
	}

	// 定义要执行的表达式
	// 这个表达式会使用环境中的变量
	code := `"你好，" + name + "，你的年龄是" + string(age) + "岁"`

	// 编译表达式
	// 编译阶段会检查表达式的语法和类型安全性
	// expr.Env(env) 参数指定了表达式可以访问的环境变量
	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		log.Fatal("编译表达式失败:", err)
	}

	// 执行编译后的表达式
	// expr.Run函数执行表达式并返回结果
	output, err := expr.Run(program, env)
	if err != nil {
		log.Fatal("执行表达式失败:", err)
	}

	// 输出结果
	fmt.Println("表达式执行结果:", output)
}