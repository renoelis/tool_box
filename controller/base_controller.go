package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/renoz/toolbox-api/model"
	"net/http"
)

// 通用响应封装函数
func responseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, model.Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// 通用错误响应
func responseError(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, model.Response{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
} 