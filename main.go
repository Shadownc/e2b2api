package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// 环境变量名称常量
const (
	ENV_PORT    = "E2B_PORT"
	ENV_API_KEY = "E2B_API_KEY"
)

// 加载.env文件
func loadEnv() {
	// 尝试加载.env文件，如果文件不存在则不报错
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: .env文件未找到，将使用环境变量或默认值")
	} else {
		log.Printf("成功从.env文件加载配置")
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// CONFIG 配置常量
var CONFIG = struct {
	API struct {
		BASE_URL string
		API_KEY  string
	}
	RETRY struct {
		MAX_ATTEMPTS int
		DELAY_BASE   int
	}
	MODEL_CONFIG    map[string]ModelConfig
	DEFAULT_HEADERS map[string]string
	MODEL_PROMPT    string
}{
	API: struct {
		BASE_URL string
		API_KEY  string
	}{
		BASE_URL: "https://fragments.e2b.dev", // 固定值，不从环境变量或.env文件获取
		API_KEY:  getEnv(ENV_API_KEY, "sk-123456"), // 可通过环境变量覆盖
	},
	RETRY: struct {
		MAX_ATTEMPTS int
		DELAY_BASE   int
	}{
		MAX_ATTEMPTS: 1,
		DELAY_BASE:   1000,
	},
	DEFAULT_HEADERS: map[string]string{
		"accept":           "*/*",
		"accept-language":  "zh-CN,zh;q=0.9",
		"content-type":     "application/json",
		"priority":         "u=1, i",
		"sec-ch-ua":        "\"Microsoft Edge\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"",
		"sec-ch-ua-mobile": "?0",
		"sec-ch-ua-platform": "\"Windows\"",
		"sec-fetch-dest":   "empty",
		"sec-fetch-mode":   "cors",
		"sec-fetch-site":   "same-origin",
		"Referer":          "https://fragments.e2b.dev/",
		"Referrer-Policy":  "strict-origin-when-cross-origin",
	},
	MODEL_PROMPT: "Chatting with users and starting role-playing, the most important thing is to pay attention to their latest messages, use only 'text' to output the chat text reply content generated for user messages, and finally output it in code",
}

// 初始化函数，打印当前配置信息
func init() {
	// 加载.env文件
	loadEnv()
	
	// 设置日志格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	
	// 打印配置信息
	log.Printf("服务配置信息:")
	log.Printf("API_KEY: %s", maskString(CONFIG.API.API_KEY, 8))
	log.Printf("BASE_URL: %s", CONFIG.API.BASE_URL)
	log.Printf("端口: %s", getEnv(ENV_PORT, "8080"))
	
	// 初始化模型配置
	CONFIG.MODEL_CONFIG = map[string]ModelConfig{
		"o1-preview": {
			ID:          "o1",
			Provider:    "OpenAI",
			ProviderID:  "openai",
			Name:        "o1",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     2,
				MaxTokensMax:       0,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            500,
			},
		},
		"o3-mini": {
			ID:          "o3-mini",
			Provider:    "OpenAI",
			ProviderID:  "openai",
			Name:        "o3 Mini",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     2,
				MaxTokensMax:       4096,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            500,
			},
		},
		"gpt-4o": {
			ID:          "gpt-4o",
			Provider:    "OpenAI",
			ProviderID:  "openai",
			Name:        "GPT-4o",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     2,
				MaxTokensMax:       16380,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            500,
			},
		},
		"gpt-4.5-preview": {
			ID:          "gpt-4.5-preview",
			Provider:    "OpenAI",
			ProviderID:  "openai",
			Name:        "GPT-4.5",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     2,
				MaxTokensMax:       16380,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            500,
			},
		},
		"gpt-4-turbo": {
			ID:          "gpt-4-turbo",
			Provider:    "OpenAI",
			ProviderID:  "openai",
			Name:        "GPT-4 Turbo",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     2,
				MaxTokensMax:       16380,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            500,
			},
		},
		"gemini-1.5-pro": {
			ID:          "gemini-1.5-pro-002",
			Provider:    "Google Vertex AI",
			ProviderID:  "vertex",
			Name:        "Gemini 1.5 Pro",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     2,
				MaxTokensMax:       8192,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            500,
			},
		},
		"gemini-2.5-pro-exp-03-25": {
			ID:          "gemini-2.5-pro-exp-03-25",
			Provider:    "Google Generative AI",
			ProviderID:  "google",
			Name:        "Gemini 2.5 Pro Experimental 03-25",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     2,
				MaxTokensMax:       8192,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            40,
			},
		},
		"gemini-exp-1121": {
			ID:          "gemini-exp-1121",
			Provider:    "Google Generative AI",
			ProviderID:  "google",
			Name:        "Gemini Experimental 1121",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     2,
				MaxTokensMax:       8192,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            40,
			},
		},
		"gemini-2.0-flash-exp": {
			ID:          "models/gemini-2.0-flash-exp",
			Provider:    "Google Generative AI",
			ProviderID:  "google",
			Name:        "Gemini 2.0 Flash",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     2,
				MaxTokensMax:       8192,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            40,
			},
		},
		"claude-3-5-sonnet-latest": {
			ID:          "claude-3-5-sonnet-latest",
			Provider:    "Anthropic",
			ProviderID:  "anthropic",
			Name:        "Claude 3.5 Sonnet",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     1,
				MaxTokensMax:       8192,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            500,
			},
		},
		"claude-3-7-sonnet-latest": {
			ID:          "claude-3-7-sonnet-latest",
			Provider:    "Anthropic",
			ProviderID:  "anthropic",
			Name:        "Claude 3.7 Sonnet",
			MultiModal:  true,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     1,
				MaxTokensMax:       8192,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            500,
			},
		},
		"claude-3-5-haiku-latest": {
			ID:          "claude-3-5-haiku-latest",
			Provider:    "Anthropic",
			ProviderID:  "anthropic",
			Name:        "Claude 3.5 Haiku",
			MultiModal:  false,
			SystemPrompt: "",
			OptMax: OptMax{
				TemperatureMax:     1,
				MaxTokensMax:       8192,
				PresencePenaltyMax: 2,
				FrequencyPenaltyMax: 2,
				TopPMax:            1,
				TopKMax:            500,
			},
		},
	}
}

