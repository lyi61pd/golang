```
cd /Users/ybbj100324/code/self/golang_demo/packages/protobuf

# 1. 安装 protoc-gen-go 插件（只需执行一次）
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 2. 生成 Go 代码
protoc --go_out=. --go_opt=paths=source_relative proto/user.proto

# 3. 下载依赖
go mod tidy

# 4. 运行示例
go run main.go
```