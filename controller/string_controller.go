package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/renoz/toolbox-api/model"
	"github.com/renoz/toolbox-api/service"
)

// 字符串分割
func SplitHandler(c *gin.Context) {
	var req model.SplitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 4001, "参数验证错误: "+err.Error())
		return
	}
	
	// 验证必要参数
	if req.Input == "" || req.Delimiter == "" {
		responseError(c, 4001, "输入字符串和分隔符不能为空")
		return
	}
	
	// 验证键值对模式参数
	if req.MapFormat && req.KeyValueDelimiter == "" {
		responseError(c, 4001, "键值对模式下必须指定键值对分隔符")
		return
	}
	
	// 检查输入长度限制
	if len(req.Input) > 100000 {
		responseError(c, 400, "输入内容过长，最大支持100000字符")
		return
	}
	
	// 调用服务处理
	result, err := service.SplitString(req)
	if err != nil {
		responseError(c, 9000, "字符串分割失败: "+err.Error())
		return
	}
	
	responseSuccess(c, result)
}

// 索引切分字符串
func SplitIndexedHandler(c *gin.Context) {
	var req model.SplitIndexedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 4001, "参数验证错误: "+err.Error())
		return
	}
	
	// 检查输入长度限制
	if len(req.Content) > 100000 {
		responseError(c, 400, "输入内容过长，最大支持100000字符")
		return
	}
	
	// 调用服务处理
	result, err := service.SplitIndexedString(req)
	if err != nil {
		responseError(c, 9000, "字符串分割失败: "+err.Error())
		return
	}
	
	responseSuccess(c, result)
}

// 字符串替换
func ReplaceHandler(c *gin.Context) {
	var req model.ReplaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 4001, "参数验证错误: "+err.Error())
		return
	}
	
	// 检查输入长度限制
	if len(req.Input) > 100000 {
		responseError(c, 400, "输入内容过长，最大支持100000字符")
		return
	}
	
	// 调用服务处理
	result, err := service.ReplaceString(req)
	if err != nil {
		responseError(c, 9000, "字符串替换失败: "+err.Error())
		return
	}
	
	responseSuccess(c, result)
}

// 命名格式转换处理函数
func CaseConversionHandler(c *gin.Context) {
	var req model.CaseConversionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 4001, "参数验证错误: "+err.Error())
		return
	}
	
	// 检查输入长度限制
	if len(req.Text) > 100000 {
		responseError(c, 400, "输入内容过长，最大支持100000字符")
		return
	}
	
	// 获取转换类型
	convType := c.Param("type")
	
	var result *model.CaseConversionResponse
	var err error
	
	// 根据转换类型调用相应函数
	switch convType {
	case "camel":
		result, err = service.ToCamelCase(req)
	case "pascal":
		result, err = service.ToPascalCase(req)
	case "snake":
		result, err = service.ToSnakeCase(req)
	case "kebab":
		result, err = service.ToKebabCase(req)
	default:
		responseError(c, 4001, "不支持的转换类型，可选值: camel, pascal, snake, kebab")
		return
	}
	
	if err != nil {
		responseError(c, 9000, "命名格式转换失败: "+err.Error())
		return
	}
	
	responseSuccess(c, result)
}

// 中文拼音首字母提取
func ExtractInitialsHandler(c *gin.Context) {
	var req model.ExtractInitialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 4001, "参数验证错误: "+err.Error())
		return
	}
	
	// 检查输入长度限制
	if len(req.Text) > 100000 {
		responseError(c, 400, "输入内容过长，最大支持100000字符")
		return
	}
	
	// 设置默认值
	// Uppercase默认为true
	if req.Uppercase == false {
		// 请求中明确指定为false时才设为false，否则保持默认值true
		req.Uppercase = false
	} else {
		req.Uppercase = true
	}
	
	// 调用服务处理
	result, err := service.ExtractInitials(req)
	if err != nil {
		responseError(c, 9000, "提取首字母失败: "+err.Error())
		return
	}
	
	responseSuccess(c, result)
}

// 日期转换为中文大写格式
func ConvertDateHandler(c *gin.Context) {
	var req model.ConvertDateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 4001, "参数验证错误: "+err.Error())
		return
	}
	
	// 调用服务处理
	result, err := service.ConvertDateToChinese(req)
	if err != nil {
		if e, ok := err.(*model.ErrorResponse); ok {
			responseError(c, e.Code, e.Message)
		} else {
			responseError(c, 9000, "日期转换失败: "+err.Error())
		}
		return
	}
	
	responseSuccess(c, result)
}

// 日期转换为中文普通格式
func ConvertDateSimpleHandler(c *gin.Context) {
	var req model.ConvertDateSimpleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, 4001, "参数验证错误: "+err.Error())
		return
	}
	
	// 调用服务处理
	result, err := service.ConvertDateToChineseSimple(req)
	if err != nil {
		if e, ok := err.(*model.ErrorResponse); ok {
			responseError(c, e.Code, e.Message)
		} else {
			responseError(c, 9000, "日期转换失败: "+err.Error())
		}
		return
	}
	
	responseSuccess(c, result)
} 