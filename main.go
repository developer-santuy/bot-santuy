package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type webHookReqBody struct {
	Message `json:"message"`
}

type Message struct {
	MessageID      int64  `json:"message_id"`
	Text           string `json:"text"`
	Chat           `json:"chat"`
	ReplyToMessage *Message `json:"reply_to_message,omitempty"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

type joke struct {
	Value struct {
		Joke string `json:"joke"`
	} `json:"value"`
}

type sendMessageReqBody struct {
	ChatID           int64  `json:"chat_id"`
	Text             string `json:"text"`
	ReplyToMessageID int64  `json:"reply_to_message_id"`
}

type idea struct {
	Content string `json:"content"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body := &webHookReqBody{}

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		fmt.Println("webhookHandler: error %w", err)
		return
	}

	err = sendReply(body.Message)

	if err != nil {
		fmt.Printf("Something happen: %v", err)
	}
}

func sendReply(m Message) error {
	fmt.Println("sendReply called")

	var botToken string = os.Getenv("BOT_TOKEN")

	commands := strings.ToLower(m.Text)
	text := ""
	var err error = nil

	switch commands {
	case "/joke":
		text, err = jokeFetcher()
	case "/idea":
		text, err = scrapIdea()
	}

	AIname := strings.Split(m.Text, ",")[0]

	if AIname == "AI" {
		text, err = friendlyBot(AIname, m)
	}

	if AIname == "Cuy" {
		text, err = sarcasticBot(AIname, m)
	}

	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("bot Commands: error %v", err)
	}

	reqBody := &sendMessageReqBody{
		ChatID:           m.Chat.ID,
		Text:             text,
		ReplyToMessageID: m.MessageID,
	}

	fmt.Printf("reqbody %v:", reqBody)
	reqBytes, err := json.Marshal(reqBody)

	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("bot Commands: error %v", err)
	}

	resp, err := http.Post(
		"https://api.telegram.org/bot"+botToken+"/sendMessage",
		"application/json",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("failed make request to telegram: error %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Unexpected Status" + resp.Status)
	}

	return err
}

func jokeFetcher() (string, error) {
	resp, err := http.Get("http://api.icndb.com/jokes/random")
	c := &joke{}
	if err != nil {
		return "", fmt.Errorf("jokeFetcher: error %w", err)
	}

	err = json.NewDecoder(resp.Body).Decode(c)
	return c.Value.Joke, err
}

func scrapIdea() (string, error) {
	url := "https://thisideadoesnotexist.com/"

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("scrapIdea: error %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("goquery unable read body: %w", err)
	}

	idea := &idea{}

	idea.Content = doc.Find("h2").Text()

	return idea.Content, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	})

	http.HandleFunc("/webhook", webhookHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	http.ListenAndServe(":"+port, nil)
}
