package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/ws-ai/pkg/dalle"
)

func main() {
	log.Info("generating an image...")

	prompt := "A cute pink dragon"
	prompt = strings.TrimSpace(prompt)

	req := dalle.DalleRequest{
		Model: "dall-e-3",
		// Model:  "dall-e-2",
		Prompt: prompt,
		Size:   "1024x1024",
		// Size: "512x512",
		// Size: "256x256",
		N: 1,
	}
	imageURL, err := dalle.GenerateImage(req)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}

	log.Info(imageURL)
	log.Info("done")
}
