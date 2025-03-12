package domain

import (
	"UnoBackend/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiURL    = "https://api.deepseek.com/v1/chat/completions"
	apiKey    = "sk-09e51faee39f4a9a9358dbd732868b1f" // 从环境变量或配置读取更安全
	modelName = "deepseek-chat"                       // 根据实际模型名称修改
)

func GetDeepSeekChatCompletion(messages []model.ChatMessage) (string, error) {
	requestBody := model.ChatCompletionRequest{
		Model:       modelName,
		Messages:    messages,
		MaxTokens:   500,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	var response model.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decode response failed: %v", err)
	}

	if response.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}
