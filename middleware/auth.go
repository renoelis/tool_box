package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/renoz/toolbox-api/config"
	"net/http"
	"strings"
)

// 免认证路径列表
var exemptPaths = []string{
	"/toolbox/system/token",        // token查询接口不需要认证
	"/toolbox/system/token/refresh", // token刷新接口不需要认证
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否为免认证路径
		path := c.Request.URL.Path
		for _, exemptPath := range exemptPaths {
			if strings.EqualFold(path, exemptPath) {
				// 免认证路径，直接放行
				c.Next()
				return
			}
		}
		
		// 检查是否启用token验证
		if !config.IsTokenEnabled() {
			c.Next()
			return
		}

		// 从Header获取accessToken
		token := c.GetHeader("accessToken")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    4030,
				"message": "缺少访问令牌，请在请求头中设置accessToken",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 验证token
		validToken := config.GetToken()
		if token != validToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    4031,
				"message": "无效的访问令牌",
				"data":    nil,
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
} 