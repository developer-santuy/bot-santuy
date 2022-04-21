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
	"github.com/joho/godotenv"
)

type webHookReqBody struct {
	Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat `json:"chat"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type joke struct {
	Value struct {
		Joke string `json:"joke"`
	} `json:"value"`
}

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type idea struct {
	Content string `json:"content"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body := &webHookReqBody{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		log.Printf("An error occured (webHookHandler)")
		log.Printf("error %v", err)
		log.Panic(err)
		return
	}

	commands := strings.ToLower(body.Message.Text)

	err := errors.New("")

	switch commands {
	case "/joke":
		err = sendReply(body.Message.Chat.ID, jokeFetcher)
	case "/idea":
		err = sendReply(body.Message.Chat.ID, scrapIdea)
	}

	if err != nil {
		log.Fatal(err)
		return
	}

}

type Commands func() (string, error)

func sendReply(chatID int64, commands Commands) error {
	fmt.Println("sendReply called")

	var botToken string = os.Getenv("BOT_TOKEN")

	text, err := commands()
	if err != nil {
		return err
	}

	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   text,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		"https://api.telegram.org/bot"+botToken+"/sendMessage",
		"application/json",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return err
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
		return "", nil
	}

	err = json.NewDecoder(resp.Body).Decode(c)
	return c.Value.Joke, err
}

func scrapIdea() (string, error) {
	url := "https://thisideadoesnotexist.com/"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	idea := &idea{}

	idea.Content = "Idea:" + doc.Find("h2").Text()

	return idea.Content, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
