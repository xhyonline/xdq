package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xdq/library"
	"github.com/xhyonline/xdq/services"
	"github.com/xhyonline/xutil/g"
)

// Push 推送数据
func Push(context *gin.Context) {
	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusOK, g.R(library.Error, "参数错误"+err.Error(), nil))
		return
	}
	params := new(services.PushParams)
	if err := json.Unmarshal(body, params); err != nil {
		context.JSON(http.StatusOK, g.R(library.Error, "参数错误"+err.Error(), nil))
		return
	}
	if err := services.Push(params); err != nil {
		context.JSON(http.StatusOK, g.R(library.Error, err.Error(), nil))
		return
	}
	context.JSON(http.StatusOK, g.R(library.Success, "", nil))
}

// DeleteTopic 删除主题
func DeleteTopic(context *gin.Context) {
	topic := context.Param("name")
	if err := services.DeleteTopic(topic); err != nil {
		context.JSON(http.StatusOK, g.R(library.Error, "删除失败"+err.Error(), nil))
		return
	}
	context.JSON(http.StatusOK, g.R(library.Success, "", nil))
}

// GetTopic 获取主题
func GetTopic(context *gin.Context) {
	context.JSON(http.StatusOK, g.R(library.Success, "", services.GetTopics()))
}

// GetWaitDataByTopic 获取等待被消费的数据
func GetWaitDataByTopic(context *gin.Context) {
	topic := context.Param("name")
	detail, err := services.GetWaitDataByTopic(topic)
	if err != nil {
		context.JSON(http.StatusOK, g.R(library.Error, err.Error(), nil))
		return
	}
	context.JSON(http.StatusOK, g.R(library.Success, "", detail))
}
