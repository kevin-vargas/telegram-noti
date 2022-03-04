package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/kevin-vargas/sidecar-log/pubsub"
)

const (
	TOPIC_NOTIFICATIONS = "notification"
	URL_BASE            = "https://api.telegram.org/bot"
	CHAT_ID             = -1001767276075
	METHOD              = "sendMessage"
)

func getURL() string {
	return URL_BASE + os.Getenv("TOKEN") + "/" + METHOD
}

type Message struct {
	Chat int    `json:"chat_id"`
	Text string `json:"text"`
}

func makeHandler(p pubsub.MQTTI) func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		payload := msg.Payload()
		message := Message{
			Chat: CHAT_ID,
			Text: string(payload),
		}
		b, err := json.Marshal(message)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(getURL())
		fmt.Println(string(b))
		resp, err := http.Post(getURL(), "application/json", bytes.NewBuffer(b))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

	}
}

func main() {
	syncChan := make(chan bool, 0)
	fmt.Println("Init telegram-noti")
	godotenv.Load()
	client := pubsub.New()
	handler := makeHandler(client)
	client.SubscribeWithCB(TOPIC_NOTIFICATIONS, handler)
	<-syncChan
}
