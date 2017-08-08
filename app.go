package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
		ID        string         `json:"id,omitempty"`
		Time      int         `json:"time,omitempty"`
		Messaging []Messaging `json:"messaging,omitempty"`
	} `json:"entry,omitempty"`
}

type Messaging struct {
	Sender    User    `json:"sender,omitempty"`
	Recipient User    `json:"recipient,omitempty"`
	Timestamp int     `json:"timestamp,omitempty"`
	Message   Message `json:"message,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type Message struct {
	MID        string `json:"mid,omitempty"`
	Text       string `json:"text,omitempty"`
	QuickReply *struct {
		Payload string `json:"payload,omitempty"`
	} `json:"quick_reply,omitempty"`
	Attachments *[]Attachment `json:"attachments,omitempty"`
	Attachment  *Attachment   `json:"attachment,omitempty"`
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
	challenge := r.URL.Query().Get("hub.challenge")
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
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
	client := &http.Client{}
	response := Response{
		Recipient: User{
			ID: "1125206514248128",
		},
		Message: Message{
			Text: "Hello world",
		},
	}
	fmt.Printf("%+v\n", response)
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&response)
	req, err := http.NewRequest("POST", FACEBOOK_API+"?access_token=EAAbAxXjuZAdgBAGaQNmhQ5NaF8q0pEWRyFx0rZCIwKDrunKwYMofxpNj6d1ILFOW3bJyOlu9m3ZApP8HGZAqQuVzhppzqOFZBCNMyOXZB7QCgxiElv0EZA6eGKYLwIqwrRVV00ZCLnwJVeP2D811ZAv2ABRDIfYt25wVPdMYSOGcktwZDZD", body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		fmt.Println("here")
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("here 2")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("here 3")
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
}

func MessagesEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("is here")
	var callback Callback
	err := json.NewDecoder(r.Body).Decode(&callback)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(data))
	fmt.Println(callback)
	fmt.Println("goes here")
	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				ProcessMessage(event)
			}
		}
		w.WriteHeader(200)
		w.Write([]byte("Got your message"))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Message not supported"))	
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/webhook", VertificationEndpoint).Methods("GET")
	r.HandleFunc("/webhook", MessagesEndpoint).Methods("POST")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal(err)
	}
}