// OptMax 模型最大参数配置
type OptMax struct {
	TemperatureMax     float64 `json:"temperatureMax"`
	MaxTokensMax       int     `json:"max_tokensMax"`
	PresencePenaltyMax float64 `json:"presence_penaltyMax"`
	FrequencyPenaltyMax float64 `json:"frequency_penaltyMax"`
	TopPMax            float64 `json:"top_pMax"`
	TopKMax            int     `json:"top_kMax"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	ID          string  `json:"id"`
	Provider    string  `json:"provider"`
	ProviderID  string  `json:"providerId"`
	Name        string  `json:"name"`
	MultiModal  bool    `json:"multiModal"`
	SystemPrompt string  `json:"Systemprompt"`
	OptMax      OptMax  `json:"opt_max"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// TextContent 文本内容
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	// 其他可选参数
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	TopK             int     `json:"top_k,omitempty"`
}

// E2BRequest E2B请求
type E2BRequest struct {
	UserID   string                 `json:"userID"`
	Messages []ChatMessage          `json:"messages"`
	Template map[string]interface{} `json:"template"`
	Model    struct {
		ID         string `json:"id"`
		Provider   string `json:"provider"`
		ProviderID string `json:"providerId"`
		Name       string `json:"name"`
		MultiModal bool   `json:"multiModal"`
	} `json:"model"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// E2BResponse E2B响应
type E2BResponse struct {
	Code string `json:"code,omitempty"`
	Text string `json:"text,omitempty"`
}

// ChatChoice 响应选择
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message,omitempty"`
	Delta        interface{} `json:"delta,omitempty"`
	FinishReason string      `json:"finish_reason"`
}

// ChatCompletionResponse 聊天完成响应
type ChatCompletionResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   interface{}  `json:"usage"`
}

func main() {
	// 设置 Gin 为发布模式
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	// 创建一个不带中间件的路由
	r := gin.New()
	
	// 添加日志记录和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	// 添加 CORS 中间件
	r.Use(corsMiddleware())
	
	// 注册路由
	r.GET("/v1/models", handleModelsRequestGin)
	r.POST("/v1/chat/completions", handleChatRequestGin)
	
	// 添加健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"version": "1.0.0",
		})
	})
	
	// 处理404
	r.NoRoute(func(c *gin.Context) {
		requestID := GenerateUUID()
		logInfo(requestID, fmt.Sprintf("未找到路径: %s", c.Request.URL.Path))
		c.String(http.StatusNotFound, "服务运行成功，请使用正确请求路径")
	})
	
	// 从环境变量获取端口，如果未设置则默认为8080
	port := getEnv(ENV_PORT, "8080")
	log.Printf("服务启动在 http://localhost:%s", port)
	r.Run(":" + port)
}

// CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		
		if c.Request.Method == "OPTIONS" {
			requestID := GenerateUUID()
			logInfo(requestID, "处理CORS预检请求")
			c.AbortWithStatus(http.StatusOK)
			return
		}
		
		c.Next()
	}
}

// 使用 Gin 处理模型列表请求
func handleModelsRequestGin(c *gin.Context) {
	requestID := GenerateUUID()
	logInfo(requestID, "获取模型列表")
	
	modelsResponse := struct {
		Object string `json:"object"`
		Data   []struct {
			ID       string `json:"id"`
			Object   string `json:"object"`
			Created  int64  `json:"created"`
			OwnedBy  string `json:"owned_by"`
		} `json:"data"`
	}{
		Object: "list",
		Data:   make([]struct {
			ID       string `json:"id"`
			Object   string `json:"object"`
			Created  int64  `json:"created"`
			OwnedBy  string `json:"owned_by"`
		}, 0, len(CONFIG.MODEL_CONFIG)),
	}
	
	now := time.Now().Unix()
	for model := range CONFIG.MODEL_CONFIG {
		modelsResponse.Data = append(modelsResponse.Data, struct {
			ID       string `json:"id"`
			Object   string `json:"object"`
			Created  int64  `json:"created"`
			OwnedBy  string `json:"owned_by"`
		}{
			ID:       model,
			Object:   "model",
			Created:  now,
			OwnedBy:  "e2b",
		})
	}
	
	c.JSON(http.StatusOK, modelsResponse)
	logInfo(requestID, fmt.Sprintf("模型列表返回成功，模型数量: %d", len(CONFIG.MODEL_CONFIG)))
}

