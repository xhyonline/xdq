package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xdq/controllers"
)

// InitRouter 初始化路由
func InitRouter(engine *gin.Engine) {
	// 推送数据
	engine.POST("/push", controllers.Push)
	// 消费数据 (删除)
	engine.POST("/delete", controllers.Push)
	// 获取某主题下可以消费的数据
	engine.GET("/get/ready/topic/:name", controllers.Push)
	// 获取某主题下所有的数据
	engine.GET("/get/topic/:name", controllers.Push)
	// 获取所有主题信息
	engine.GET("/get/topics", controllers.Push)
}
