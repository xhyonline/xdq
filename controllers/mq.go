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
		context.JSON(http.StatusOK, g.R(library.Error, "参数错误"+err.Error(), nil))
		return
	}
	context.JSON(http.StatusOK, g.R(library.Success, "", nil))
}

// Get
func Get(context *gin.Context) {

}
