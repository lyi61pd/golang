package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/fx"
)

// 2. 进阶示例：带生命周期管理的 HTTP 服务器

// Config 配置结构
type Config struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Port: ":8080",
	}
}

// HTTPHandler 处理 HTTP 请求
type HTTPHandler struct {
	logger Logger
}

func NewHTTPHandler(logger Logger) *HTTPHandler {
	return &HTTPHandler{logger: logger}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Log(fmt.Sprintf("收到请求: %s %s", r.Method, r.URL.Path))

	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "欢迎使用 fx HTTP 服务器!\n")
	case "/health":
		fmt.Fprintf(w, "OK\n")
	case "/time":
		fmt.Fprintf(w, "当前时间: %s\n", time.Now().Format(time.RFC3339))
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found\n")
	}
}

// HTTPServer 封装 HTTP 服务器
type HTTPServer struct {
	server *http.Server
	logger Logger
}

func NewHTTPServer(lc fx.Lifecycle, config *Config, handler *HTTPHandler, logger Logger) *HTTPServer {
	srv := &http.Server{
		Addr:    config.Port,
		Handler: handler,
	}

	httpServer := &HTTPServer{
		server: srv,
		logger: logger,
	}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Log(fmt.Sprintf("启动 HTTP 服务器，监听端口 %s", config.Port))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Log(fmt.Sprintf("HTTP 服务器错误: %v", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Log("关闭 HTTP 服务器")
			return srv.Shutdown(ctx)
		},
	})

	return httpServer
}

// 要运行此示例：
// 1. 创建新目录：mkdir -p ../fxx-http && cd ../fxx-http
// 2. 复制文件：cp ../fxx/go.mod . && cp ../fxx/main.go logger.go && cp ../fxx/02_http_server.go main.go
// 3. 取消注释下面的 main 函数
// 4. go run .

/*
func main() {
	fmt.Println("=== fx HTTP 服务器示例 ===\n")

	app := fx.New(
		fx.Provide(
			NewSimpleLogger,
			NewConfig,
			NewHTTPHandler,
			NewHTTPServer,
		),
		fx.Invoke(func(server *HTTPServer) {
			// server 已经通过生命周期钩子启动
		}),
	)

	app.Run()
	// 服务器会一直运行，直到收到停止信号（Ctrl+C）
}
*/
