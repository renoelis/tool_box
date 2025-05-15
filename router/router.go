package router

import (
	"github.com/gin-gonic/gin"
	"github.com/renoz/toolbox-api/controller"
	"github.com/renoz/toolbox-api/middleware"
	"net/http"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	
	// 添加中间件
	r.Use(middleware.RateLimiterMiddleware())
	r.Use(middleware.AuthMiddleware())
	r.Use(middleware.QueryCaseMiddleware())
	
	// 处理404错误 - 将NoRoute移到这里，因为RouterGroup不支持NoRoute
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "接口不存在",
			"data":    nil,
		})
	})
	
	// 创建路由组，添加统一前缀
	api := r.Group("/toolbox")
	{
		// 系统相关接口（以system/开头）
		systemGroup := api.Group("/system")
		{
			systemGroup.GET("/token", controller.GetDeployTokenHandler)
			systemGroup.POST("/token/refresh", controller.RefreshTokenHandler)
			
			// 添加非GET方法的处理
			notSupportedHandler := func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":    4005,
					"message": "请求方法错误: 系统信息接口仅支持GET请求",
					"data":    nil,
				})
			}
			systemGroup.POST("/token", notSupportedHandler)
			systemGroup.PUT("/token", notSupportedHandler)
			systemGroup.DELETE("/token", notSupportedHandler)
			systemGroup.PATCH("/token", notSupportedHandler)
			systemGroup.OPTIONS("/token", notSupportedHandler)
		}
		
		// 字符串处理相关路由（以string/开头）
		stringGroup := api.Group("/string")
		{
			stringGroup.POST("/split", controller.SplitHandler)
			stringGroup.POST("/split-indexed", controller.SplitIndexedHandler)
			stringGroup.POST("/replace", controller.ReplaceHandler)
			stringGroup.POST("/case-conversion/:type", controller.CaseConversionHandler)
			stringGroup.POST("/extract-initials", controller.ExtractInitialsHandler)
			stringGroup.POST("/convert-date", controller.ConvertDateHandler)
			stringGroup.POST("/convert-date-simple", controller.ConvertDateSimpleHandler)
			
			// 添加非POST方法的处理，为每个HTTP方法单独设置处理函数
			notSupportedHandler := func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":    4005,
					"message": "请求方法错误: 字符串相关接口仅支持POST请求",
					"data":    nil,
				})
			}
			stringGroup.GET("/*path", notSupportedHandler)
			stringGroup.PUT("/*path", notSupportedHandler)
			stringGroup.DELETE("/*path", notSupportedHandler)
			stringGroup.PATCH("/*path", notSupportedHandler)
			stringGroup.OPTIONS("/*path", notSupportedHandler)
		}
		
		// 随机数生成相关路由（以random/开头）
		randomGroup := api.Group("/random")
		{
			randomGroup.GET("/integer", controller.RandomIntegerHandler)
			
			// 添加非GET方法的处理，为每个HTTP方法单独设置处理函数
			notSupportedHandler := func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":    4005,
					"message": "请求方法错误: 随机数相关接口仅支持GET请求",
					"data":    nil,
				})
			}
			randomGroup.POST("/*path", notSupportedHandler)
			randomGroup.PUT("/*path", notSupportedHandler)
			randomGroup.DELETE("/*path", notSupportedHandler)
			randomGroup.PATCH("/*path", notSupportedHandler)
			randomGroup.OPTIONS("/*path", notSupportedHandler)
		}
		
		// 时间处理相关路由（以time/开头）
		timeGroup := api.Group("/time")
		{
			timeGroup.POST("/workday-range", controller.WorkdayRangeHandler)
			timeGroup.POST("/current", controller.CurrentTimeHandler)
			timeGroup.POST("/convert", controller.TimeConvertHandler)
			
			timeGroup.GET("/is-weekend", controller.IsWeekendHandler)
			timeGroup.GET("/week-number", controller.WeekNumberHandler)
			timeGroup.GET("/timezone-info", controller.TimezoneInfoHandler)
			
			// 为POST接口添加方法不支持的处理
			postNotSupportedHandler := func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":    4005,
					"message": "请求方法错误: 该接口仅支持POST请求",
					"data":    nil,
				})
			}
			
			// 为每个POST路由添加其他HTTP方法的处理
			timeGroup.GET("/workday-range", postNotSupportedHandler)
			timeGroup.PUT("/workday-range", postNotSupportedHandler)
			timeGroup.DELETE("/workday-range", postNotSupportedHandler)
			timeGroup.PATCH("/workday-range", postNotSupportedHandler)
			timeGroup.OPTIONS("/workday-range", postNotSupportedHandler)
			
			timeGroup.GET("/current", postNotSupportedHandler)
			timeGroup.PUT("/current", postNotSupportedHandler)
			timeGroup.DELETE("/current", postNotSupportedHandler)
			timeGroup.PATCH("/current", postNotSupportedHandler)
			timeGroup.OPTIONS("/current", postNotSupportedHandler)
			
			timeGroup.GET("/convert", postNotSupportedHandler)
			timeGroup.PUT("/convert", postNotSupportedHandler)
			timeGroup.DELETE("/convert", postNotSupportedHandler)
			timeGroup.PATCH("/convert", postNotSupportedHandler)
			timeGroup.OPTIONS("/convert", postNotSupportedHandler)
			
			// 为GET接口添加方法不支持的处理
			getNotSupportedHandler := func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":    4005,
					"message": "请求方法错误: 该接口仅支持GET请求",
					"data":    nil,
				})
			}
			
			// 为每个GET路由添加其他HTTP方法的处理
			timeGroup.POST("/is-weekend", getNotSupportedHandler)
			timeGroup.PUT("/is-weekend", getNotSupportedHandler)
			timeGroup.DELETE("/is-weekend", getNotSupportedHandler)
			timeGroup.PATCH("/is-weekend", getNotSupportedHandler)
			timeGroup.OPTIONS("/is-weekend", getNotSupportedHandler)
			
			timeGroup.POST("/week-number", getNotSupportedHandler)
			timeGroup.PUT("/week-number", getNotSupportedHandler)
			timeGroup.DELETE("/week-number", getNotSupportedHandler)
			timeGroup.PATCH("/week-number", getNotSupportedHandler)
			timeGroup.OPTIONS("/week-number", getNotSupportedHandler)
			
			timeGroup.POST("/timezone-info", getNotSupportedHandler)
			timeGroup.PUT("/timezone-info", getNotSupportedHandler)
			timeGroup.DELETE("/timezone-info", getNotSupportedHandler)
			timeGroup.PATCH("/timezone-info", getNotSupportedHandler)
			timeGroup.OPTIONS("/timezone-info", getNotSupportedHandler)
		}
	}
	
	return r
} 