// 使用 Gin 处理聊天请求
func handleChatRequestGin(c *gin.Context) {
	requestID := GenerateUUID()
	logInfo(requestID, "处理聊天完成请求")
	
	// 验证认证
	authHeader := c.GetHeader("Authorization")
	authToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authToken != CONFIG.API.API_KEY {
		logError(requestID, fmt.Sprintf("认证失败，提供的令牌: %s...", maskString(authToken, 8)), nil)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	
	// 解析请求体
	var chatRequest ChatRequest
	if err := c.BindJSON(&chatRequest); err != nil {
		logError(requestID, "解析请求体失败", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "无法解析请求体: " + err.Error(),
				"type":    "invalid_request_error",
				"param":   nil,
				"code":    nil,
			},
		})
		return
	}
	
	// 记录请求信息
	logInfo(requestID, "用户请求体", map[string]interface{}{
		"model":          chatRequest.Model,
		"messages_count": len(chatRequest.Messages),
		"stream":         chatRequest.Stream,
		"temperature":    chatRequest.Temperature,
		"max_tokens":     chatRequest.MaxTokens,
	})
	
	// 检查模型是否支持
	modelConfig, ok := CONFIG.MODEL_CONFIG[chatRequest.Model]
	if !ok {
		logError(requestID, "不支持的模型: "+chatRequest.Model, nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "不支持的模型: " + chatRequest.Model,
				"type":    "invalid_request_error",
				"param":   "model",
				"code":    nil,
			},
		})
		return
	}
	
	// 配置选项
	params := map[string]interface{}{
		"temperature":       chatRequest.Temperature,
		"max_tokens":        chatRequest.MaxTokens,
		"presence_penalty":  chatRequest.PresencePenalty,
		"frequency_penalty": chatRequest.FrequencyPenalty,
		"top_p":             chatRequest.TopP,
		"top_k":             chatRequest.TopK,
	}
	configOpt := ConfigOpt(params, modelConfig)
	
	// 准备E2B请求
	e2bRequest, err := PrepareChatRequest(modelConfig, requestID, chatRequest, configOpt)
	if err != nil {
		logError(requestID, "准备聊天请求失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "准备请求失败: " + err.Error(),
				"type":    "server_error",
				"param":   nil,
				"code":    nil,
			},
		})
		return
	}
	
	logInfo(requestID, "发送到E2B的请求", map[string]interface{}{
		"model":          e2bRequest.Model.Name,
		"messages_count": len(e2bRequest.Messages),
		"config":         e2bRequest.Config,
	})
	
	// 发送请求到E2B
	requestData, err := json.Marshal(e2bRequest)
	if err != nil {
		logError(requestID, "请求序列化失败", err)
		handleInternalErrorGin(c, requestID, "请求序列化失败: "+err.Error())
		return
	}
	
	client := &http.Client{}
	req, err := http.NewRequest("POST", CONFIG.API.BASE_URL+"/api/chat", bytes.NewBuffer(requestData))
	if err != nil {
		logError(requestID, "创建HTTP请求失败", err)
		handleInternalErrorGin(c, requestID, "创建HTTP请求失败: "+err.Error())
		return
	}
	
	// 设置请求头
	for key, value := range CONFIG.DEFAULT_HEADERS {
		req.Header.Set(key, value)
	}
	
	// 发送请求并记录时间
	fetchStartTime := time.Now()
	resp, err := client.Do(req)
	fetchEndTime := time.Now()
	
	if err != nil {
		logError(requestID, "请求E2B失败", err)
		handleInternalErrorGin(c, requestID, "请求上游服务失败: "+err.Error())
		return
	}
	defer resp.Body.Close()
	
	// 解析E2B响应
	var e2bResponse E2BResponse
	if err := json.NewDecoder(resp.Body).Decode(&e2bResponse); err != nil {
		logError(requestID, "解析E2B响应失败", err)
		handleInternalErrorGin(c, requestID, "解析上游服务响应失败: "+err.Error())
		return
	}
	
	logInfo(requestID, fmt.Sprintf("收到E2B的响应: %d, 耗时: %dms", resp.StatusCode, fetchEndTime.Sub(fetchStartTime).Milliseconds()), map[string]interface{}{
		"status":           resp.StatusCode,
		"has_code":         e2bResponse.Code != "",
		"has_text":         e2bResponse.Text != "",
		"response_preview": truncateString(e2bResponse.Code+e2bResponse.Text, 100),
	})
	
	// 提取响应内容
	chatMessage := strings.TrimSpace(e2bResponse.Code)
	if chatMessage == "" {
		chatMessage = strings.TrimSpace(e2bResponse.Text)
	}
	
	if chatMessage == "" {
		logError(requestID, "E2B没有返回有效响应", nil)
		handleInternalErrorGin(c, requestID, "未从上游服务获取到响应")
		return
	}
	
	// 根据请求类型返回流式或普通响应
	if chatRequest.Stream {
		handleStreamResponseGin(c, chatMessage, chatRequest.Model, requestID)
	} else {
		handleNormalResponseGin(c, chatMessage, chatRequest.Model, requestID)
	}
}

// 使用 Gin 处理内部错误
func handleInternalErrorGin(c *gin.Context, requestID, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"message": message + " 请求失败，可能是上下文超出限制或其他错误，请稍后重试。",
			"type":    "server_error",
			"param":   nil,
			"code":    nil,
		},
	})
}

// 使用 Gin 处理普通响应
func handleNormalResponseGin(c *gin.Context, chatMessage string, model string, requestID string) {
	logInfo(requestID, fmt.Sprintf("处理普通响应，内容长度: %d 字符", len(chatMessage)))
	
	response := ChatCompletionResponse{
		ID:      GenerateUUID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []ChatChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: chatMessage,
				},
				FinishReason: "stop",
			},
		},
		Usage: nil,
	}
	
	c.JSON(http.StatusOK, response)
	logInfo(requestID, "返回普通响应成功")
}

