package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	TOKEN        = "you know nothing, john snow"
	FACEBOOK_API = "https://graph.facebook.com/v2.6/me/messages"
	IMAGE        = "http://37.media.tumblr.com/e705e901302b5925ffb2bcf3cacb5bcd/tumblr_n6vxziSQD11slv6upo3_500.gif"
)

type Callback struct {
	Object string `json:"object,omitempty"`
	Entry  []struct {
		ID        string      `json:"id,omitempty"`
		Time      time.Time   `json:"time,omitempty"`
		Messaging []Messaging `json:"messaging,omitempty"`
	} `json:"entry,omitempty"`
}

type Messaging struct {
	Sender    User      `json:"sender,omitempty"`
	Recipient User      `json:"recipient,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Message   Message   `json:"message,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type Message struct {
	MID        string `json:"mid,omitempty"`
	Text       string `json:"text,omitempty"`
	QuickReply struct {
		Payload string `json:"payload,omitempty"`
	} `json:"quick_reply,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Attachment  Attachment   `json:"attachment,omitempty"`
}

type Attachment struct {
	Type    string  `json:"type,omitempty"`
	Payload Payload `json:"payload,omitempty"`
}

type Response struct {
	Recipient User    `json:"recipient,omitempty"`
	Message   Message `json:"message,omitempty"`
}

type Payload struct {
	URL string `json:"url,omitempty"`
}

func VertificationEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	challenge := params["hub.challenge"]
	mode := params["hub.mode"]
	token := params["hub.verify_token"]
	
	fmt.Println(params)
	fmt.Println("im here")

	if mode != "" && token == TOKEN {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Error, wrong validation token"))
	}
}

func ProcessMessage(event Messaging) {
	response := Response{
		Recipient: event.Recipient,
		Message: Message{
			Attachment: Attachment{
				Type: "image",
				Payload: Payload{
					URL: IMAGE,
				},
			},
		},
	}

	client := &http.Client{}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(response)
	req, err := http.NewRequest("POST", FACEBOOK_API, body)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func MessagesEndpoint(w http.ResponseWriter, r *http.Request) {
	var callback Callback
	json.NewDecoder(r.Body).Decode(callback)
	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				ProcessMessage(event)
			}
		}
		w.WriteHeader(200)
		w.Write([]byte("Got your message"))
	}
	w.WriteHeader(404)
	w.Write([]byte("Message not supported"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/verify", VertificationEndpoint).Methods("GET")
	r.HandleFunc("/messages", MessagesEndpoint).Methods("POST")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal(err)
	}
}
