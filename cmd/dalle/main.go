package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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

// Corrected Generate Image function with correct API endpoint
func generateImage(apiKey string, req DalleRequest) (string, error) {
	url := "https://api.openai.com/v1/images/generations" // Correct API endpoint

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

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable not set")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your prompt: ")
	prompt, _ := reader.ReadString('\n')
	prompt = strings.TrimSpace(prompt)

	if prompt == "" {
		fmt.Println("Prompt cannot be empty")
		return
	}

	req := DalleRequest{
		Model:  "dall-e-3", // You can allow users to choose between "dall-e-3" and "dall-e-2" if needed
		Prompt: prompt,
		Size:   "1024x1024", // You can allow users to customize the size if needed
		N:      1,
	}

	imageURL, err := generateImage(apiKey, req)
	if err != nil {
		fmt.Printf("Error generating image: %v\n", err)
		return
	}

	fmt.Printf("Generated image URL: %s\n", imageURL)
}