// 使用 Gin 处理流式响应
func handleStreamResponseGin(c *gin.Context, chatMessage string, model string, requestID string) {
	logInfo(requestID, fmt.Sprintf("处理流式响应，内容长度: %d 字符", len(chatMessage)))
	
	// 设置响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)
	
	// 分段发送响应
	index := 0
	for index < len(chatMessage) {
		// 计算随机块大小
		chunkSize := rand.Intn(15) + 15 // 15到29之间
		if index+chunkSize > len(chatMessage) {
			chunkSize = len(chatMessage) - index
		}
		
		chunk := chatMessage[index:index+chunkSize]
		index += chunkSize
		
		// 创建事件数据
		eventData := struct {
			ID      string `json:"id"`
			Object  string `json:"object"`
			Created int64  `json:"created"`
			Model   string `json:"model"`
			Choices []struct {
				Index        int         `json:"index"`
				Delta        interface{} `json:"delta"`
				FinishReason *string     `json:"finish_reason"`
			} `json:"choices"`
		}{
			ID:      GenerateUUID(),
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   model,
			Choices: []struct {
				Index        int         `json:"index"`
				Delta        interface{} `json:"delta"`
				FinishReason *string     `json:"finish_reason"`
			}{
				{
					Index: 0,
					Delta: map[string]string{
						"content": chunk,
					},
					FinishReason: nil,
				},
			},
		}
		
		// 如果是最后一个分块，设置finish_reason
		if index >= len(chatMessage) {
			finishReason := "stop"
			eventData.Choices[0].FinishReason = &finishReason
		}
		
		// 序列化事件数据
		eventJSON, err := json.Marshal(eventData)
		if err != nil {
			logError(requestID, "序列化事件数据失败", err)
			return
		}
		
		// 写入事件流
		fmt.Fprintf(c.Writer, "data: %s\n\n", eventJSON)
		c.Writer.Flush()
		
		// 添加小延迟模拟真实速度
		time.Sleep(50 * time.Millisecond)
	}
	
	// 发送结束标记
	fmt.Fprint(c.Writer, "data: [DONE]\n\n")
	c.Writer.Flush()
	
	logInfo(requestID, "流式响应完成")
}

// 日志工具函数
func logInfo(requestID string, message string, data ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	if len(data) > 0 {
		jsonData, err := json.Marshal(data[0])
		if err != nil {
			log.Printf("[%s] INFO: %s - 无法序列化日志数据: %v", timestamp, message, err)
			return
		}
		
		// 对于大型响应体，可能需要截断
		dataStr := string(jsonData)
		if len(dataStr) > 500 {
			dataStr = dataStr[:500] + "...(truncated)"
		}
		log.Printf("[%s][%s] INFO: %s - %s", timestamp, requestID, message, dataStr)
	} else {
		log.Printf("[%s][%s] INFO: %s", timestamp, requestID, message)
	}
}

func logError(requestID string, message string, err error) {
	timestamp := time.Now().Format(time.RFC3339)
	log.Printf("[%s][%s] ERROR: %s - %v", timestamp, requestID, message, err)
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	
	// 设置版本和变体位
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 1
	
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// ConfigOpt 配置选项
func ConfigOpt(params map[string]interface{}, modelConfig ModelConfig) map[string]interface{} {
	if modelConfig.OptMax == (OptMax{}) {
		return nil
	}
	
	optionsMap := map[string]string{
		"temperature":       "TemperatureMax",
		"max_tokens":        "MaxTokensMax",
		"presence_penalty":  "PresencePenaltyMax",
		"frequency_penalty": "FrequencyPenaltyMax",
		"top_p":             "TopPMax",
		"top_k":             "TopKMax",
	}
	
	constrainedParams := make(map[string]interface{})
	for key, value := range params {
		maxKey, ok := optionsMap[key]
		if !ok || value == nil {
			continue
		}
		
		switch maxKey {
		case "TemperatureMax":
			if temp, ok := value.(float64); ok && modelConfig.OptMax.TemperatureMax > 0 {
				constrainedParams[key] = math.Min(temp, modelConfig.OptMax.TemperatureMax)
			}
		case "MaxTokensMax":
			if tokens, ok := value.(int); ok && modelConfig.OptMax.MaxTokensMax > 0 {
				constrainedParams[key] = int(math.Min(float64(tokens), float64(modelConfig.OptMax.MaxTokensMax)))
			}
		case "PresencePenaltyMax":
			if penalty, ok := value.(float64); ok && modelConfig.OptMax.PresencePenaltyMax > 0 {
				constrainedParams[key] = math.Min(penalty, modelConfig.OptMax.PresencePenaltyMax)
			}
		case "FrequencyPenaltyMax":
			if penalty, ok := value.(float64); ok && modelConfig.OptMax.FrequencyPenaltyMax > 0 {
				constrainedParams[key] = math.Min(penalty, modelConfig.OptMax.FrequencyPenaltyMax)
			}
		case "TopPMax":
			if topP, ok := value.(float64); ok && modelConfig.OptMax.TopPMax > 0 {
				constrainedParams[key] = math.Min(topP, modelConfig.OptMax.TopPMax)
			}
		case "TopKMax":
			if topK, ok := value.(int); ok && modelConfig.OptMax.TopKMax > 0 {
				constrainedParams[key] = int(math.Min(float64(topK), float64(modelConfig.OptMax.TopKMax)))
			}
		}
	}
	
	return constrainedParams
}

// ProcessMessageContent 处理消息内容
func ProcessMessageContent(content interface{}) string {
	switch v := content.(type) {
	case string:
		return v
	case []interface{}:
		var textParts []string
		for _, item := range v {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if itemMap["type"] == "text" {
					if text, ok := itemMap["text"].(string); ok {
						textParts = append(textParts, text)
					}
				}
			}
		}
		return strings.Join(textParts, "\n")
	case map[string]interface{}:
		if text, ok := v["text"].(string); ok {
			return text
		}
	}
	
	return ""
}

