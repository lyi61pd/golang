package mathutil // 包名通常与文件夹名一致

// Add 是导出函数（首字母大写），可被其他包调用
func Add(a, b int) int {
	return a + b
}
