// Package learn 提供expr包的学习示例
// 这个文件是学习系列的第四部分：expr中的自定义函数
package learn

import (
	"fmt"
	"log"
	"strings"

	"github.com/expr-lang/expr"
)

// CustomFunctions 展示expr包中的自定义函数使用
// 这是学习expr的第四步，展示如何在表达式中注册和使用自定义函数
func CustomFunctions() {
	fmt.Println("=== 4. expr中的自定义函数 ===")

	// 定义自定义函数
	// 这些函数可以在表达式中使用
	env := map[string]interface{}{
		// 自定义的字符串处理函数
		"toUpperCase": strings.ToUpper,
		"toLowerCase": strings.ToLower,
		"join":        strings.Join,
		
		// 自定义的数学函数
		"square": func(x int) int {
			return x * x
		},
		
		// 自定义的条件函数
		"isAdult": func(age int) bool {
			return age >= 18
		},
		
		// 自定义的格式化函数
		"formatUserInfo": func(name string, age int) string {
			return fmt.Sprintf("用户: %s, 年龄: %d", name, age)
		},
		
		// 数据变量
		"name": "WangWu",
		"age":  20,
		"tags": []string{"Go", "expr", "tutorial"},
	}

	// 示例1: 使用字符串处理函数
	fmt.Println("--- 使用字符串处理函数 ---")
	strCode := `toUpperCase(name) + " " + toLowerCase(name)`
	
	strProgram, err := expr.Compile(strCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译字符串处理表达式失败:", err)
	}

	strOutput, err := expr.Run(strProgram, env)
	if err != nil {
		log.Fatal("执行字符串处理表达式失败:", err)
	}

	fmt.Printf("字符串处理结果: %s\n", strOutput)

	// 示例2: 使用自定义数学函数
	fmt.Println("\n--- 使用自定义数学函数 ---")
	mathCode := `square(age) + 10`
	
	mathProgram, err := expr.Compile(mathCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译数学表达式失败:", err)
	}

	mathOutput, err := expr.Run(mathProgram, env)
	if err != nil {
		log.Fatal("执行数学表达式失败:", err)
	}

	fmt.Printf("%d的平方加10等于: %v\n", env["age"], mathOutput)

	// 示例3: 使用条件函数
	fmt.Println("\n--- 使用条件函数 ---")
	condCode := `isAdult(age) ? "成年人" : "未成年人"`
	
	condProgram, err := expr.Compile(condCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译条件表达式失败:", err)
	}

	condOutput, err := expr.Run(condProgram, env)
	if err != nil {
		log.Fatal("执行条件表达式失败:", err)
	}

	fmt.Printf("年龄%d是: %s\n", env["age"], condOutput)

	// 示例4: 使用格式化函数
	fmt.Println("\n--- 使用格式化函数 ---")
	formatCode := `formatUserInfo(name, age)`
	
	formatProgram, err := expr.Compile(formatCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译格式化表达式失败:", err)
	}

	formatOutput, err := expr.Run(formatProgram, env)
	if err != nil {
		log.Fatal("执行格式化表达式失败:", err)
	}

	fmt.Printf("格式化结果: %s\n", formatOutput)

	// 示例5: 使用数组和字符串函数组合
	fmt.Println("\n--- 使用数组和字符串函数组合 ---")
	arrayCode := `join(tags, ", ")`
	
	arrayProgram, err := expr.Compile(arrayCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译数组处理表达式失败:", err)
	}

	arrayOutput, err := expr.Run(arrayProgram, env)
	if err != nil {
		log.Fatal("执行数组处理表达式失败:", err)
	}

	fmt.Printf("标签列表: %s\n", arrayOutput)
}