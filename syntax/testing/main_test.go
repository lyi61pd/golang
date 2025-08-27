package testing

import "testing"

// Add 是被测试的函数，返回两数之和
func Add(a, b int) int {
	return a + b
}

// TestAdd 是基础单元测试
func TestAdd(t *testing.T) {
	result := Add(1, 2)
	if result != 3 {
		t.Errorf("期望 3，得到 %d", result)
	}
}

// TestAddTableDriven 是表驱动测试的示例
func TestAddTableDriven(t *testing.T) {
	// 定义一组测试用例，每个用例包含输入和期望输出
	cases := []struct {
		a, b   int // 输入参数
		expect int // 期望的结果
	}{
		{1, 2, 3},       // 1+2=3
		{0, 0, 0},       // 0+0=0
		{-1, 1, 0},      // -1+1=0
		{100, 200, 300}, // 100+200=300
	}

	// 遍历每个测试用例
	for _, c := range cases {
		got := Add(c.a, c.b) // 调用被测试函数
		// 检查实际结果是否等于期望结果
		if got != c.expect {
			// 如果不等，报告错误，显示输入和实际/期望输出
			t.Errorf("Add(%d, %d) 期望 %d，得到 %d", c.a, c.b, c.expect, got)
		}
	}
}

// 运行测试：在 testing 目录下执行 go test
// go test 会自动查找所有 _test.go 文件并运行以 Test 开头的函数
// go test 会自动查找所有 _test.go 文件并运行以 Test 开头的函数
