# 工具箱API服务

这是一个多功能工具箱API服务，提供字符串处理、随机数生成和时间处理等功能。

## 功能特点

- 字符串处理服务：分割、替换、命名格式转换等
- 随机数生成服务：生成随机整数
- 时间处理服务：工作日计算、时区转换等
- 支持Token认证
- 请求频率限制(每秒240次)，IP级别隔离，支持自动清理过期记录

## 配置和部署

### 环境需求

- Go 1.21+
- Docker (可选)

### 本地运行

```bash
# 克隆仓库
git clone https://github.com/yourusername/toolbox-api.git
cd toolbox-api

# 安装依赖
go mod download

# 运行服务
go run cmd/main.go
```

### Docker部署

```bash
# 构建镜像
docker build -t toolbox-api:latest .

# 运行容器（带必要配置）
# 方式1：使用当前目录下的toolbox_data（确保你在项目根目录）
docker run -d \
  -p 4005:4005 \
  --name toolbox-api \
  -v ./toolbox_data:/app/toolbox_data \
  -e GIN_MODE=release \
  --restart=always \
  toolbox-api:latest

# 方式2：使用绝对路径（替换为你的实际路径）
# docker run -d \
#   -p 4005:4005 \
#   --name toolbox-api \
#   -v /path/to/your/toolbox-api/toolbox_data:/app/toolbox_data \
#   -e GIN_MODE=release \
#   --restart=always \
#   toolbox-api:latest

# 使用docker-compose
docker-compose -f docker-compose-toolbox-api.yml up -d
```

docker-compose-toolbox-api.yml 配置说明：
- 服务名：toolbox-api
- 容器名：toolbox-api
- 端口映射：4005:4005
- 卷挂载：./toolbox_data:/app/toolbox_data（保证Token持久化）
- 环境变量：GIN_MODE=release
- 重启策略：always
- 日志配置：JSON格式，单文件最大20MB，最多保留5个旧文件
- 网络：使用外部网络 api-proxy_proxy_net

### 导出镜像

```bash
# 保存镜像为tar文件
docker save -o toolbox-api.tar toolbox-api:latest
```

## API认证

API服务使用Token进行认证，默认启用。Token存储在`toolbox_data/toolbox_token.json`文件中，首次运行时会自动生成。

API调用必须在HTTP请求头中提供token：

```
accessToken: YOUR_TOKEN
```

所有没有提供有效accessToken的请求都将被拒绝访问。系统接口(/toolbox/system/*)除外。

## 主要API接口

所有API接口均以`/toolbox`为统一前缀，接口根据功能分组。

### 系统相关

- `GET /toolbox/system/token` - 获取部署Token
- `POST /toolbox/system/token/refresh` - 刷新Token

### 字符串处理

- `POST /toolbox/string/split` - 字符串分割
- `POST /toolbox/string/split-indexed` - 索引切分字符串
- `POST /toolbox/string/replace` - 字符串替换
- `POST /toolbox/string/case-conversion/{type}` - 命名格式转换
- `POST /toolbox/string/extract-initials` - 中文拼音首字母提取
- `POST /toolbox/string/convert-date` - 日期转换为中文大写格式
- `POST /toolbox/string/convert-date-simple` - 日期转换为中文普通格式

### 随机数生成

- `GET /toolbox/random/integer` - 生成随机整数

### 时间处理

- `POST /toolbox/time/workday-range` - 工作日计算
- `POST /toolbox/time/current` - 获取当前时间
- `POST /toolbox/time/convert` - 时间格式转换
- `GET /toolbox/time/is-weekend` - 检查是否为周末
- `GET /toolbox/time/week-number` - 获取周数信息
- `GET /toolbox/time/timezone-info` - 获取时区信息

## 错误处理

服务使用统一的JSON格式返回错误信息：

```json
{
  "code": 错误码,
  "message": "错误信息",
  "data": null
}
```

常见错误码：
- 404: 接口不存在
- 4005: 请求方法错误
- 4030: 缺少访问令牌
- 4031: 无效的访问令牌
- 4290: 请求频率过高

## 性能优化

- 速率限制使用令牌桶算法，支持每IP限流
- 非活跃IP的限流器数据会自动清理（30分钟无活动）
- 使用读写锁优化并发性能
- 支持请求参数字段名的驼峰命名和下划线命名自动转换