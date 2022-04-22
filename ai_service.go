package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type UpdateMessage struct {
	Message Message
}

type PayloadAI struct {
	Prompt      string `json:"prompt"`
	MaxTokens   uint8  `json:"max_tokens"`
	Temperature uint8  `json:"temperature"`
	TopP        uint8  `json:"top_p"`
}

type AIResponse struct {
	Id      string    `json:"id"`
	Choices []Choices `json:"choices"`
}

type Choices struct {
	Text string `json:"text"`
}

func friendlyBot(commands string) (string, error) {
	url := "https://api.openai.com/v1/engines/text-davinci-002/completions"

	question := commands[4:]
	payload := &PayloadAI{
		Prompt:      question,
		MaxTokens:   100,
		Temperature: 0,
		TopP:        1,
	}

	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("friendlyBot - marshalling: error %w", err)
	}

	client := &http.Client{}
	aiResponse := &AIResponse{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshalledPayload))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("AI_API_KEY"))

	if err != nil {
		return "", fmt.Errorf("friendlyBot - create request: error %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("friendlyBot - do request: error %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(aiResponse)

	if err != nil {
		return "", fmt.Errorf("friendlyBot - decode body: error %w", err)
	}

	fmt.Printf("airesponse %v", aiResponse)
	return aiResponse.Choices[0].Text, nil
}
