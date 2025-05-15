package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/renoz/toolbox-api/model"
	"github.com/renoz/toolbox-api/service"
	"strconv"
)

// 工作日计算
func WorkdayRangeHandler(c *gin.Context) {
	var req model.WorkdayRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 4001, "参数验证错误: "+err.Error())
		return
	}
	
	// 调用服务处理
	result, err := service.CalculateWorkdays(req)
	if err != nil {
		if e, ok := err.(*model.ErrorResponse); ok {
			responseError(c, e.Code, e.Message)
		} else {
			responseError(c, 9000, "计算工作日时发生错误: "+err.Error())
		}
		return
	}
	
	responseSuccess(c, result)
}

// 获取当前时间
func CurrentTimeHandler(c *gin.Context) {
	var req struct {
		Format       string `json:"format" form:"format"`
		Timezone     string `json:"timezone" form:"timezone"`
		TzOffset     *int   `json:"tz_offset" form:"tz_offset"`
		CustomFormat string `json:"custom_format" form:"custom_format"`
	}

	// 绑定JSON参数
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 1001, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务函数获取当前时间
	response, err := service.GetCurrentTime(req.Format, req.Timezone, req.TzOffset, req.CustomFormat)
	if err != nil {
		if e, ok := err.(*model.ErrorResponse); ok {
			responseError(c, e.Code, e.Message)
		} else {
			responseError(c, 9000, "获取当前时间时发生错误: "+err.Error())
		}
		return
	}

	responseSuccess(c, response)
}

// 检查是否为周末
func IsWeekendHandler(c *gin.Context) {
	// 解析参数
	dateStr := c.DefaultQuery("date", "")
	
	// 调用服务处理
	result, err := service.CheckIsWeekend(dateStr)
	if err != nil {
		if e, ok := err.(*model.ErrorResponse); ok {
			responseError(c, e.Code, e.Message)
		} else {
			responseError(c, 9000, "检查周末时发生错误: "+err.Error())
		}
		return
	}
	
	responseSuccess(c, result)
}

// 时间格式转换
func TimeConvertHandler(c *gin.Context) {
	var req model.TimeConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 4001, "参数验证错误: "+err.Error())
		return
	}
	
	// 调用服务处理
	result, err := service.ConvertTime(req)
	if err != nil {
		if e, ok := err.(*model.ErrorResponse); ok {
			responseError(c, e.Code, e.Message)
		} else {
			responseError(c, 9000, "时间转换失败: "+err.Error())
		}
		return
	}
	
	responseSuccess(c, result)
}

// 获取周数信息
func WeekNumberHandler(c *gin.Context) {
	// 解析参数
	dateStr := c.DefaultQuery("date", "")
	
	// 解析是否计算月内周数
	inMonthStr := c.DefaultQuery("in_month", "false")
	inMonth := false
	if inMonthStr == "true" {
		inMonth = true
	}
	
	// 解析周起始日
	startWithMondayStr := c.DefaultQuery("start_with_monday", "true")
	startWithMonday := true
	if startWithMondayStr == "false" {
		startWithMonday = false
	}
	
	// 调用服务处理
	result, err := service.GetWeekNumber(dateStr, inMonth, startWithMonday)
	if err != nil {
		if e, ok := err.(*model.ErrorResponse); ok {
			responseError(c, e.Code, e.Message)
		} else {
			responseError(c, 9000, "获取周数信息失败: "+err.Error())
		}
		return
	}
	
	responseSuccess(c, result)
}

// 获取时区信息
func TimezoneInfoHandler(c *gin.Context) {
	// 解析参数
	timezone := c.DefaultQuery("timezone", "")
	
	// 解析时区偏移量
	var tzOffset *int
	tzOffsetStr := c.Query("tz_offset")
	if tzOffsetStr != "" {
		offset, err := strconv.Atoi(tzOffsetStr)
		if err == nil {
			tzOffset = &offset
		}
	}
	
	// 调用服务处理
	result, err := service.GetTimezoneInfo(timezone, tzOffset)
	if err != nil {
		if e, ok := err.(*model.ErrorResponse); ok {
			responseError(c, e.Code, e.Message)
		} else {
			responseError(c, 9000, "获取时区信息失败: "+err.Error())
		}
		return
	}
	
	responseSuccess(c, result)
} 