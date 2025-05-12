package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // 引入 CORS 中间件
	"ycd-platform/api"
	"ycd-platform/config"
)

func main() {
	// 1. 加载配置
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	log.Printf("📂 加载配置文件: %s", configPath)
	err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}
	log.Println("✅ 配置加载成功")

	// 2. 初始化 Fiber
	app := fiber.New()

	// 3. 启用 CORS 中间件
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // 允许的前端地址
		AllowMethods: "GET,POST,PUT,DELETE",   // 允许的 HTTP 方法
	}))

	// 4. 注册 API 路由
	apiGroup := app.Group("/api")

	// 健康检查路由
	healthGroup := apiGroup.Group("/health")
	healthGroup.Get("/", api.HealthCheck)

	// 项目相关路由
	projectGroup := apiGroup.Group("/projects")
	projectGroup.Get("/", api.GetProjects)
	projectGroup.Get("/environments", api.GetEnvironments)

	// 镜像相关路由
	imageGroup := apiGroup.Group("/images")
	imageGroup.Get("/", api.GetImages)

	// Kubernetes 相关路由
	k8sGroup := apiGroup.Group("/k8s")
	k8sGroup.Get("/clusters", api.GetClusters)
	k8sGroup.Get("/namespaces", api.GetNamespaces)

	// 部署相关路由
	deployGroup := apiGroup.Group("/deploy")
	deployGroup.Post("/", api.Deploy)
	deployGroup.Get("/status", api.GetDeploymentStatus)

	// 配置管理路由
	configGroup := apiGroup.Group("/config")
	configGroup.Get("/", api.GetConfig)                  // 获取所有配置
	configGroup.Put("/project", api.UpdateProjectConfig) // 更新项目配置
	configGroup.Post("/project", api.AddProjectConfig)   // 添加新项目
	configGroup.Delete("/project", api.DeleteProjectConfig) // 删除项目

	// 5. 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 服务启动在 http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
