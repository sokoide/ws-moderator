package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	Prompt string
}

type Response struct {
	Text      string
	Succeeded bool
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

type ClaudeResponse struct {
	Content []ClaudeContent `json:"content"`
}

type ClaudeContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func callClaudeAPI(apiKey string, model string, history []ClaudeMessage) (string, error) {
	url := "https://api.anthropic.com/v1/messages"

	requestBody, err := json.Marshal(ClaudeRequest{
		Model:     model,
		MaxTokens: 1024,
		Messages:  history,
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
		// Update conversation history with the new assistant message
		history = append(history, ClaudeMessage{
			Role:    "assistant",
			Content: claudeResp.Content[0].Text,
		})
		return claudeResp.Content[0].Text, nil
	}

	return "", fmt.Errorf("no content in response")
}

// id is an email
func StartConversation(id string, cin chan Request, cout chan Response) {
	fmt.Printf("Starting ID:%s\n", id)
	model := "claude-3-5-sonnet-20240620"
	// TODO: store history in mongodb
	var history []ClaudeMessage

	for {
		req, ok := <-cin
		if !ok {
			fmt.Printf("ID:%s channel closed, exiting\n", id)
			return
		}

		userInput := strings.TrimSpace(req.Prompt)

		// Update conversation history with the new user message
		history = append(history, ClaudeMessage{
			Role:    "user",
			Content: userInput,
		})

		response, err := callClaudeAPI(apiKey, model, history)
		if err != nil {
			fmt.Printf("ID:%s error calling Claude API: %v\n", id, err)
			cout <- Response{Succeeded: false, Text: ""}
		}

		cout <- Response{Succeeded: true, Text: response}
	}
}

var apiKey string

func init() {
	apiKey = os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		panic("CLAUDE_API_KEY environment variable not set")
	}
}
