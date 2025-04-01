# E2B API Gateway

这是一个用Go语言实现的API网关，用于代理LLM请求到E2B服务。它提供了与OpenAI API兼容的接口，支持多种大语言模型，包括Claude、GPT-4o、Gemini等。基于 Gin Web 框架构建，提供更高的性能和更好的 API 处理能力。

## 特性

- 支持多种LLM模型 (OpenAI, Anthropic, Google)
- 兼容OpenAI API格式
- 支持流式响应和普通响应
- 详细的日志记录
- 自动参数约束和验证
- 内置CORS支持
- 基于 Gin 框架，性能更佳
- 中间件支持，易于扩展
- .env文件配置支持，便于管理敏感信息
- 环境变量配置，无需修改代码
- API密钥掩码保护，不会在日志中暴露完整密钥

## 配置

### 配置文件 (.env)

服务支持通过 `.env` 文件进行配置，无需修改代码：

1. 复制 `.env.example` 文件为 `.env`
2. 根据需要修改 `.env` 文件中的配置项

`.env` 文件示例:
```
# E2B API Gateway 配置文件

# API密钥，用于访问E2B服务
E2B_API_KEY=sk-123456

# 服务运行端口，默认为8080
E2B_PORT=8080

# Gin框架模式: debug 或 release
GIN_MODE=release
```

### 环境变量

除了 `.env` 文件外，服务也支持通过环境变量进行配置：

- `E2B_API_KEY`: API密钥，用于访问E2B服务
- `E2B_PORT`: 服务运行端口，默认为"8080"

例如：
```bash
# Linux/macOS 设置环境变量
export E2B_API_KEY="your-api-key"
export E2B_PORT="3000"

# 然后运行服务
go run main.go
```

在Windows PowerShell中设置环境变量：
```powershell
# Windows PowerShell 设置环境变量
$env:E2B_API_KEY = "your-api-key"
$env:E2B_PORT = "3000"

# 然后运行服务
go run main.go
```

在Windows命令提示符(CMD)中设置环境变量：
```cmd
:: Windows CMD 设置环境变量
set E2B_API_KEY=your-api-key
set E2B_PORT=3000

:: 然后运行服务
go run main.go
```

### 代码配置

主要配置在代码顶部的`CONFIG`结构中，但推荐使用.env文件或环境变量进行配置：

1. API密钥: 优先使用环境变量`E2B_API_KEY`，否则使用代码中的默认值
2. E2B基础URL: 固定为"https://fragments.e2b.dev"，直接写在代码中，不可配置
3. 服务端口: 优先使用环境变量`E2B_PORT`，否则使用默认值"8080"
4. 根据需要调整重试参数和模型配置

## 安装依赖

```bash
# 安装依赖
go mod tidy
```

## 运行

```bash
# 编译
go build -o e2b-gateway main.go

# 运行
./e2b-gateway
```

服务默认在`http://localhost:8080`上运行。

## API使用

### 获取可用模型列表

```
GET /v1/models
```

### 聊天完成请求

```
POST /v1/chat/completions
```

请求格式示例:

```json
{
  "model": "claude-3-7-sonnet-latest",
  "messages": [
    {"role": "system", "content": "你是一个有用的AI助手。"},
    {"role": "user", "content": "你好，介绍一下自己"}
  ],
  "temperature": 0.7,
  "max_tokens": 500,
  "stream": true
}
```

## 注意

- 所有请求都需要包含`Authorization: Bearer YOUR_API_KEY`头。
- `YOUR_API_KEY`必须与配置中的`CONFIG.API.API_KEY`匹配。

## 支持的模型

本服务支持多种模型，包括但不限于：

- OpenAI: o1-preview, o3-mini, gpt-4o, gpt-4.5-preview, gpt-4-turbo
- Anthropic: claude-3-5-sonnet-latest, claude-3-7-sonnet-latest, claude-3-5-haiku-latest
- Google: gemini-1.5-pro, gemini-2.5-pro-exp-03-25, gemini-exp-1121, gemini-2.0-flash-exp

有关完整的模型列表，请参阅`/v1/models`响应。

## 扩展

通过 Gin 框架，可以轻松添加自定义中间件和路由，例如：

- 添加速率限制
- 实现自定义认证
- 添加请求日志中间件
- 实现指标收集
- 添加健康检查端点

## 安全特性

### .env 文件配置

使用 .env 文件进行配置有以下安全优势：

- 敏感信息（如 API 密钥）可以存储在本地文件中，不会暴露在代码或命令行历史中
- .env 文件可以被添加到 .gitignore 中，防止敏感信息被意外提交到版本控制系统
- 便于在不同环境中（开发、测试、生产）使用不同的配置
- 所有敏感配置集中在一个文件中，便于安全审计和管理

记得将 .env 文件添加到 .gitignore 文件中：
```
# .gitignore
.env
```

### API密钥保护

为了防止敏感信息泄漏，该服务在日志记录中实现了API密钥掩码保护：

- 在日志中，API密钥只会显示前几个字符，其余部分会被替换为`...`
- 这样可以确保在排查问题时知道使用了哪个密钥，同时不会泄露完整的密钥信息
- 同样的保护也应用于认证失败时的日志记录

示例日志输出：
```
[2023-08-15T12:34:56Z][f7a8b9c0] INFO: 服务配置信息:
[2023-08-15T12:34:56Z][f7a8b9c0] INFO: API_KEY: sk-12345...
[2023-08-15T12:34:56Z][f7a8b9c0] INFO: BASE_URL: https://fragments.e2b.dev
[2023-08-15T12:34:56Z][f7a8b9c0] INFO: 端口: 8080
``` 