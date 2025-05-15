package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/renoz/toolbox-api/config"
	"github.com/renoz/toolbox-api/router"
	"log"
)

func main() {
	// 只设置启动模式为Release，这只会禁用启动时的调试信息
	// 但会保留API请求的访问日志
	gin.SetMode(gin.ReleaseMode)
	
	// 禁用Gin的控制台彩色输出
	gin.DisableConsoleColor()
	
	// 加载配置
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}
	
	// 设置路由
	r := router.SetupRouter()
	
	// 获取端口
	port := config.GetPort()
	
	// 启动服务
	log.Printf("工具箱API服务启动, 监听端口: %s", port)
	err = r.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
} 