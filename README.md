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
go build -o e2b2api main.go

# 运行
./e2b2api
```

服务默认在`http://localhost:8080`上运行。

## API使用

### 获取可用模型列表

```
GET /v1/models
```

**curl示例:**
```bash
curl -X GET http://localhost:8080/v1/models \
  -H "Authorization: Bearer sk-123456"
```

### 聊天完成请求

```
POST /v1/chat/completions
```

**请求格式示例:**
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

**curl示例 (普通请求):**
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-123456" \
  -d '{
    "model": "claude-3-7-sonnet-latest",
    "messages": [
      {"role": "system", "content": "你是一个有用的AI助手。"},
      {"role": "user", "content": "你好，介绍一下自己"}
    ],
    "temperature": 0.7,
    "max_tokens": 500
  }'
```

**curl示例 (流式请求):**
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-123456" \
  --no-buffer \
  -d '{
    "model": "claude-3-7-sonnet-latest",
    "messages": [
      {"role": "system", "content": "你是一个有用的AI助手。"},
      {"role": "user", "content": "你好，介绍一下自己"}
    ],
    "temperature": 0.7,
    "max_tokens": 500,
    "stream": true
  }'
```

## 开发者集成示例

以下是几种常用编程语言的集成示例，展示如何在您的应用中调用E2B API Gateway。

### JavaScript (Node.js)

```javascript
// 普通请求示例
async function sendChatRequest() {
  const response = await fetch('http://localhost:8080/v1/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer sk-123456'
    },
    body: JSON.stringify({
      model: 'claude-3-7-sonnet-latest',
      messages: [
        { role: 'system', content: '你是一个有用的AI助手。' },
        { role: 'user', content: '你好，介绍一下自己' }
      ],
      temperature: 0.7,
      max_tokens: 500
    })
  });
  
  const data = await response.json();
  console.log(data.choices[0].message.content);
}

// 流式请求示例
async function streamChatRequest() {
  const response = await fetch('http://localhost:8080/v1/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer sk-123456'
    },
    body: JSON.stringify({
      model: 'claude-3-7-sonnet-latest',
      messages: [
        { role: 'system', content: '你是一个有用的AI助手。' },
        { role: 'user', content: '你好，介绍一下自己' }
      ],
      temperature: 0.7,
      max_tokens: 500,
      stream: true
    })
  });
  
  const reader = response.body.getReader();
  const decoder = new TextDecoder('utf-8');
  
  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    
    const chunk = decoder.decode(value);
    const lines = chunk.split('\n\n');
    
    for (const line of lines) {
      if (line.startsWith('data: ') && line !== 'data: [DONE]') {
        const data = JSON.parse(line.substring(6));
        if (data.choices[0].delta.content) {
          process.stdout.write(data.choices[0].delta.content);
        }
      }
    }
  }
}
```

### Python

```python
import requests
import json
import sseclient

# 普通请求示例
def send_chat_request():
    url = 'http://localhost:8080/v1/chat/completions'
    headers = {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer sk-123456'
    }
    data = {
        'model': 'claude-3-7-sonnet-latest',
        'messages': [
            {'role': 'system', 'content': '你是一个有用的AI助手。'},
            {'role': 'user', 'content': '你好，介绍一下自己'}
        ],
        'temperature': 0.7,
        'max_tokens': 500
    }
    
    response = requests.post(url, headers=headers, json=data)
    result = response.json()
    print(result['choices'][0]['message']['content'])

# 流式请求示例
def stream_chat_request():
    url = 'http://localhost:8080/v1/chat/completions'
    headers = {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer sk-123456'
    }
    data = {
        'model': 'claude-3-7-sonnet-latest',
        'messages': [
            {'role': 'system', 'content': '你是一个有用的AI助手。'},
            {'role': 'user', 'content': '你好，介绍一下自己'}
        ],
        'temperature': 0.7,
        'max_tokens': 500,
        'stream': True
    }
    
    response = requests.post(url, headers=headers, json=data, stream=True)
    client = sseclient.SSEClient(response)
    
    for event in client.events():
        if event.data != '[DONE]':
            chunk = json.loads(event.data)
            if chunk['choices'][0]['delta'].get('content'):
                print(chunk['choices'][0]['delta']['content'], end='', flush=True)
```

### Go

```go
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 聊天请求结构
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 聊天响应结构
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// 普通请求示例
func sendChatRequest() error {
	url := "http://localhost:8080/v1/chat/completions"
	
	chatReq := ChatRequest{
		Model: "claude-3-7-sonnet-latest",
		Messages: []ChatMessage{
			{Role: "system", Content: "你是一个有用的AI助手。"},
			{Role: "user", Content: "你好，介绍一下自己"},
		},
		Temperature: 0.7,
		MaxTokens:   500,
	}
	
	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-123456")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return err
	}
	
	fmt.Println(chatResp.Choices[0].Message.Content)
	return nil
}

// 流式请求示例
func streamChatRequest() error {
	url := "http://localhost:8080/v1/chat/completions"
	
	chatReq := ChatRequest{
		Model: "claude-3-7-sonnet-latest",
		Messages: []ChatMessage{
			{Role: "system", Content: "你是一个有用的AI助手。"},
			{Role: "user", Content: "你好，介绍一下自己"},
		},
		Temperature: 0.7,
		MaxTokens:   500,
		Stream:      true,
	}
	
	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-123456")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// 处理SSE流
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}
			
			var streamResp map[string]interface{}
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}
			
			choices, ok := streamResp["choices"].([]interface{})
			if !ok || len(choices) == 0 {
				continue
			}
			
			choice, ok := choices[0].(map[string]interface{})
			if !ok {
				continue
			}
			
			delta, ok := choice["delta"].(map[string]interface{})
			if !ok {
				continue
			}
			
			content, ok := delta["content"].(string)
			if ok {
				fmt.Print(content)
			}
		}
	}
	
	return nil
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

## Docker 部署

本项目提供了完整的 Docker 支持，您可以使用 Docker 来快速部署此服务。

### 使用预构建镜像

可以直接从 Docker Hub 或 GitHub Container Registry 拉取预构建镜像:

```bash
# 从 Docker Hub 拉取
docker pull lmyself/e2b2api-go:latest
```

### 运行 Docker 容器

```bash
# 运行容器，映射8080端口，并设置API密钥
docker run -d \
  --name e2b2api \
  -p 8080:8080 \
  -e E2B_API_KEY="your-api-key" \
  lmyself/e2b2api-go:latest
```

### 使用自定义配置

```bash
# 使用自定义.env文件运行
docker run -d \
  --name e2b2api \
  -p 8080:8080 \
  -v $(pwd)/.env:/app/.env \
  lmyself/e2b2api-go:latest
```

### 使用Docker Compose

项目提供了docker-compose.yml文件，可以更简单地部署服务：

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

您可以通过在启动前设置环境变量来配置服务：

```bash
# 设置API密钥
export E2B_API_KEY="your-api-key"

# 启动服务
docker-compose up -d
```

### 从源码构建镜像

如果您想自行构建 Docker 镜像，可以执行:

```bash
# 构建镜像
docker build -t e2b2api-go:latest .

# 运行容器
docker run -d \
  --name e2b2api \
  -p 8080:8080 \
  -e E2B_API_KEY="your-api-key" \
  e2b2api-go:latest
```

# 单行命令
```
docker run -d --name e2b2api -p 8080:8080 -v /your_file_path/.env:/app/.env lmyself/e2b2api-go:latest
```

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