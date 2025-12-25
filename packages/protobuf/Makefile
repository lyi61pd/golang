# Protobuf 代码生成和运行的 Makefile（增量构建版本）

# ========== 增量构建规则 ==========
# 声明：目标文件 依赖于 源文件
# 只有当 user.proto 比 user.pb.go 新时，才会重新生成

proto/user.pb.go: proto/user.proto
	@echo ">>> user.proto 有变化，重新生成 user.pb.go"
	protoc --go_out=. --go_opt=paths=source_relative proto/user.proto
	@echo ">>> 生成完成"

# gen 是别名，依赖真实的目标文件
gen: proto/user.pb.go
	@echo ">>> 代码生成检查完毕"

# ========== 普通任务（每次都执行） ==========

# 下载依赖
deps:
	go mod tidy

# 运行示例
run: proto/user.pb.go
	go run main.go

# 一键执行：生成代码 -> 下载依赖 -> 运行
all: gen deps run

# 安装 protoc-gen-go 插件
install-plugin:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 清理生成的文件（用于测试增量构建）
clean:
	rm -f proto/user.pb.go
	@echo ">>> 已删除 proto/user.pb.go"

# ========== 伪目标声明 ==========
# 注意：gen 也是 .PHONY，但它依赖的 proto/user.pb.go 是真实文件
# 所以增量构建逻辑在 proto/user.pb.go 规则里生效

.PHONY: gen deps run all install-plugin clean
