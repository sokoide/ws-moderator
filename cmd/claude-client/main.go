package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeRequest struct {
	Model     string         `json:"model"`
	MaxTokens int            `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

type ClaudeResponse struct {
	Content []ClaudeContent `json:"content"`
}

type ClaudeContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func callClaudeAPI(apiKey string, model string, prompt string) (string, error) {
	url := "https://api.anthropic.com/v1/messages"

	messages := []ClaudeMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody, err := json.Marshal(ClaudeRequest{
		Model:     model,
		MaxTokens: 1024,
		Messages:  messages,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01") // Example version, adjust as necessary

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", body)
	}

	var claudeResp ClaudeResponse
	err = json.Unmarshal(body, &claudeResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(claudeResp.Content) > 0 {
		return claudeResp.Content[0].Text, nil
	}

	return "", fmt.Errorf("no content in response")
}

func main() {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		fmt.Println("CLAUDE_API_KEY environment variable not set")
		return
	}

	model := "claude-3-5-sonnet-20240620" // Replace with the appropriate model
	prompt := "Hello, world"
	response, err := callClaudeAPI(apiKey, model, prompt)
	if err != nil {
		fmt.Printf("Error calling Claude API: %v\n", err)
		return
	}

	fmt.Printf("Claude response: %s\n", response)
}

