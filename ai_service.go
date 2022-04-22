package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

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

var url string = "https://api.openai.com/v1/engines/text-davinci-002/completions"

func friendlyBot(ainame string, m Message) (string, error) {
	question := m.Text[4:]

	prompt :=
		"The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly." +
			m.Firstname + ": Hello, who are you?" +
			ainame + ": I am an AI created by Developer santuy. How can I help you today?" +
			m.Firstname + ":" + question +
			ainame + ":"
	payload := &PayloadAI{
		Prompt:      prompt,
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

func sarcasticBot(ainame string, m Message) (string, error) {
	question := m.Text[4:]

	prompt :=
		ainame + " is a chatbot that reluctantly answers questions with sarcastic responses:" +
			m.Firstname + ": Sekilo berapa pound?" +
			ainame + ": Ini lagi? Sekilo itu 2.2 pound. Nanti catet ya" +
			m.Firstname + ": Apa kepanjangan dari HTML?" +
			ainame + ": Emang google sibuk? Hypertext Markup Language. T nya untuk TOLOL" +
			m.Firstname + ": Apa arti kehidupan?" +
			ainame + ": Gak yakin. Nanti tanya ke Google." +
			m.Firstname + ":" + question +
			ainame + ":"

	payload := &PayloadAI{
		Prompt:      prompt,
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
