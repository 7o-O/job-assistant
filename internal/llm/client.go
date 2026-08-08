package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 调用模型
type Client struct {
	BaseURL string
	APIKey  string
	Model   string
	HTTP    *http.Client
}

// 对话信息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 发给大模型请求数据
type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// 大模型返回数据
type chatRespose struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// 创建客户端
func NewClient(baseURL, apiKey, model string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
		HTTP:    &http.Client{Timeout: 180 * time.Second},
	}
}

// 调用大模型，返回回答
func (c *Client) Chat(ctx context.Context, messages []Message) (string, error) {
	requestBody := chatRequest{
		Model:    c.Model,
		Messages: messages,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("失败，%w", err)
	}

	//创建post请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	//设置请求头
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	//发送请求
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("调用失败，%w", err)
	}
	defer resp.Body.Close()

	//读取响应内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("大模型返回错误，状态码=%d，内容=%s", resp.StatusCode, string(data))
	}

	//解析响应
	var result chatRespose
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("解析失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("无返回")
	}
	return result.Choices[0].Message.Content, nil
}
