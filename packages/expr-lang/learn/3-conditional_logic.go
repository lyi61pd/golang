// Package learn 提供expr包的学习示例
// 这个文件是学习系列的第三部分：expr中的条件逻辑
package learn

import (
	"fmt"
	"log"

	"github.com/expr-lang/expr"
)

// ConditionalLogic 展示expr包中的条件逻辑
// 这是学习expr的第三步，展示如何在表达式中使用条件判断
func ConditionalLogic() {
	fmt.Println("=== 3. expr中的条件逻辑 ===")

	// 定义环境变量，包含条件判断所需的数据
	env := map[string]interface{}{
		"user": map[string]interface{}{
			"name":   "李四",
			"age":    17,
			"active": true,
			"roles":  []string{"user", "editor"},
		},
		"vipThreshold": 18,
	}

	// 定义包含条件逻辑的表达式
	// 使用三元运算符 ?: 进行条件判断
	// 如果用户年龄大于等于vipThreshold且账户活跃，则返回VIP消息，否则返回普通用户消息
	code := `user.age >= vipThreshold && user.active ? 
             "欢迎VIP用户 " + user.name + "!" : 
             "欢迎普通用户 " + user.name + "，请继续努力!"`

	// 编译表达式
	// expr会检查条件表达式的语法和类型
	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		log.Fatal("编译表达式失败:", err)
	}

	// 执行表达式
	output, err := expr.Run(program, env)
	if err != nil {
		log.Fatal("执行表达式失败:", err)
	}

	// 输出结果
	fmt.Printf("用户信息: %+v\n", env["user"])
	fmt.Printf("判断结果: %s\n", output)

	// 更复杂的条件示例 - 检查用户角色
	fmt.Println("\n--- 检查用户角色示例 ---")
	roleCheckCode := `"user" in user.roles ? "用户拥有基础权限" : "用户没有基础权限"`
	
	roleProgram, err := expr.Compile(roleCheckCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译角色检查表达式失败:", err)
	}

	roleOutput, err := expr.Run(roleProgram, env)
	if err != nil {
		log.Fatal("执行角色检查表达式失败:", err)
	}

	fmt.Printf("角色检查结果: %s\n", roleOutput)
}