package middleware

import (
	"net/http"
	"strings"

	"github.com/xhyonline/xdq/services"

	"github.com/xhyonline/xutil/helper"

	"github.com/gin-gonic/gin"
	"github.com/xhyonline/xdq/library"
	"github.com/xhyonline/xutil/g"
)

// Auth 鉴权
func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		var needAbort = true
		defer func() {
			if needAbort {
				context.Abort()
				return
			}
			context.Next()
		}()
		// isNotCheck 第三方回调,均不需要拦截
		if isNotCheck(context) {
			needAbort = false
			return
		}
		if context.Request.Header.Get("x-request-type") != "internal" {
			context.JSON(http.StatusOK, g.R(library.Error, "非法调用", nil))
			return
		}
		token := strings.Replace(context.Request.Header.Get("Authorization"), "Bearer ", "", 1)
		if _, err := services.ParseToken(token); err != nil {
			context.JSON(http.StatusOK, g.R(library.Error, err.Error(), nil))
			return
		}
		needAbort = false
	}
}

// isNotCheck 是否不需要检查
func isNotCheck(context *gin.Context) bool {
	return helper.InArray(context.Request.URL.Path, []string{})
}
