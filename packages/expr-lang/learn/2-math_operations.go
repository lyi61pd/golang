// Package learn 提供expr包的学习示例
// 这个文件是学习系列的第二部分：expr中的数学运算
package learn

import (
	"fmt"
	"log"
	"math"

	"github.com/expr-lang/expr"
)

// MathOperations 展示expr包中的数学运算
// 这是学习expr的第二步，展示如何在表达式中使用数学函数
func MathOperations() {
	fmt.Println("=== 2. expr中的数学运算 ===")

	// 定义环境变量，包含数学运算所需的数据和函数
	// 这里我们将Go标准库中的数学函数注册到环境中
	env := map[string]interface{}{
		// 基本数值
		"price": 99.99,
		"tax":   0.08,
		"count": 3,

		// 注册数学函数，这样在表达式中就可以使用这些函数了
		"ceil":  math.Ceil,  // 向上取整
		"floor": math.Floor, // 向下取整
		"round": math.Round, // 四舍五入
		"max":   math.Max,   // 取最大值
		"min":   math.Min,   // 取最小值
		"pow":   math.Pow,   // 幂运算
		"sqrt":  math.Sqrt,  // 开平方根
	}

	// 定义包含数学运算的表达式
	// 这个表达式计算含税总价并向上取整
	code := `ceil(price * (1 + tax) * count)`

	// 编译表达式
	// 在编译时，expr会检查表达式中使用的函数和变量是否在环境中定义
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
	fmt.Printf("单价: %.2f, 税率: %.2f, 数量: %d\n", env["price"], env["tax"], env["count"])
	fmt.Printf("含税总价(向上取整): %.0f\n", output)
}