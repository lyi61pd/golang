// Package learn 提供expr包的学习示例
// 这个文件是学习系列的第五部分：expr中对结构体和数组的操作
package learn

import (
	"fmt"
	"log"

	"github.com/expr-lang/expr"
)

// User 定义用户结构体
// 在expr中可以访问结构体的字段
type User struct {
	Name   string
	Age    int
	Active bool
	Score  float64
}

// Department 定义部门结构体
type Department struct {
	Name  string
	Users []User
}

// StructsAndArrays 展示expr包中对结构体和数组的操作
// 这是学习expr的第五步，展示如何在表达式中访问结构体字段和操作数组
func StructsAndArrays() {
	fmt.Println("=== 5. expr中对结构体和数组的操作 ===")

	// 创建测试数据
	users := []User{
		{"张三", 25, true, 95.5},
		{"李四", 17, false, 88.0},
		{"王五", 30, true, 92.3},
		{"赵六", 20, true, 78.9},
	}

	dept := Department{
		Name:  "技术部",
		Users: users,
	}

	// 定义环境变量
	env := map[string]interface{}{
		"user":       users[0],     // 单个用户
		"users":      users,        // 用户数组
		"department": dept,         // 部门对象
		"threshold":  18,           // 年龄阈值
		"minScore":   80.0,         // 最低分数
	}

	// 示例1: 访问结构体字段
	fmt.Println("--- 访问结构体字段 ---")
	fieldCode := `user.Name + "的年龄是" + string(user.Age) + "岁"`
	
	fieldProgram, err := expr.Compile(fieldCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译字段访问表达式失败:", err)
	}

	fieldOutput, err := expr.Run(fieldProgram, env)
	if err != nil {
		log.Fatal("执行字段访问表达式失败:", err)
	}

	fmt.Printf("字段访问结果: %s\n", fieldOutput)

	// 示例2: 数组索引访问
	fmt.Println("\n--- 数组索引访问 ---")
	indexCode := `users[1].Name + "的分数是" + string(users[1].Score)`
	
	indexProgram, err := expr.Compile(indexCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译索引访问表达式失败:", err)
	}

	indexOutput, err := expr.Run(indexProgram, env)
	if err != nil {
		log.Fatal("执行索引访问表达式失败:", err)
	}

	fmt.Printf("索引访问结果: %s\n", indexOutput)

	// 示例3: 数组过滤
	fmt.Println("\n--- 数组过滤 ---")
	// filter函数用于过滤数组，只保留满足条件的元素
	filterCode := `filter(users, {.Age >= threshold})`
	
	filterProgram, err := expr.Compile(filterCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译过滤表达式失败:", err)
	}

	filterOutput, err := expr.Run(filterProgram, env)
	if err != nil {
		log.Fatal("执行过滤表达式失败:", err)
	}

	fmt.Printf("年龄大于等于%d的用户: %v\n", env["threshold"], filterOutput)

	// 示例4: 数组映射
	fmt.Println("\n--- 数组映射 ---")
	// map函数用于将数组中的每个元素转换为新的值
	mapCode := `map(users, {.Name})`
	
	mapProgram, err := expr.Compile(mapCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译映射表达式失败:", err)
	}

	mapOutput, err := expr.Run(mapProgram, env)
	if err != nil {
		log.Fatal("执行映射表达式失败:", err)
	}

	fmt.Printf("所有用户名: %v\n", mapOutput)

	// 示例5: 数组归约
	fmt.Println("\n--- 数组归约 ---")
	// 由于expr内置不支持reduce函数，我们通过自定义函数实现
	env["sumScores"] = func(users []User) float64 {
		sum := 0.0
		for _, user := range users {
			sum += user.Score
		}
		return sum
	}
	
	reduceCode := `sumScores(users)`
	
	reduceProgram, err := expr.Compile(reduceCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译归约表达式失败:", err)
	}

	reduceOutput, err := expr.Run(reduceProgram, env)
	if err != nil {
		log.Fatal("执行归约表达式失败:", err)
	}

	fmt.Printf("所有用户总分: %.2f\n", reduceOutput)

	// 示例6: 复杂的嵌套结构访问
	fmt.Println("\n--- 复杂的嵌套结构访问 ---")
	nestedCode := `department.Name + "有" + string(len(department.Users)) + "名用户"`
	
	nestedProgram, err := expr.Compile(nestedCode, expr.Env(env))
	if err != nil {
		log.Fatal("编译嵌套结构表达式失败:", err)
	}

	nestedOutput, err := expr.Run(nestedProgram, env)
	if err != nil {
		log.Fatal("执行嵌套结构表达式失败:", err)
	}

	fmt.Printf("嵌套结构访问结果: %s\n", nestedOutput)
}