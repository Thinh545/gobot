package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	FBMessageURL    = "https://graph.facebook.com/v3.1/me/messages"
	PageToken       = "EAAMp67qUXs0BAMZBOaEKPnJhIsmQP3ZAHwhnYqBnOoExoJrsmZCd04uRPpxvkPSelYTM1EMQZCp6iDHXEZBR3HFnwMwxnG68IrQ7KD272Pu60n1lSGVRL4yRkwtQbgli098wkbGaax9PZCCdIRau4qHcZBAowoZBid6TW8oGaREDnSvgyvKwZAQ7E"
	MessageResponse = "RESPONSE"
	MarkSeen        = "mark_seen"
	TypingOn        = "typing_on"
	TypingOff       = "typing_off"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", chatbotHandler)
	if err := http.ListenAndServe("", r); err != nil {
		log.Fatal(err.Error())
	}
}

func chatbotHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		verifyWebhook(w, r)

	case "POST":
		processWebhook(w, r)

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalf("Không hỗ trợ phương thức HTTP %v", r.Method)
	}
}

func verifyWebhook(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	challenge := r.URL.Query().Get("hub.challenge")
	token := r.URL.Query().Get("hub.verify_token")

	if mode == "subscribe" && token == "Go_TestBot" {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Error, wrong validation token"))
	}
}

func processWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Message not supported"))
		return
	}

	if req.Object == "page" {
		for _, entry := range req.Entry {
			for _, event := range entry.Messaging {
				if event.Message != nil {
					sendAction(event.Sender, MarkSeen)
					sendAction(event.Sender, TypingOn)
					sendText(event.Sender, strings.ToUpper(event.Message.Text))
					sendAction(event.Sender, TypingOff)
				}
			}
		}

		w.WriteHeader(200)
		w.Write([]byte("Got your message"))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Message not supported"))
	}

}

func sendFBRequest(url string, m interface{}) error {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(&m)
	if err != nil {
		log.Fatalf("sendFBRequest:json.NewEncoder: " + err.Error())
		return err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalf("sendFBRequest:http.NewRequest:" + err.Error())
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.URL.RawQuery = "access_token=" + PageToken
	client := &http.Client{Timeout: time.Second * 30}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("sendFBRequest:client.Do: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	return nil
}

func sendText(recipient *User, message string) error {
	m := ResponseMessage{
		MessageType: MessageResponse,
		Recipient:   recipient,
		Message: &ResMessage{
			Text: message,
		},
	}
	return sendFBRequest(FBMessageURL, &m)
}

func sendAction(recipient *User, action string) error {
	m := ResponseMessage{
		MessageType: MessageResponse,
		Recipient:   recipient,
		Action:      action,
	}
	return sendFBRequest(FBMessageURL, &m)
}
