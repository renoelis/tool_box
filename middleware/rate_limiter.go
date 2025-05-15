package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/renoz/toolbox-api/config"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

// 限流器客户端的过期时间
const (
	cleanupInterval = 10 * time.Minute  // 清理间隔
	clientExpiry    = 30 * time.Minute  // 客户端过期时间
)

// 客户端限流器
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// 全局限流器
var (
	clients   = make(map[string]*client)
	clientsMu sync.RWMutex
	cleanupOnce sync.Once
)

// 启动清理goroutine
func startCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		
		for range ticker.C {
			// 使用当前时间
			now := time.Now()
			
			// 加锁清理
			clientsMu.Lock()
			for ip, cl := range clients {
				if now.Sub(cl.lastSeen) > clientExpiry {
					delete(clients, ip)
				}
			}
			clientsMu.Unlock()
		}
	}()
}

// RateLimiterMiddleware 请求频率限制中间件
func RateLimiterMiddleware() gin.HandlerFunc {
	// 获取配置的请求速率限制
	rateLimit := config.GetRateLimit()

	// 仅启动一次清理例程
	cleanupOnce.Do(startCleanupRoutine)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		// 使用读锁检查客户端是否存在
		clientsMu.RLock()
		cl, exists := clients[ip]
		clientsMu.RUnlock()
		
		if !exists {
			// 需要创建新客户端，使用写锁
			clientsMu.Lock()
			// 双重检查 - 避免在获取写锁期间其他goroutine已创建
			if cl, exists = clients[ip]; !exists {
				cl = &client{
					limiter:  rate.NewLimiter(rate.Limit(rateLimit), rateLimit),
					lastSeen: time.Now(),
				}
				clients[ip] = cl
			}
			clientsMu.Unlock()
		} else {
			// 仅更新访问时间，使用写锁
			clientsMu.Lock()
			cl.lastSeen = time.Now()
			clientsMu.Unlock()
		}

		// 检查是否允许当前请求
		if !cl.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    4290,
				"message": "请求频率过高，请稍后再试",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
} 