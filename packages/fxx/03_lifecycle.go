package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

// 3. 生命周期管理示例

// Database 模拟数据库连接
type Database struct {
	connected bool
	logger    Logger
}

func NewDatabase(lc fx.Lifecycle, logger Logger) *Database {
	db := &Database{logger: logger}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Log("正在连接数据库...")
			time.Sleep(500 * time.Millisecond) // 模拟连接延迟
			db.connected = true
			logger.Log("✓ 数据库连接成功")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Log("正在断开数据库连接...")
			time.Sleep(200 * time.Millisecond) // 模拟断开延迟
			db.connected = false
			logger.Log("✓ 数据库连接已关闭")
			return nil
		},
	})

	return db
}

func (db *Database) Query(sql string) string {
	if !db.connected {
		return "错误：数据库未连接"
	}
	db.logger.Log(fmt.Sprintf("执行查询: %s", sql))
	return "查询结果: [数据...]"
}

// Cache 模拟缓存服务
type Cache struct {
	started bool
	logger  Logger
}

func NewCache(lc fx.Lifecycle, logger Logger) *Cache {
	cache := &Cache{logger: logger}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Log("正在启动缓存服务...")
			time.Sleep(300 * time.Millisecond)
			cache.started = true
			logger.Log("✓ 缓存服务启动成功")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Log("正在停止缓存服务...")
			time.Sleep(100 * time.Millisecond)
			cache.started = false
			logger.Log("✓ 缓存服务已停止")
			return nil
		},
	})

	return cache
}

func (c *Cache) Get(key string) string {
	if !c.started {
		return "错误：缓存服务未启动"
	}
	c.logger.Log(fmt.Sprintf("从缓存获取: %s", key))
	return fmt.Sprintf("缓存值: %s", key)
}

// UserService 依赖数据库和缓存
type UserService struct {
	db     *Database
	cache  *Cache
	logger Logger
}

func NewUserService(db *Database, cache *Cache, logger Logger) *UserService {
	logger.Log("✓ UserService 已创建")
	return &UserService{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

func (s *UserService) GetUser(id string) {
	s.logger.Log(fmt.Sprintf("\n--- 获取用户 %s ---", id))

	// 先查缓存
	cacheResult := s.cache.Get(id)
	s.logger.Log(cacheResult)

	// 再查数据库
	dbResult := s.db.Query(fmt.Sprintf("SELECT * FROM users WHERE id='%s'", id))
	s.logger.Log(dbResult)
}

// 要运行此示例：
// 创建新目录并复制必要的代码，然后取消注释下面的 main 函数

/*
func main() {
	fmt.Println("=== fx 生命周期管理示例 ===\n")

	app := fx.New(
		fx.Provide(
			NewSimpleLogger,
			NewDatabase,
			NewCache,
			NewUserService,
		),
		fx.Invoke(func(service *UserService) {
			// 等待所有组件启动完成
			time.Sleep(1 * time.Second)

			fmt.Println("\n=== 应用程序运行中 ===")
			service.GetUser("user123")
			service.GetUser("user456")
		}),
	)

	// 启动应用（会按顺序调用所有 OnStart 钩子）
	app.Run()
	// 停止应用（会按相反顺序调用所有 OnStop 钩子）
}
*/
