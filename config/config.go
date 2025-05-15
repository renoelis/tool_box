package config

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 应用配置结构
type Config struct {
	Port         string `json:"port"`
	RateLimit    int    `json:"rate_limit"`    // 每秒请求限制
	TokenEnabled bool   `json:"token_enabled"` // 是否启用token验证
}

// TokenConfig token配置
type TokenConfig struct {
	Token string `json:"token"`
}

// 默认配置
var defaultConfig = Config{
	Port:         "4005",
	RateLimit:    240, // 每秒240次请求
	TokenEnabled: true,
}

// 全局配置实例
var AppConfig Config
var TokenConf TokenConfig

// LoadConfig 加载配置
func LoadConfig() error {
	// 初始化默认配置
	AppConfig = defaultConfig

	// 确保toolbox_data目录存在
	if err := ensureDataDir(); err != nil {
		return err
	}

	// 如果存在配置文件则加载
	configPath := "config.json"
	if _, err := os.Stat(configPath); err == nil {
		configFile, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(configFile, &AppConfig)
		if err != nil {
			return err
		}
	} else {
		// 创建默认配置文件
		configData, err := json.MarshalIndent(AppConfig, "", "  ")
		if err != nil {
			return err
		}
		err = os.WriteFile(configPath, configData, 0644)
		if err != nil {
			return err
		}
	}

	// 加载或创建token配置
	tokenPath := "toolbox_data/toolbox_token.json"
	if _, err := os.Stat(tokenPath); err == nil {
		tokenFile, err := os.ReadFile(tokenPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(tokenFile, &TokenConf)
		if err != nil {
			return err
		}
	} else {
		// 生成默认token
		TokenConf = TokenConfig{
			Token: generateRandomToken(),
		}
		tokenData, err := json.MarshalIndent(TokenConf, "", "  ")
		if err != nil {
			return err
		}
		err = os.WriteFile(tokenPath, tokenData, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// 确保toolbox_data目录存在
func ensureDataDir() error {
	if _, err := os.Stat("toolbox_data"); os.IsNotExist(err) {
		err = os.Mkdir("toolbox_data", 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// 生成随机token
func generateRandomToken() string {
	// 生成32字节的随机数据
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// 如果随机数生成失败，则使用默认token
		return "toolbox-api-default-token-2024"
	}
	
	// 将随机字节转换为base64字符串
	token := base64.StdEncoding.EncodeToString(randomBytes)
	
	// 返回格式化的token
	return "toolbox-api-" + token[:24]
}

// GetPort 获取服务端口
func GetPort() string {
	return AppConfig.Port
}

// GetRateLimit 获取限流设置
func GetRateLimit() int {
	return AppConfig.RateLimit
}

// GetRootPath 获取项目根路径
func GetRootPath() (string, error) {
	return filepath.Abs(".")
}

// IsTokenEnabled 是否启用token验证
func IsTokenEnabled() bool {
	return AppConfig.TokenEnabled
}

// GetToken 获取有效token
func GetToken() string {
	return TokenConf.Token
}

// UpdateToken 更新Token
func UpdateToken(newToken string) {
	// 更新内存中的token
	TokenConf.Token = newToken
} 