// TransformMessages 转换消息
func TransformMessages(messages []ChatMessage) []ChatMessage {
	if len(messages) == 0 {
		return messages
	}
	
	var mergedMessages []ChatMessage
	var lastMessage *ChatMessage
	
	for _, current := range messages {
		currentContent := ProcessMessageContent(current.Content)
		if currentContent == "" {
			continue
		}
		
		if lastMessage != nil && lastMessage.Role == current.Role {
			lastContent := ProcessMessageContent(lastMessage.Content)
			if lastContent != "" {
				lastMessage.Content = lastContent + "\n" + currentContent
				continue
			}
		}
		
		messageCopy := current
		mergedMessages = append(mergedMessages, messageCopy)
		lastMessage = &mergedMessages[len(mergedMessages)-1]
	}
	
	// 转换为E2B要求的格式
	var transformed []ChatMessage
	for _, msg := range mergedMessages {
		content := ProcessMessageContent(msg.Content)
		switch msg.Role {
		case "system", "user":
			transformed = append(transformed, ChatMessage{
				Role: "user",
				Content: []TextContent{
					{
						Type: "text",
						Text: content,
					},
				},
			})
		case "assistant":
			transformed = append(transformed, ChatMessage{
				Role: "assistant",
				Content: []TextContent{
					{
						Type: "text",
						Text: content,
					},
				},
			})
		default:
			transformed = append(transformed, msg)
		}
	}
	
	return transformed
}

// PrepareChatRequest 准备聊天请求
func PrepareChatRequest(modelConfig ModelConfig, requestID string, request ChatRequest, config map[string]interface{}) (E2BRequest, error) {
	logInfo(requestID, fmt.Sprintf("准备聊天请求, 模型: %s, 消息数: %d", modelConfig.Name, len(request.Messages)))
	
	transformedMessages := TransformMessages(request.Messages)
	logInfo(requestID, fmt.Sprintf("转换后的消息数量: %d", len(transformedMessages)))
	
	if config == nil {
		config = map[string]interface{}{
			"model": modelConfig.ID,
		}
	}
	
	e2bRequest := E2BRequest{
		UserID:   GenerateUUID(),
		Messages: transformedMessages,
		Template: map[string]interface{}{
			"text": map[string]interface{}{
				"name":         CONFIG.MODEL_PROMPT,
				"lib":          []string{""},
				"file":         "pages/ChatWithUsers.txt",
				"instructions": modelConfig.SystemPrompt,
				"port":         nil,
			},
		},
		Model: struct {
			ID         string `json:"id"`
			Provider   string `json:"provider"`
			ProviderID string `json:"providerId"`
			Name       string `json:"name"`
			MultiModal bool   `json:"multiModal"`
		}{
			ID:         modelConfig.ID,
			Provider:   modelConfig.Provider,
			ProviderID: modelConfig.ProviderID,
			Name:       modelConfig.Name,
			MultiModal: modelConfig.MultiModal,
		},
		Config: config,
	}
	
	return e2bRequest, nil
}

// 辅助函数：截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// 辅助函数：掩盖字符串
func maskString(s string, showLen int) string {
	if len(s) <= showLen {
		return s
	}
	return s[:showLen] + "..."
}
