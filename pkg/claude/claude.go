package claude

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

// types
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

// globals
var apiKey string
var claudeConns int32

// functions
func GetConns() int32 {
	return atomic.LoadInt32(&claudeConns)
}

func callClaudeAPI(apiKey string, model string, history *[]ClaudeMessage) (string, error) {
	url := "https://api.anthropic.com/v1/messages"

	requestBody, err := json.Marshal(ClaudeRequest{
		Model:     model,
		MaxTokens: 1024,
		Messages:  *history,
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

	body, err := io.ReadAll(resp.Body)
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
		*history = append(*history, ClaudeMessage{
			Role:    "assistant",
			Content: claudeResp.Content[0].Text,
		})
		return claudeResp.Content[0].Text, nil
	}

	return "", fmt.Errorf("no content in response")
}

func serializeHistory(id string, history []ClaudeMessage) error {
	filepath := fmt.Sprintf("%s.gob", id)
	file, err := os.Create(filepath)
	if err != nil {
		log.Errorf("Error creating file: %v", err)
		return err
	}
	defer file.Close()

	// Create a new encoder
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(history); err != nil {
		log.Errorf("Error encoding GOB: %v", err)
		return err
	}

	log.Infof("Serialization of %s succeeded", filepath)
	return nil
}

func deserializeHistory(id string) []ClaudeMessage {
	var history []ClaudeMessage

	filepath := fmt.Sprintf("%s.gob", id)
	file, err := os.Open(filepath)
	if err != nil {
		log.Errorf("Error opening file: %v", err)
		return history
	}
	defer file.Close()

	// Create a new decoder
	decoder := gob.NewDecoder(file)

	// Decode the data into the variable
	if err := decoder.Decode(&history); err != nil {
		log.Errorf("Error decoding GOB: %v", err)
		return history
	}

	log.Infof("Deserialization of %s succeeded", filepath)
	return history
}

// id is an email
func StartConversation(id string, cin chan Request, cout chan Response) {
	fmt.Printf("Starting ID:%s\n", id)
	model := "claude-3-5-sonnet-20240620"
	var history []ClaudeMessage
	atomic.AddInt32(&claudeConns, 1)
	defer atomic.AddInt32(&claudeConns, -1)

	history = deserializeHistory(id)

	for {
		req, ok := <-cin
		if !ok {
			fmt.Printf("ID:%s channel closed, exiting\n", id)
			err := serializeHistory(id, history)
			if err == nil {
				log.Info("serialization succeeded")
			} else {
				log.Errorf("serialization failed. %v", err)
			}
			return
		}

		userInput := strings.TrimSpace(req.Prompt)

		// Update conversation history with the new user message
		history = append(history, ClaudeMessage{
			Role:    "user",
			Content: userInput,
		})

		response, err := callClaudeAPI(apiKey, model, &history)
		if err != nil {
			fmt.Printf("ID:%s error calling Claude API: %v\n", id, err)
			cout <- Response{Succeeded: false, Text: fmt.Sprintf("err: %v", err)}
		}

		cout <- Response{Succeeded: true, Text: response}
	}
}

func init() {
	apiKey = os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		panic("CLAUDE_API_KEY environment variable not set")
	}
}
