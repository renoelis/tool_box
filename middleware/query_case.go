package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
)

// QueryCaseMiddleware 处理请求参数中的驼峰命名和下划线命名转换
func QueryCaseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只处理POST和PUT请求
		if c.Request.Method != "POST" && c.Request.Method != "PUT" {
			c.Next()
			return
		}
		
		// 检查Content-Type
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			c.Next()
			return
		}
		
		// 读取请求体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}
		
		// 关闭原始请求体
		c.Request.Body.Close()
		
		// 解析JSON
		var requestData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
			// 重置请求体并继续
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			c.Next()
			return
		}
		
		// 转换字段名称 - 支持驼峰命名映射到下划线命名
		normalizedRequestData := make(map[string]interface{})
		for key, value := range requestData {
			// 检查常见的驼峰命名模式
			if key == "mapFormat" {
				normalizedRequestData["map_format"] = value
				normalizedRequestData["mapFormat"] = value
			} else if key == "keyValueDelimiter" {
				normalizedRequestData["key_value_delimiter"] = value
				normalizedRequestData["keyValueDelimiter"] = value
			} else if key == "useRegex" {
				normalizedRequestData["use_regex"] = value
				normalizedRequestData["useRegex"] = value
			} else if key == "uppercase" || key == "Uppercase" {
				// 确保同时支持大小写两种形式
				normalizedRequestData["uppercase"] = value
				normalizedRequestData["Uppercase"] = value
			} else {
				// 保留原始键名
				normalizedRequestData[key] = value
			}
		}
		
		// 将转换后的数据重新编码为JSON
		newBodyBytes, err := json.Marshal(normalizedRequestData)
		if err != nil {
			// 出错时使用原始请求体
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		} else {
			// 替换请求体
			c.Request.Body = io.NopCloser(bytes.NewBuffer(newBodyBytes))
			// 更新内容长度
			c.Request.ContentLength = int64(len(newBodyBytes))
		}
		
		c.Next()
	}
} 