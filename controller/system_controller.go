package controller

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/renoz/toolbox-api/config"
	"os"
	"path/filepath"
)

// Token信息结构
type TokenInfo struct {
	Token string `json:"token"`
}

// 获取部署Token
func GetDeployTokenHandler(c *gin.Context) {
	// 获取token文件路径
	// 尝试相对路径
	tokenFile := "toolbox_data/toolbox_token.json"
	
	// 如果文件不存在，尝试从上一级目录查找
	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		exePath, _ := os.Executable()
		exeDir := filepath.Dir(exePath)
		tokenFile = filepath.Join(exeDir, "toolbox_data/toolbox_token.json")
	}
	
	// 读取token文件
	fileData, err := os.ReadFile(tokenFile)
	if err != nil {
		responseError(c, 9001, "无法读取Token信息: "+err.Error())
		return
	}
	
	// 解析JSON
	var tokenInfo TokenInfo
	if err := json.Unmarshal(fileData, &tokenInfo); err != nil {
		responseError(c, 9002, "Token信息解析失败: "+err.Error())
		return
	}
	
	// 返回token信息
	responseSuccess(c, gin.H{
		"token": tokenInfo.Token,
		"description": "项目部署Token，用于API认证",
	})
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

// 重置刷新Token
func RefreshTokenHandler(c *gin.Context) {
	// 获取token文件路径
	tokenFile := "toolbox_data/toolbox_token.json"
	
	// 如果文件不存在，尝试从上一级目录查找
	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		exePath, _ := os.Executable()
		exeDir := filepath.Dir(exePath)
		tokenFile = filepath.Join(exeDir, "toolbox_data/toolbox_token.json")
	}
	
	// 先保存当前的token作为旧token
	oldToken := config.GetToken()
	
	// 生成新的token
	newToken := generateRandomToken()
	
	// 确保新token和旧token不同
	for newToken == oldToken {
		newToken = generateRandomToken()
	}
	
	// 创建新的token配置
	tokenInfo := TokenInfo{
		Token: newToken,
	}
	
	// 将新token写入配置文件
	tokenData, err := json.MarshalIndent(tokenInfo, "", "  ")
	if err != nil {
		responseError(c, 9003, "生成Token数据失败: "+err.Error())
		return
	}
	
	// 确保toolbox_data目录存在
	if err := ensureDataDir(); err != nil {
		responseError(c, 9005, "创建数据目录失败: "+err.Error())
		return
	}
	
	err = os.WriteFile(tokenFile, tokenData, 0644)
	if err != nil {
		responseError(c, 9004, "保存Token数据失败: "+err.Error())
		return
	}
	
	// 更新config包中的token配置
	config.UpdateToken(newToken)
	
	// 返回新的token信息
	responseSuccess(c, gin.H{
		"new_token": newToken,
		"message": "Token已成功重置，请使用新Token进行API认证",
	})
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