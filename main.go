package main

import (
	"net/http"

	"github.com/xhyonline/xdq/configs"
	"github.com/xhyonline/xdq/internal"
	"github.com/xhyonline/xdq/middleware"
	"github.com/xhyonline/xdq/router"

	"github.com/xhyonline/xutil/sig"


	"github.com/xhyonline/xdq/component"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	// 初始化配置
	configs.Init(configs.WithRedis())
	// 初始化 mysql 、redis 等服务组件
	component.Init(component.RegisterRedis())
	// 中间件
	g.Use(middleware.Cors())
	// 初始化路由
	router.InitRouter(g)
	// 启动 HTTP 服务
	httpServer := &internal.HTTPServer{Server: &http.Server{Addr: "0.0.0.0:8081", Handler: g}}
	go httpServer.Run()
	// 注册优雅退出
	ctx := sig.Get().RegisterClose(httpServer)
	<-ctx.Done()
}
