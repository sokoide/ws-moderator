package dalle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync/atomic"
)

// types
type Request struct {
	Prompt string
}

type Response struct {
	Url       string
	Succeeded bool
}

type DalleRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt,omitempty"`
	Size      string `json:"size"`
	Quality   string `json:"quality,omitempty"`
	N         int    `json:"n"`
	Image     string `json:"image,omitempty"`
	ImageData []byte `json:"-"`
}

type DalleResponse struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

// globals
var apiKey string
var dalleConns int32

// functions
func GetConns() int32 {
	return atomic.LoadInt32(&dalleConns)
}

func GenerateImage(req DalleRequest) (string, error) {
	url := "https://api.openai.com/v1/images/generations" // Correct API endpoint
	atomic.AddInt32(&dalleConns, 1)
	defer atomic.AddInt32(&dalleConns, -1)

	requestBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
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

	var dalleResp DalleResponse
	err = json.Unmarshal(body, &dalleResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(dalleResp.Data) > 0 {
		return dalleResp.Data[0].URL, nil
	}

	return "", fmt.Errorf("no image URL returned")
}

func init() {
	apiKey = os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable not set")
	}
}
