package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User 结构体用于表示用户信息
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var initialUsers = []User{
	{ID: 1, Name: "Alice"},
	{ID: 2, Name: "Bob"},
}

var users = append([]User{}, initialUsers...)

// Logger 是一个简单的中间件示例
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("请求路径:", c.Request.URL.Path)
		c.Next()
	}
}

// LoginForm 用于参数绑定示例
type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// gin-jwt 认证中间件配置
var identityKey = "id"

// GORM模型定义（可与 User 结构体一致或更丰富）
type GormUser struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

func main() {
	// 设置 Gin 运行模式，可选 gin.DebugMode/gin.ReleaseMode/gin.TestMode
	gin.SetMode(gin.ReleaseMode)

	// r := gin.Default() // 原有代码
	r := gin.New() // 使用 gin.New() 不自动注册 Logger/Recovery

	// 注册全局中间件
	r.Use(Logger())
	r.Use(gin.Recovery())     // 推荐加上 Recovery 中间件，防止 panic 导致服务崩溃
	r.Use(TimingMiddleware()) // 请求耗时统计中间件
	// r.Use(AuthMiddleware()) // 简单鉴权中间件

	// 静态文件服务，将 ./static 目录映射到 /static 路径
	// 访问方式: http://localhost:8080/static/文件名
	r.Static("/static", "./static")

	// 路由分组示例
	api := r.Group("/api")
	{
		// 示例: GET http://localhost:8080/api/ping
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "api pong"})
		})
	}

	// /ping 路由，返回 pong 消息
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// /hello 路由，接收 URL 查询参数，返回问候消息
	// 调用方式: GET http://localhost:8080/hello?name=YourName
	r.GET("/hello", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			name = "world"
		}
		c.JSON(200, gin.H{
			"message": "Hello " + name,
		})
	})

	// /login 路由，处理 POST 请求，接收表单参数
	// 调用方式: POST http://localhost:8080/login
	// 表单参数: username, password
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		c.JSON(200, gin.H{
			"username": username,
			"password": password,
		})
	})

	// /user/:name 路由，获取路径参数，返回用户名称
	// 调用方式: GET http://localhost:8080/user/YourName
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(200, gin.H{
			"user": name,
		})
	})

	// 获取用户列表
	// 调用方式: curl http://localhost:8080/users
	r.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, users)
	})

	// 添加用户
	// 调用方式: curl -X POST -H "Content-Type: application/json" -d '{"id":3,"name":"Charlie"}' http://localhost:8080/users
	r.POST("/users", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		users = append(users, newUser)
		c.JSON(http.StatusOK, newUser)
	})

	// 根据用户ID获取用户详情
	// 调用方式: curl http://localhost:8080/users/1
	r.GET("/users/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}
		for _, user := range users {
			if user.ID == id {
				c.JSON(http.StatusOK, user)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	})

	// 删除用户
	// 调用方式: curl -X DELETE http://localhost:8080/users/1
	r.DELETE("/users/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}
		for i, user := range users {
			if user.ID == id {
				// 从切片中删除该用户
				users = append(users[:i], users[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	})

	// 更新用户信息
	// 调用方式: curl -X PUT -H "Content-Type: application/json" -d '{"name":"NewName"}' http://localhost:8080/users/1
	r.PUT("/users/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}
		var updateData struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for i, user := range users {
			if user.ID == id {
				users[i].Name = updateData.Name
				c.JSON(http.StatusOK, users[i])
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	})

	// 按用户名模糊查询用户
	// 调用方式: curl "http://localhost:8080/search?name=al"
	r.GET("/search", func(c *gin.Context) {
		query := c.Query("name")
		var result []User
		for _, user := range users {
			if strings.Contains(strings.ToLower(user.Name), strings.ToLower(query)) {
				result = append(result, user)
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// 统计用户数量
	// 调用方式: curl http://localhost:8080/users/count
	r.GET("/users/count", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"count": len(users),
		})
	})

	// 重置用户列表为初始状态
	// 调用方式: curl -X POST http://localhost:8080/users/reset
	r.POST("/users/reset", func(c *gin.Context) {
		users = append([]User{}, initialUsers...)
		c.JSON(http.StatusOK, gin.H{
			"message": "User list reset",
			"users":   users,
		})
	})

	// 参数绑定示例接口
	// 支持 application/json 或 application/x-www-form-urlencoded
	// 调用方式:
	// curl -X POST -H "Content-Type: application/json" -d '{"username":"admin","password":"123456"}' http://localhost:8080/bind
	// 或
	// curl -X POST -d "username=admin&password=123456" http://localhost:8080/bind
	r.POST("/bind", func(c *gin.Context) {
		var form LoginForm
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"username": form.Username})
	})

	// 请求上下文示例
	// 调用方式: curl -H "Authorization: token123" http://localhost:8080/context
	r.GET("/context", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		c.Header("X-App", "gin-demo")
		c.JSON(200, gin.H{"token": token})
	})

	// 重定向示例
	// 调用方式: curl -i http://localhost:8080/redirect
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(302, "https://www.baidu.com")
	})

	// 单文件上传示例
	// 调用方式: curl -F "file=@/path/to/your/file.txt" http://localhost:8080/upload
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// 保存文件到当前目录
		c.SaveUploadedFile(file, "./"+file.Filename)
		c.JSON(200, gin.H{"filename": file.Filename})
	})

	// 返回 XML 示例
	// 调用方式: curl http://localhost:8080/xml
	r.GET("/xml", func(c *gin.Context) {
		c.XML(200, gin.H{"message": "hello", "status": "ok"})
	})

	// 返回纯文本示例
	// 调用方式: curl http://localhost:8080/text
	r.GET("/text", func(c *gin.Context) {
		c.String(200, "hello, world")
	})

	// 参数获取演示接口
	// 调用方式:
	// curl "http://localhost:8080/params-demo/123?query=abc" -H "X-Token: mytoken" -d "formval=formdata" -X POST
	r.POST("/params-demo/:id", func(c *gin.Context) {
		// 获取路径参数
		id := c.Param("id")
		// 获取查询参数
		query := c.Query("query")
		// 获取表单参数
		formval := c.PostForm("formval")
		// 获取请求头
		token := c.GetHeader("X-Token")

		c.JSON(200, gin.H{
			"id":      id,
			"query":   query,
			"formval": formval,
			"token":   token,
		})
	})

	// 获取所有请求头、所有查询参数、所有表单参数的演示接口
	// 调用方式:
	// curl -X POST "http://localhost:8080/all-params?foo=bar&baz=qux" -H "X-Test: testval" -d "a=1&b=2"
	r.POST("/all-params", func(c *gin.Context) {
		// 获取所有请求头
		headers := map[string]string{}
		for k, v := range c.Request.Header {
			headers[k] = strings.Join(v, ",")
		}
		// 获取所有查询参数
		querys := map[string]string{}
		for k, v := range c.Request.URL.Query() {
			querys[k] = strings.Join(v, ",")
		}
		// 获取所有表单参数
		c.Request.ParseForm()
		forms := map[string]string{}
		for k, v := range c.Request.PostForm {
			forms[k] = strings.Join(v, ",")
		}
		c.JSON(200, gin.H{
			"headers": headers,
			"query":   querys,
			"form":    forms,
		})
	})

	// 演示如何获取客户端IP和请求方法
	// 调用方式: curl -X GET http://localhost:8080/request-info
	r.GET("/request-info", func(c *gin.Context) {
		clientIP := c.ClientIP()
		method := c.Request.Method
		c.JSON(200, gin.H{
			"client_ip": clientIP,
			"method":    method,
		})
	})

	// 演示如何设置和获取 Cookie
	// 设置 Cookie: curl -X GET "http://localhost:8080/set-cookie"
	r.GET("/set-cookie", func(c *gin.Context) {
		// 设置名为 "mycookie" 的 Cookie，值为 "hello", 有效期 1 小时
		c.SetCookie("mycookie", "hello", 3600, "/", "", false, true)
		c.JSON(200, gin.H{"message": "cookie set"})
	})

	// 获取 Cookie: curl -X GET --cookie "mycookie=hello" "http://localhost:8080/get-cookie"
	r.GET("/get-cookie", func(c *gin.Context) {
		val, err := c.Cookie("mycookie")
		if err != nil {
			c.JSON(400, gin.H{"error": "cookie not found"})
			return
		}
		c.JSON(200, gin.H{"mycookie": val})
	})

	// 演示如何获取和设置自定义 context 变量（在中间件和 handler 之间传递数据）
	// 调用方式: curl http://localhost:8080/context-demo
	r.GET("/context-demo", func(c *gin.Context) {
		// 在 context 中设置自定义变量
		c.Set("user_id", 1001)
		c.Next()
	}, func(c *gin.Context) {
		// 在下一个 handler 中获取变量
		val, exists := c.Get("user_id")
		if !exists {
			c.JSON(400, gin.H{"error": "user_id not found"})
			return
		}
		c.JSON(200, gin.H{"user_id": val})
	})

	// 演示如何获取原始请求体（raw body）
	// 调用方式: curl -X POST http://localhost:8080/raw-body -d 'rawdata=abc'
	r.POST("/raw-body", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"raw_body": string(body)})
	})

	// 演示如何返回 YAML 格式响应
	// 调用方式: curl http://localhost:8080/yaml
	r.GET("/yaml", func(c *gin.Context) {
		c.YAML(200, gin.H{
			"message": "hello",
			"status":  "ok",
		})
	})

	// 统一处理未匹配的路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "接口不存在"})
	})

	// 通过环境变量 PORT 设置端口，默认 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// gin-jwt 中间件实例
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "example zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"name":      v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				ID:   int(claims[identityKey].(float64)),
				Name: claims["name"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals LoginForm
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			// 简单演示，实际应查数据库
			if (username == "admin" && password == "123456") || (username == "alice" && password == "123456") {
				return &User{
					ID:   1,
					Name: username,
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// 可自定义权限
			if _, ok := data.(*User); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// 登录接口（自动生成token）
	// curl -X POST -d "username=admin&password=123456" http://localhost:8080/login-jwt
	r.POST("/login-jwt", authMiddleware.LoginHandler)

	// 刷新token接口
	r.GET("/refresh-token", authMiddleware.RefreshHandler)

	// 受保护的路由分组
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// curl -H "Authorization: Bearer <token>" http://localhost:8080/auth/profile
		auth.GET("/profile", func(c *gin.Context) {
			user, _ := c.Get(identityKey)
			c.JSON(200, gin.H{
				"user":   user,
				"claims": jwt.ExtractClaims(c),
			})
		})
	}

	// 路由分组示例：以 /api/v1 为前缀，分组管理 RESTful 资源
	v1 := r.Group("/api/v1")
	{
		// 用户资源 RESTful 路由
		// GET /api/v1/users      - 获取用户列表
		// POST /api/v1/users     - 创建用户
		// GET /api/v1/users/:id  - 获取单个用户
		// PUT /api/v1/users/:id  - 更新用户
		// DELETE /api/v1/users/:id - 删除用户

		v1.GET("/users", func(c *gin.Context) {
			// ...可复用上面 /users 逻辑...
			c.JSON(http.StatusOK, users)
		})

		v1.POST("/users", func(c *gin.Context) {
			var newUser User
			if err := c.ShouldBindJSON(&newUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			users = append(users, newUser)
			c.JSON(http.StatusOK, newUser)
		})

		v1.GET("/users/:id", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
				return
			}
			for _, user := range users {
				if user.ID == id {
					c.JSON(http.StatusOK, user)
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		})

		v1.PUT("/users/:id", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
				return
			}
			var updateData struct {
				Name string `json:"name"`
			}
			if err := c.ShouldBindJSON(&updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			for i, user := range users {
				if user.ID == id {
					users[i].Name = updateData.Name
					c.JSON(http.StatusOK, users[i])
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		})

		v1.DELETE("/users/:id", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
				return
			}
			for i, user := range users {
				if user.ID == id {
					users = append(users[:i], users[i+1:]...)
					c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		})
	}

	// 初始化 GORM（以 SQLite 为例，实际可用 MySQL/Postgres）
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	db.AutoMigrate(&GormUser{})

	// GORM 高级API分组
	gormApi := r.Group("/gorm")
	{
		// 创建用户
		// curl -X POST -H "Content-Type: application/json" -d '{"name":"Tom"}' http://localhost:8080/gorm/users
		gormApi.POST("/users", func(c *gin.Context) {
			var user GormUser
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if err := db.Create(&user).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, user)
		})

		// 查询所有用户
		// curl http://localhost:8080/gorm/users
		gormApi.GET("/users", func(c *gin.Context) {
			var users []GormUser
			if err := db.Find(&users).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, users)
		})

		// 查询单个用户
		// curl http://localhost:8080/gorm/users/1
		gormApi.GET("/users/:id", func(c *gin.Context) {
			var user GormUser
			if err := db.First(&user, c.Param("id")).Error; err != nil {
				c.JSON(404, gin.H{"error": "not found"})
				return
			}
			c.JSON(200, user)
		})

		// 更新用户
		// curl -X PUT -H "Content-Type: application/json" -d '{"name":"Jerry"}' http://localhost:8080/gorm/users/1
		gormApi.PUT("/users/:id", func(c *gin.Context) {
			var user GormUser
			if err := db.First(&user, c.Param("id")).Error; err != nil {
				c.JSON(404, gin.H{"error": "not found"})
				return
			}
			var update struct{ Name string }
			if err := c.ShouldBindJSON(&update); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			user.Name = update.Name
			db.Save(&user)
			c.JSON(200, user)
		})

		// 删除用户
		// curl -X DELETE http://localhost:8080/gorm/users/1
		gormApi.DELETE("/users/:id", func(c *gin.Context) {
			if err := db.Delete(&GormUser{}, c.Param("id")).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"message": "deleted"})
		})

		// 条件查询与分页
		// curl "http://localhost:8080/gorm/query?name=Tom&page=1&page_size=2"
		gormApi.GET("/query", func(c *gin.Context) {
			var users []GormUser
			name := c.Query("name")
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
			if page < 1 {
				page = 1
			}
			if pageSize < 1 {
				pageSize = 10
			}
			query := db.Model(&GormUser{})
			if name != "" {
				query = query.Where("name LIKE ?", "%"+name+"%")
			}
			var total int64
			query.Count(&total)
			query = query.Offset((page - 1) * pageSize).Limit(pageSize)
			if err := query.Find(&users).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{
				"total":     total,
				"page":      page,
				"page_size": pageSize,
				"data":      users,
			})
		})

		// 排序
		// curl "http://localhost:8080/gorm/sorted?order=desc"
		gormApi.GET("/sorted", func(c *gin.Context) {
			var users []GormUser
			order := c.DefaultQuery("order", "asc")
			if order != "asc" && order != "desc" {
				order = "asc"
			}
			if err := db.Order("id " + order).Find(&users).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, users)
		})

		// 事务示例
		// curl -X POST -H "Content-Type: application/json" -d '{"name":"TxUser"}' http://localhost:8080/gorm/tx
		gormApi.POST("/tx", func(c *gin.Context) {
			var user GormUser
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(&user).Error; err != nil {
					return err
				}
				// 可以在这里做更多操作，出错会自动回滚
				return nil
			})
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, user)
		})

		// 批量插入
		// curl -X POST -H "Content-Type: application/json" -d '[{"name":"A"},{"name":"B"}]' http://localhost:8080/gorm/batch
		gormApi.POST("/batch", func(c *gin.Context) {
			var users []GormUser
			if err := c.ShouldBindJSON(&users); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if err := db.Create(&users).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, users)
		})
	}

	// 启动 HTTP 服务（非阻塞）
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	fmt.Println("服务已启动，监听端口:", port)

	// 等待中断信号以优雅关停
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("收到退出信号，正在优雅关停...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	fmt.Println("服务已安全关闭")
}

// TimingMiddleware 统计请求耗时的中间件
func TimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		cost := time.Since(start)
		fmt.Printf("请求 %s 耗时: %v\n", c.Request.URL.Path, cost)
	}
}

// AuthMiddleware 简单的鉴权中间件，要求请求头 X-Auth=secret
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Auth") != "secret" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
