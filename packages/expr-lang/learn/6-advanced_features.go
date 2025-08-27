// Package learn 提供expr包的学习示例
// 这个文件是学习系列的第六部分：expr的高级特性
package learn

import (
	"fmt"
	"log"
	"regexp"

	"github.com/expr-lang/expr"
)

// AdvancedFeatures 展示expr包的高级特性
// 这是学习expr的第六步，展示一些更高级的功能和用法
func AdvancedFeatures() {
	fmt.Println("=== 6. expr的高级特性 ===")

	// 示例1: 正则表达式匹配
	fmt.Println("--- 正则表达式匹配 ---")
	env1 := map[string]interface{}{
		"text": "hello@example.com",
		"match": func(pattern, text string) bool {
			matched, _ := regexpMatch(pattern, text)
			return matched
		},
	}

	// 检查邮箱格式
	// 注意：在Go字符串中需要转义反斜杠
	code1 := `match("[a-zA-Z0-9._%+\\-]+@[a-zA-Z0-9.\\-]+\\.[a-zA-Z]{2,}", text)`
	
	program1, err := expr.Compile(code1, expr.Env(env1))
	if err != nil {
		log.Fatal("编译正则表达式失败:", err)
	}

	output1, err := expr.Run(program1, env1)
	if err != nil {
		log.Fatal("执行正则表达式失败:", err)
	}

	fmt.Printf("邮箱格式检查结果: %v\n", output1)

	// 示例2: 错误处理和验证
	fmt.Println("\n--- 表达式验证和错误处理 ---")
	env2 := map[string]interface{}{
		"a": 10,
		"b": 5,
	}

	// 正确的表达式
	correctCode := `a + b * 2`
	
	// 验证表达式（不执行）
	_, err = expr.Compile(correctCode, expr.Env(env2))
	if err != nil {
		fmt.Printf("正确表达式验证失败: %v\n", err)
	} else {
		fmt.Printf("表达式 '%s' 语法正确\n", correctCode)
	}

	// 错误的表达式示例
	wrongCode := `a + b *` // 不完整的表达式
	_, err = expr.Compile(wrongCode, expr.Env(env2))
	if err != nil {
		fmt.Printf("错误表达式 '%s' 验证失败: %v\n", wrongCode, err)
	}

	// 示例3: 动态表达式执行
	fmt.Println("\n--- 动态表达式执行 ---")
	env3 := map[string]interface{}{
		"x": 10,
		"y": 3,
		"operations": map[string]string{
			"add":      "x + y",
			"subtract": "x - y",
			"multiply": "x * y",
			"divide":   "x / y",
		},
	}

	// 动态选择并执行表达式
	for opName, opExpr := range env3["operations"].(map[string]string) {
		program, err := expr.Compile(opExpr, expr.Env(env3))
		if err != nil {
			log.Printf("编译%s表达式失败: %v", opName, err)
			continue
		}

		result, err := expr.Run(program, env3)
		if err != nil {
			log.Printf("执行%s表达式失败: %v", opName, err)
			continue
		}

		fmt.Printf("%s: %s = %v\n", opName, opExpr, result)
	}

	// 示例4: 使用选项自定义编译器行为
	fmt.Println("\n--- 使用编译器选项 ---")
	env4 := map[string]interface{}{
		"name": "expr学习者",
		"age":  25,
	}

	// 限制表达式中允许的函数
	// 只允许使用len函数
	code4 := `len(name) > 5`
	
	program4, err := expr.Compile(code4, 
		expr.Env(env4))
	if err != nil {
		log.Fatal("编译受限表达式失败:", err)
	}

	output4, err := expr.Run(program4, env4)
	if err != nil {
		log.Fatal("执行受限表达式失败:", err)
	}

	fmt.Printf("表达式 '%s' 的结果: %v\n", code4, output4)

	// 示例5: 表达式缓存（编译一次，多次执行）
	fmt.Println("\n--- 表达式缓存 ---")
	env5 := map[string]interface{}{
		"price":    100,
		"discount": 0.1,
		"tax":      0.08,
	}

	// 编译一次表达式
	code5 := `price * (1 - discount) * (1 + tax)`
	program5, err := expr.Compile(code5, expr.Env(env5))
	if err != nil {
		log.Fatal("编译缓存表达式失败:", err)
	}

	// 多次执行同一个编译后的表达式（提高性能）
	testCases := []map[string]interface{}{
		{"price": 100, "discount": 0.1, "tax": 0.08},
		{"price": 200, "discount": 0.15, "tax": 0.08},
		{"price": 50, "discount": 0.05, "tax": 0.08},
	}

	for i, testCase := range testCases {
		result, err := expr.Run(program5, testCase)
		if err != nil {
			log.Printf("执行第%d个测试用例失败: %v", i+1, err)
			continue
		}
		fmt.Printf("测试用例 %d: %+v => 最终价格: %.2f\n", 
			i+1, testCase, result)
	}
}

// regexpMatch 简单的正则表达式匹配函数
// 用于演示如何在expr中使用正则表达式
func regexpMatch(pattern, text string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(text), nil
}