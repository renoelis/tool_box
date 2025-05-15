package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/renoz/toolbox-api/model"
	"github.com/renoz/toolbox-api/service"
	"strconv"
)

// 生成随机整数
func RandomIntegerHandler(c *gin.Context) {
	// 解析参数
	minStr := c.Query("min")
	maxStr := c.Query("max")
	countStr := c.Query("count")
	allowDuplicatesStr := c.DefaultQuery("allow_duplicates", "true")
	
	// 验证参数
	min, err := strconv.Atoi(minStr)
	if err != nil {
		responseError(c, 400, "min参数必须是整数")
		return
	}
	
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		responseError(c, 400, "max参数必须是整数")
		return
	}
	
	count, err := strconv.Atoi(countStr)
	if err != nil {
		responseError(c, 400, "count参数必须是整数")
		return
	}
	
	allowDuplicates := true
	if allowDuplicatesStr == "false" {
		allowDuplicates = false
	}
	
	// 调用服务处理
	result, err := service.GenerateRandomIntegers(min, max, count, allowDuplicates)
	if err != nil {
		if e, ok := err.(*model.ErrorResponse); ok {
			responseError(c, e.Code, e.Message)
		} else {
			responseError(c, 9000, "生成随机数失败: "+err.Error())
		}
		return
	}
	
	responseSuccess(c, result)
} 