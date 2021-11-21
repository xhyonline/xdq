package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xdq/controllers"
)

// InitRouter 初始化路由
func InitRouter(engine *gin.Engine) {
	// 推送数据
	engine.POST("/push", controllers.Push)
	// 删除主题
	engine.GET("/delete/:name", controllers.DeleteTopic)
	// 获取所有主题名
	engine.GET("/get/topics", controllers.GetTopic)
	// 获取某主题下等待被消费的数据
	engine.GET("/get/topic/:name", controllers.GetWaitDataByTopic)
}
